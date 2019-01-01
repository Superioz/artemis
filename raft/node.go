package raft

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/pkg/transport"
	"github.com/superioz/artemis/pkg/uid"
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
	id    uid.UID
	state State

	leader uuid.UUID

	currentTerm uint64  // stable storage
	votedFor    uid.UID // stable storage
	log         Log     // stable storage

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
		logrus.Errorln(fmt.Sprintf("couldn't connect to broker %s", brokerUrl), err)
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
	n.votedFor = uid.UID{}

followerLoop:
	for n.state == Follower {
		select {
		case p := <-pc:
			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				logrus.Errorln("couldn't decode packet", err)
				break
			}

			switch m.(type) {
			case *protocol.RequestVoteCall:
				// * check if he can grant the vote and respond with
				// his answer.

				reqVote := *m.(*protocol.RequestVoteCall)

				logrus.Debugln("follower.vote.req.inc", n.id, m)

				n.processRequestVote(reqVote, p.Source.String())
				break
			case *protocol.RequestVoteRespond:
				// * ignore as followers can't process this packet. only
				// candidates can.
				break
			case *protocol.AppendEntriesCall:
				// * the leader sends an append entries call, so
				// update the leader and reset the timeout.

				appendEntr := *m.(*protocol.AppendEntriesCall)
				logrus.Debugln("follower.appendEntries", n.id, appendEntr)

				// reset timeout
				timeout = n.generateTimeout()
				tc.Reset(timeout)

				// append entries
				err := n.processAppendEntries(appendEntr, p.Source.String())
				n.sendAppendEntriesResponse(err == nil, p.Source.String())
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
			logrus.Debugln("follower.timeout", n.id, timeout)

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
	logrus.Debugln("candidate.vote.req", n.id)
	n.sendRequestVote()

candidateLoop:
	for n.state == Candidate {
		select {
		case p := <-packetChan:
			m, err := transport.Decode(p.Packet.Data)
			if err != nil {
				logrus.Errorln("couldn't decode packet", err)
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

				n.processRequestVote(reqVote, p.Source.String())
				break
			case *protocol.RequestVoteRespond:
				// TODO count negative and positive responds and calculate if he got the majority ..

				n.state = Leader
				logrus.Debugln("candidate.vote.inc", n.id)
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

			logrus.Debugln("candidate.timeout", n.id, timeout)
			break candidateLoop
		case <-hardBeetTimer.C:
			// * try again to receive votes from followers, cause maybe not every follower
			// received the packet or responded yet.

			logrus.Debugln("candidate.heartbeat", n.id)
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
				logrus.Errorln("couldn't decode packet", err)
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

				appendEntrResp := *m.(*protocol.AppendEntriesRespond)

				logrus.Debugln("leader.appendEntries.respond", n.id, appendEntrResp)

				break
			}
			break
		case <-hb.C:
			// * send normal heartbeat with no information
			// just to keep his authority.
			logrus.Debugln("leader.heartbeat", n.id)
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

func (n *Node) processRequestVote(req protocol.RequestVoteCall, source string) {
	var res bool

	if n.votedFor != uid.Nil {
		res = n.votedFor.String() == req.CandidateId
	} else if req.Term < n.currentTerm {
		res = false
	} else {
		ls := n.log.MirrorState()
		res = req.LastLogIndex >= ls.LastIndex && req.LastLogTerm >= ls.LastTerm
	}

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

func (n *Node) processAppendEntries(appendEntr protocol.AppendEntriesCall, source string) error {
	// update term if
	if appendEntr.Term > n.currentTerm {
		n.currentTerm = appendEntr.Term
	}

	// set leader
	n.leader, _ = uuid.FromString(appendEntr.LeaderId)

	// check for prevLogIndex
	if n.log.HasEntry(appendEntr.PrevLogIndex, appendEntr.PrevLogTerm) {
		err := fmt.Errorf("log does not contain prev log context")
		logrus.Errorln(err)
		return err
	}

	// apply entries
	if appendEntr.CommitIndex > n.lastApplied {
		err := n.log.ApplyEntries(n.lastApplied, appendEntr.CommitIndex)
		if err != nil {
			logrus.Errorln("error while applying entries", err)
			return err
		} else {
			n.lastApplied = appendEntr.CommitIndex
		}
	}

	// append entries
	if len(appendEntr.Entries) > 0 {
		err := n.log.AppendEntries(appendEntr.Entries...)
		if err != nil {
			logrus.Errorln("error while appending entries", err)
			return err
		}
	}

	return nil
}

func (n *Node) sendAppendEntriesResponse(success bool, source string) {
	d, _ := transport.Encode(&protocol.AppendEntriesRespond{
		Term:    n.currentTerm,
		Success: success,
	})
	m := transport.OutgoingMessage{
		RoutingKey: source,
		Data:       d,
	}
	n.transport.Send() <- &m
}
