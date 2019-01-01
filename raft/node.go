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

	exchange       = "artemis"
	broadcastRoute = "broadcast.all"

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

	Passive bool
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
			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				logger.Err("couldn't decode packet", err)
				break
			}

			switch m.(type) {
			case *protocol.RequestVoteCall:
				// * check if he can grant the vote and respond with
				// his answer.

				reqVote := *m.(*protocol.RequestVoteCall)

				logger.Debug("follower.vote.req.inc", n.id, m)

				n.responseRequestVote(reqVote, p.Source.String())
				break
			case *protocol.RequestVoteRespond:
				// * ignore as followers can't process this packet. only
				// candidates can.
				break
			case *protocol.AppendEntriesCall:
				// * the leader sends an append entries call, so
				// update the leader and reset the timeout.

				appendEntr := *m.(*protocol.AppendEntriesCall)

				// update term if
				if appendEntr.Term > n.currentTerm {
					n.currentTerm = appendEntr.Term
				}

				// reset timeout
				timeout = n.generateTimeout()
				tc.Reset(timeout)

				// set leader
				n.leader, _ = uuid.FromString(appendEntr.LeaderId)

				// apply entries
				if appendEntr.CommitIndex > n.lastApplied {
					err := n.log.ApplyEntries(n.lastApplied, appendEntr.CommitIndex)
					if err != nil {
						logger.Err("error while applying entries", err)
						break
					}

					n.lastApplied = appendEntr.CommitIndex
				}

				// append entries
				if len(appendEntr.Entries) > 0 {
					err := n.log.AppendEntries(appendEntr.Entries...)
					if err != nil {
						logger.Err("error while appending entries", err)
						break
					}
				}
				break
			case *protocol.AppendEntriesRespond:
				// * ignore as followers can't process this packet. only
				// the leader can.
				break
			}
			break
		case <-tc.C:
			// * timeout and try to become leader by sending request votes
			// also step up to being candidate
			logger.Debug("follower.timeout", n.id, timeout)

			// if the node is only passive, don't try to ever get leader
			if n.Passive {
				break followerLoop
			}
			n.state = Candidate
			break
		}
	}
}

func (n *Node) candidateLoop() {
	timeout := n.generateTimeout()

	timeoutTimer := time.NewTimer(timeout)
	packetChan := n.transport.Receive()
	hardBeetTimer := time.NewTicker(n.heartbeatInterval * time.Millisecond)

	// increment term and vote for himself
	n.currentTerm++
	n.votedFor = n.id

	// send request vote packet function
	logger.Debug("candidate.vote.req", n.id)
	n.sendRequestVote()

candidateLoop:
	for n.state == Candidate {
		select {
		case p := <-packetChan:
			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				logger.Err("couldn't decode packet", err)
				break
			}

			switch m.(type) {
			case *protocol.RequestVoteCall:
				// * sent back false, as we already voted for ourself

				reqVote := *m.(*protocol.RequestVoteCall)

				// if term of packet is higher than current term
				// update and step back
				if reqVote.Term > n.currentTerm {
					n.currentTerm = reqVote.Term
					n.state = Follower
					break candidateLoop
				}

				n.responseRequestVote(reqVote, p.Source.String())
				break
			case *protocol.RequestVoteRespond:
				// TODO count negative and positive responds and calculate if he got the majority ..

				n.state = Leader
				logger.Debug("leader.vote.inc", n.id)
				break candidateLoop
			case *protocol.AppendEntriesCall:
				// * step back from being candidate, as there is already a leader
				// sending append entries.
				// * also reset timeout

				appendEntr := *m.(*protocol.AppendEntriesCall)

				// reset timeout
				timeout = n.generateTimeout()
				timeoutTimer.Reset(timeout)

				// set leader
				n.leader, _ = uuid.FromString(appendEntr.LeaderId)

				n.state = Follower
				break candidateLoop
			case *protocol.AppendEntriesRespond:
				// * ignore packet, as we are a candidate and can't
				// process these packets. Only the leader can.
				break
			}
			break
		case <-timeoutTimer.C:
			// * begin a new term and try again receiving votes.

			logger.Debug("candidate.timeout", n.id, timeout)
			break candidateLoop
		case <-hardBeetTimer.C:
			// * try again to receive votes from followers, cause maybe not every follower
			// received the packet or responded yet.

			logger.Debug("candidate.heartbeat", n.id)
			n.sendRequestVote()
			break
		}
	}

	timeoutTimer.Stop()
	hardBeetTimer.Stop()
}

func (n *Node) leaderLoop() {
	pc := n.transport.Receive()
	hb := time.NewTicker(n.heartbeatInterval * time.Millisecond)

	// send packet
	n.sendHeartbeat()

	for n.state == Leader {
		select {
		case p := <-pc:
			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				logger.Err("couldn't decode packet", err)
				break
			}

			switch m.(type) {
			case *protocol.RequestVoteCall:
				// * that can't possibly happen. But if so just ignore until
				// another heartbeat gets sent
				break
			case *protocol.RequestVoteRespond:
				// * we didn't call for a vote so just ignore this
				// packet.
				break
			case *protocol.AppendEntriesCall:
				// * a client redirected an append entries call.

				break
			case *protocol.AppendEntriesRespond:
				// * a client responded with the result of the append entries
				// call.

				break
			}
			break
		case <-hb.C:
			// * send normal heartbeat with no information
			// just to keep his authority.
			logger.Debug("leader.heartbeat", n.id)
			n.sendHeartbeat()
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

func (n *Node) sendRequestVote() {
	ls := n.log.MirrorState()
	d, _ := transport.Encode(&protocol.RequestVoteCall{
		Term:         n.currentTerm,
		CandidateId:  n.id.String(),
		LastLogIndex: ls.LastIndex,
		LastLogTerm:  ls.LastTerm,
	})
	m := transport.OutgoingMessage{
		RoutingKey: broadcastRoute,
		Data:       d,
	}
	n.transport.Send() <- &m
}

func (n *Node) responseRequestVote(req protocol.RequestVoteCall, source string) {
	var res bool

	if n.votedFor != uuid.Nil {
		res = n.votedFor.String() == req.CandidateId
	} else if req.Term < n.currentTerm {
		res = false
	}
	ls := n.log.MirrorState()
	res = req.LastLogIndex >= ls.LastIndex && req.LastLogTerm >= ls.LastTerm

	d, _ := transport.Encode(&protocol.RequestVoteRespond{
		Term:        n.currentTerm,
		VoteGranted: res,
	})
	om := transport.OutgoingMessage{
		RoutingKey: source,
		Data:       d,
	}
	n.transport.Send() <- &om
}

func (n *Node) sendHeartbeat() {
	logState := n.log.MirrorState()
	d, _ := transport.Encode(&protocol.AppendEntriesCall{
		Term:         n.currentTerm,
		LeaderId:     n.id.String(),
		PrevLogIndex: logState.LastIndex,
		PrevLogTerm:  logState.LastTerm,
		Entries:      []*protocol.LogEntry{},
		CommitIndex:  n.commitIndex,
	})
	m := transport.OutgoingMessage{
		RoutingKey: broadcastRoute,
		Data:       d,
	}
	n.transport.Send() <- &m
}
