package raft

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/superioz/artemis/pkg/logger"
	"github.com/superioz/artemis/pkg/transport"
	"github.com/superioz/artemis/pkg/util"
	"github.com/superioz/artemis/raft/protocol"
	"time"
)

type State string

const (
	Follower  = "follower"
	Candidate = "candidate"
	Leader    = "leader"

	exchange = "artemis"

	defaultHeartbeat       = 1500 // in ms, def: 50
	defaultElectionTimeout = 2000 // in ms, def: 150
)

type Node struct {
	id    uuid.UUID
	state State

	leader uuid.UUID

	currentTerm uint64    // stable storage
	votedFor    uuid.UUID // stable storage
	log         Log       // stable storage

	commitIndex uint64               // volatile
	lastApplied uint64               // volatile
	nextIndex   map[uuid.UUID]uint64 // volatile for leaders
	matchIndex  map[uuid.UUID]uint64 // volatile for leaders

	transport transport.Interface

	heartbeatInterval time.Duration
	electionTimeout   time.Duration

	// TODO test value
	Inactive bool
}

func NewNode() Node {
	trans := transport.NewAMQPInterface(exchange)

	n := Node{
		id:                trans.State().Id(),
		state:             Follower,
		currentTerm:       0,
		log:               Log{},
		commitIndex:       0,
		lastApplied:       0,
		nextIndex:         make(map[uuid.UUID]uint64),
		matchIndex:        make(map[uuid.UUID]uint64),
		transport:         &trans,
		heartbeatInterval: defaultHeartbeat,
	}
	return n
}

func (n *Node) Up(brokerUrl string) {
	err := n.transport.Connect(brokerUrl)
	if err != nil || !n.transport.State().Connected() {
		logger.Err(fmt.Sprintf("couldn't connect to broker %s", brokerUrl), err)
		return
	}

	for n.transport.State().Connected() {
		switch n.state {
		case Follower:
			n.followerLoop()
			break
		case Candidate:
			n.candidateLoop()
			break
		case Leader:
			n.leaderLoop()
			break
		}
	}
}

func (n *Node) Down() {
	_ = n.transport.Disconnect()
}

func (n *Node) followerLoop() {
	timeout := n.generateTimeout()
	tc := time.NewTimer(timeout)
	pc := n.transport.Receive()

	// reset some values
	n.votedFor = uuid.UUID{}

followerLoop:
	for n.state == Follower {
		select {
		case p := <-pc:
			fmt.Println(n.id.String()+" | Follower: Received packet: ", p)

			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				fmt.Println("couldn't decode packet", err)
				break
			}

			switch m.(type) {
			case *protocol.RequestVoteCall:
				// TODO if requestvote: check if voteFor != null
				fmt.Println(n.id.String()+" | Follower: Received request vote: ", m)

				d, _ := transport.Encode(&protocol.RequestVoteRespond{
					Term:        n.currentTerm,
					VoteGranted: n.processVoteRequest(*m.(*protocol.RequestVoteCall)),
				})
				om := transport.OutgoingMessage{
					RoutingKey: p.Source.String(),
					Data:       d,
				}
				n.transport.Send() <- &om
				break
			case *protocol.RequestVoteRespond:
				// ignore
				break
			case *protocol.AppendEntriesCall:
				appendEntr := *m.(*protocol.AppendEntriesCall)

				// reset timeout
				timeout = n.generateTimeout()
				tc.Reset(timeout)

				// set leader
				n.leader, _ = uuid.FromString(appendEntr.LeaderId)

				// TODO if append entries: appendentries
				break
			case *protocol.AppendEntriesRespond:
				// ignore
				break
			}
			break
		case <-tc.C:
			fmt.Printf(n.id.String()+" | Follower: Timeout! (After %s)\n", timeout)

			// TODO remove test value
			if n.Inactive {
				break followerLoop
			}
			// TODO set state to candidate and send request vote rpc's
			n.state = Candidate
			break
		}
	}
}

