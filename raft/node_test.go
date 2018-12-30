package raft

import "testing"

// makes sure that a node can be elected as
// leader. In this test we focus on the election timeout
// and the candidacy leading to one being a leader in the
// cluster.
func TestNodeElection(t *testing.T) {
	node := NewNode()
	go node.Up("amqp://guest:guest@localhost:5672")

	for {
		select {}
	}
}
