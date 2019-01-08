package raft

import (
	"fmt"
	"github.com/superioz/artemis/pkg/util"
	"github.com/superioz/artemis/raft/protocol"
	"math"
)

type Log struct {
	entries []*protocol.LogEntry
}

type LogState struct {
	LastIndex uint64
	LastTerm  uint64
	Size      int
}

// creates a new log instance
func NewLog(entries []*protocol.LogEntry) Log {
	return Log{entries: entries}
}

// returns the current log entries
func (l *Log) Entries() []*protocol.LogEntry {
	return l.entries
}

// gets the entry of the log slice
func (l *Log) GetEntry(index int) *protocol.LogEntry {
	if index < 0 || index >= len(l.entries) {
		return nil
	}
	return l.entries[index]
}

// get the last x entries of the log
func (l *Log) LastEntries(size int) ([]*protocol.LogEntry, error) {
	entries, err := l.SubSlice(int(math.Max(float64(len(l.entries)-size), 0)), int(float64(len(l.entries)-1)))
	return entries, err
}

// returns a sub slice of the log entry
// from and to including
func (l *Log) SubSlice(from int, to int) ([]*protocol.LogEntry, error) {
	if from < 0 || from >= len(l.entries) || to < 0 || to >= len(l.entries) {
		return nil, fmt.Errorf("index exceeds slice size")
	}

	var slice []*protocol.LogEntry
	for i := from; i <= to; i++ {
		slice = append(slice, l.entries[i])
	}
	return slice, nil
}

// checks if an entry at given index with given term
// exists.
func (l *Log) HasEntry(logIndex uint64, logTerm uint64) bool {
	if int(logIndex) >= len(l.entries) {
		return false
	}
	e := l.entries[logIndex]
	return e != nil && e.Term == logTerm
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
	if i != 0 && !(i < 0 || i >= uint64(len(l.entries))) {
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
	l.entries = append(l.entries, entry)
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
