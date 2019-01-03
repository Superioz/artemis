package raft

import (
	"github.com/superioz/artemis/config"
	"sync"
	"testing"
)

// makes sure that a node can be elected as
// leader. In this test we focus on the election timeout
// and the candidacy leading to one being a leader in the
// cluster.
func TestNodeElection(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Fatal(err)
	}

	node := NewNode(cfg)
	go node.Up("amqp://guest:guest@localhost:5672")

	res := false
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)
	AddHook(func(e *Event) error {
		if e.Node.id != node.id || e.Type != ChangeStateEvent {
			return nil
		}
		if e.Node.state != Candidate {
			res = e.Node.state == Leader
			waitGroup.Done()
		}
		return nil
	})

	node2 := NewNode()
	node2.Passive = true
	go node2.Up("amqp://guest:guest@localhost:5672")

	waitGroup.Wait()
	if !res {
		t.Fatal("expected node1 to be leader but is not")
	}
}
