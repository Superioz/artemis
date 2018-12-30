package raft

import (
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/superioz/artemis/pkg/logger"
	"github.com/superioz/artemis/pkg/transport"
	"github.com/superioz/artemis/pkg/util"
	"time"
)

type State string

const (
	Follower  = "follower"
	Candidate = "candidate"
	Leader    = "leader"

	exchange = "artemis"

	defaultHeartbeat       = 250 // in ms, def: 50
	defaultElectionTimeout = 750 // in ms, def: 150
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

func (n *Node) generateTimeout() time.Duration {
	timeout := util.RandInt(
		defaultElectionTimeout,
		defaultElectionTimeout*2,
		time.Now().UnixNano(),
		n.transport.State().Id().String(),
	)
	return time.Duration(timeout) * time.Millisecond
}

func (n *Node) Up(brokerUrl string) {
	err := n.transport.Connect(brokerUrl)
	if err != nil || !n.transport.State().Connected() {
		logger.Err(fmt.Sprintf("couldn't connect to broker %s", brokerUrl), err)
		return
	}

	// TODO currently only follower loop
mainLoop:
	for n.transport.State().Connected() {
		timeout := n.generateTimeout()
		tc := time.After(timeout)
		pc := n.transport.Receive()

		for n.state == Follower {
			select {
			case p := <-pc:
				fmt.Println("Received packet: ", p)

				break
			case <-tc:
				fmt.Printf("Timeout! (After %s)\n", timeout)

				// TODO set state to candidate and send request vote rpc's
				continue mainLoop
			}
		}
	}
}

func (n *Node) Down() {
	_ = n.transport.Disconnect()
}
