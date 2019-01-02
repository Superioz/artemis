package raft

import (
	"testing"
)

// makes sure that hooks get registered
// fired and unregistered
func TestHook(t *testing.T) {
	res := false
	h := func(e *Event) error {
		res = true
		return nil
	}
	AddHook(h)
	Fire(EventType("test"), Node{})
	if !res {
		t.Fatal("hook did not get fired")
	}
}