func (n *Node) candidateLoop() {
	timeout := n.generateTimeout()

	tc := time.NewTimer(timeout)
	pc := n.transport.Receive()
	hb := time.NewTicker(n.heartbeatInterval * time.Millisecond)

	// increment term
	n.currentTerm++
	n.votedFor = n.id

	// send packet
	sendp := func() {
		d, _ := transport.Encode(&protocol.RequestVoteCall{
			Term:         n.currentTerm,
			CandidateId:  n.id.String(),
			LastLogIndex: n.log.LastLogEntry().Index,
			LastLogTerm:  n.log.LastLogEntry().Term,
		})
		m := transport.OutgoingMessage{
			RoutingKey: "broadcast.all",
			Data:       d,
		}
		n.transport.Send() <- &m
	}
	sendp()

candidateLoop:
	for n.state == Candidate {
		select {
		case p := <-pc:
			fmt.Println(n.id.String()+" | Candidate: Received packet: ", p)

			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				fmt.Println("couldn't decode packet", err)
				break
			}

			switch m.(type) {
			case *protocol.RequestVoteCall:

				break
			case *protocol.RequestVoteRespond:
				// TODO count negative and positive responds and calculate if he got the majority ..

				n.state = Leader
				fmt.Println(n.id.String() + " | Leader: Received vote.")
				break candidateLoop
			case *protocol.AppendEntriesCall:
				// TODO oh, there is already a leader? le me step back!
				break
			case *protocol.AppendEntriesRespond:

				break
			}
			// TODO if AppendEntries -> go back to being follower
			// TODO if RequestVote -> return false, as voteFor != null
			break
		case <-tc.C:
			fmt.Printf(n.id.String()+" | Candidate: Timeout! (After %s)\n", timeout)

			// TODO set state back to candidate and resent the packets ..
			break candidateLoop
		case <-hb.C:
			fmt.Println(n.id.String() + " | Candidate: Heartbeat!")
			sendp()

			break
		}
	}

	tc.Stop()
	hb.Stop()
}

func (n *Node) leaderLoop() {
	pc := n.transport.Receive()
	hb := time.NewTicker(n.heartbeatInterval * time.Millisecond)

	// send packet
	sendp := func() {
		lastEntry := n.log.LastLogEntry()
		d, _ := transport.Encode(&protocol.AppendEntriesCall{
			Term:         n.currentTerm,
			LeaderId:     n.id.String(),
			PrevLogIndex: lastEntry.Index,
			PrevLogTerm:  lastEntry.Term,
			Entries:      []*protocol.AppendEntry{},
			CommitIndex:  n.commitIndex,
		})
		m := transport.OutgoingMessage{
			RoutingKey: "broadcast.all",
			Data:       d,
		}
		n.transport.Send() <- &m
	}
	sendp()

	for n.state == Leader {
		select {
		case p := <-pc:
			fmt.Println(n.id.String()+" | Leader: Received packet: ", p)

			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				fmt.Println("couldn't decode packet", err)
				break
			}

			switch m.(type) {
			case *protocol.RequestVoteCall:
				break
			case *protocol.RequestVoteRespond:
				break
			case *protocol.AppendEntriesCall:
				break
			case *protocol.AppendEntriesRespond:
				break
			}
			break
		case <-hb.C:
			fmt.Println(n.id.String() + " | Leader: Heartbeat!")
			sendp()

			break
		}
	}
}

func (n *Node) generateTimeout() time.Duration {
	timeout := util.RandInt(
		defaultElectionTimeout,
		defaultElectionTimeout*2,
		time.Now().UnixNano(),
		n.transport.State().Id().String(),
	)
	return time.Duration(timeout) * time.Millisecond
}

func (n *Node) processVoteRequest(req protocol.RequestVoteCall) bool {
	if req.Term < n.currentTerm {
		return false
	}
	if n.votedFor != uuid.Nil {
		return n.votedFor.String() == req.CandidateId
	}
	lle := n.log.LastLogEntry()
	return req.LastLogIndex >= lle.Index && req.LastLogTerm >= lle.Term
}
