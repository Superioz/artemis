package raft

import (
	"github.com/superioz/artemis/raft/protocol"
)

type Log []*protocol.LogEntry

// returns the log entry which has been
// added to the log recently. empty log entry
// if not found.
func (l Log) LastLogEntry() *protocol.LogEntry {
	if len(l) == 0 {
		return &protocol.LogEntry{}
	}
	return l[len(l)-1]
}
