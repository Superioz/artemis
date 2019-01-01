package raft

import (
	"fmt"
	"github.com/superioz/artemis/pkg/logc"
	"testing"
)

// makes sure that a node can be elected as
// leader. In this test we focus on the election timeout
// and the candidacy leading to one being a leader in the
// cluster.
func TestNodeElection(t *testing.T) {
	logc.ApplyConfig(logc.DefaultConfig)
	node := NewNode()
	go node.Up("amqp://guest:guest@localhost:5672")

	node2 := NewNode()
	node2.Passive = true
	fmt.Println("Passive is: " + node2.id.String())
	go node2.Up("amqp://guest:guest@localhost:5672")

	for {
		select {}
	}
}
