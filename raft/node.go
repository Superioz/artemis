package raft

import (
	"github.com/satori/go.uuid"
	"github.com/superioz/artemis/pkg/transport"
	"time"
)

type State string

const (
	None      = "none"
	Follower  = "follower"
	Candidate = "candidate"
	Leader    = "leader"
)

type Node struct {
	Id    uuid.UUID
	State State

	Leader uuid.UUID

	CurrentTerm uint64    // stable storage
	VotedFor    uuid.UUID // stable storage
	Log         Log       // stable storage

	CommitIndex uint64               // volatile
	LastApplied uint64               // volatile
	NextIndex   map[uuid.UUID]uint64 // volatile for leaders
	MatchIndex  map[uuid.UUID]uint64 // volatile for leaders

	Transport *transport.Interface

	HeartbeatInterval time.Duration
	ElectionTimeout   time.Duration
}

func (n *Node) RequestVote() {

}

func (n *Node) RequestAppendEntries() {

}
