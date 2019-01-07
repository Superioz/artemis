package raft

import (
	"fmt"
	"github.com/superioz/artemis/raft/protocol"
	"testing"
)

func TestLog_AppendEntries(t *testing.T) {
	log := Log{
		entries: []*protocol.LogEntry{},
	}
	err := log.AppendEntries(&protocol.LogEntry{
		Term: 1,
	}, &protocol.LogEntry{
		Term: 2,
	})

	if err != nil {
		t.Fatal(err)
	}
	if len(log.entries) == 0 {
		t.Fatal("couldn't append to log")
	}

	entry := log.GetEntry(0)
	if entry == nil {
		t.Fatal("couldn't find log entry index=0")
	}
	if entry.Term != 1 {
		t.Fatal("stored entry is not expected value")
	}
}

func TestLog_LastEntries(t *testing.T) {
	log := Log{
		entries: []*protocol.LogEntry{},
	}
	err := log.AppendEntries(&protocol.LogEntry{
		Term: 1,
	}, &protocol.LogEntry{
		Term: 2,
	}, &protocol.LogEntry{
		Term: 3,
	}, &protocol.LogEntry{
		Term: 4,
	}, &protocol.LogEntry{
		Term: 5,
	}, )

	if err != nil {
		t.Fatal(err)
	}
	if len(log.entries) == 0 {
		t.Fatal("couldn't append to log")
	}

	lastEntries, err := log.LastEntries(3)
	if err != nil {
		t.Fatal(err)
	}
	if len(lastEntries) != 3 {
		t.Fatal(fmt.Sprintf("unexpected size of slice, %d", len(lastEntries)))
	}

	entry := lastEntries[0]
	if entry.Term != 3 {
		t.Fatal(fmt.Sprintf("stored entry is not expected value, %d", entry.Term))
	}

	lastEntries, err = log.LastEntries(10)
	if err != nil {
		t.Fatal(err)
	}
	if len(lastEntries) != len(log.entries) {
		t.Fatal(fmt.Sprintf("unexpected slice size %d", len(lastEntries)))
	}
}
