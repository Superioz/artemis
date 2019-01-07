package raft

import (
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/pkg/uid"
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

	waitGroup := sync.WaitGroup{}
	node := NewNode(cfg, uid.NewUID())
	go node.Up("amqp://guest:guest@localhost:5672", &waitGroup)

	res := false
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

	node2 := NewNode(cfg, uid.NewUID())
	node2.Passive = true
	go node2.Up("amqp://guest:guest@localhost:5672", &waitGroup)

	waitGroup.Wait()
	if !res {
		t.Fatal("expected node1 to be leader but is not")
	}
}
