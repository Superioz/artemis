package raft

import (
	"github.com/superioz/artemis/pkg/util"
	"github.com/superioz/artemis/raft/protocol"
)

type Log struct {
	entries []*protocol.LogEntry
}

type LogState struct {
	LastIndex uint64
	LastTerm  uint64
	Size      int
}

// mirrors the current state of the log
// with replicating the information about
// the state.
func (l *Log) MirrorState() LogState {
	var lastEntry protocol.LogEntry
	if len(l.entries) == 0 {
		lastEntry = protocol.LogEntry{}
	} else {
		lastEntry = *l.entries[len(l.entries)-1]
	}

	return LogState{
		LastIndex: lastEntry.Index,
		LastTerm:  lastEntry.Term,
		Size:      len(l.entries),
	}
}

// appends given entry to the log.
// if the index of the given entry is `0`, we append
// the entry to the last index. Otherwise we try to set
// the entry to given index. Normally, a leader calls for an
// `AppendEntries` with specific index values.
func (l *Log) AppendEntry(entry *protocol.LogEntry) error {
	i := entry.Index
	t := entry.Term

	// if index is set to `0` just append the entry. otherwise
	// we try to set the entry to specific index.
	if i != 0 {
		// if there is already an entry at this index..
		cont := l.entries[i]

		// ..but with a different term..
		if cont != nil && cont.Term != t {
			// ..delete everything from this index on..
			sl, err := util.Delete(l.entries, int(i), len(l.entries))
			if err != nil {
				return err
			}
			l.entries = sl.([]*protocol.LogEntry)

			// ..and set the index to the new entry
			l.entries[i] = entry
		} else if cont == nil {
			// just set the entry to the index
			l.entries[i] = entry
		}
		return nil
	}

	// append entry to log
	index := uint64(len(l.entries))
	entry.Index = index
	l.entries[index] = entry
	return nil
}

// appends many entries to the log.
// successful if error != nil
func (l *Log) AppendEntries(e ...*protocol.LogEntry) error {
	for _, entry := range e {
		err := l.AppendEntry(entry)
		if err != nil {
			return err
		}
	}
	return nil
}

// applies the entries up until given index.
func (l *Log) ApplyEntries(indexf uint64, indext uint64) error {
	for i := indexf; i <= indext; i++ {

	}
	return nil
}
