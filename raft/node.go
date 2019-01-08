package raft

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/pkg/transport"
	"github.com/superioz/artemis/pkg/uid"
	"github.com/superioz/artemis/pkg/util"
	"github.com/superioz/artemis/raft/protocol"
	"sync"
	"time"
)

type State string

const (
	Follower  = "follower"
	Candidate = "candidate"
	Leader    = "leader"
)

type Node struct {
	id    uid.UID
	state State

	leader uuid.UUID

	currentTerm uint64  // stable storage
	votedFor    uid.UID // stable storage
	log         *Log    // stable storage

	currentVotes map[uid.UID]bool   // volatile for candidates
	commitIndex  uint64             // volatile
	lastApplied  uint64             // volatile
	nextIndex    map[uid.UID]uint64 // volatile for leaders
	matchIndex   map[uid.UID]uint64 // volatile for leaders

	transport transport.Interface

	heartbeatInterval time.Duration
	electionTimeout   time.Duration

	Passive bool

	config config.NodeConfig
}

func NewNode(config config.NodeConfig, id uid.UID) Node {
	trans := transport.NewAMQPInterface(config.Broker.ExchangeKey, id)

	n := Node{
		id:                id,
		state:             Follower,
		currentTerm:       0,
		log:               &Log{},
		currentVotes:      make(map[uid.UID]bool),
		commitIndex:       0,
		lastApplied:       0,
		nextIndex:         make(map[uid.UID]uint64),
		matchIndex:        make(map[uid.UID]uint64),
		transport:         &trans,
		heartbeatInterval: time.Duration(config.HeartbeatInterval),
		config:            config,
	}
	return n
}

func (n *Node) Log() *Log {
	return n.log
}

func (n *Node) CurrentTerm() uint64 {
	return n.currentTerm
}

func (n *Node) State() State {
	return n.state
}

func (n *Node) TransportState() *transport.State {
	return n.transport.State()
}

func (n *Node) SetState(state State) {
	n.state = state

	// call event
	logrus.WithFields(logrus.Fields{
		"state": state,
		"id":    n.id,
	}).Infoln("node changed state")
	Fire(ChangeStateEvent, *n)
}

func (n *Node) QuorumSize() int {
	return (n.config.ClusterSize / 2) + 1
}

func (n *Node) BrokerConnected() bool {
	return n.transport.State().Connected()
}

func (n *Node) Up(brokerUrl string, group *sync.WaitGroup) {
	err := n.transport.Connect(brokerUrl)

	if err != nil || !n.transport.State().Connected() {
		logrus.Fatalln(fmt.Sprintf("couldn't connect to broker %s", brokerUrl), err)
	}
	group.Done()

	// fire event
	Fire(StartupEvent, *n)

	// main loop
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

	// fire event
	Fire(ShutdownEvent, *n)
}

func (n *Node) generateTimeout() time.Duration {
	timeout := util.RandInt(
		n.config.ElectionTimeout,
		n.config.ElectionTimeout*2,
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
		RoutingKey: n.config.Broker.BroadcastRoute,
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
		RoutingKey: n.config.Broker.BroadcastRoute,
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
