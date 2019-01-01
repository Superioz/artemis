package raft

import (
	"testing"
)

type testEventHook struct {
	val bool
}

func (t *testEventHook) Fire(event *Event) error {
	if event.Type == "test" {
		t.val = !t.val
	}
	return nil
}

// makes sure that hooks get registered
// fired and unregistered
func TestHook(t *testing.T) {
	h := testEventHook{}
	AddHook(&h)
	Fire(EventType("test"), Node{})
	if !h.val {
		t.Fatal("hook did not get fired")
	}
}
