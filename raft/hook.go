package raft

import (
	"github.com/sirupsen/logrus"
	"time"
)

const (
	StartupEvent              = "startup"
	ShutdownEvent             = "shutdown"
	ChangeStateEvent          = "change_state"
)

type EventType string

// the event information to be passed to
// the event hook
type Event struct {
	Type      EventType
	Node      Node
	Timestamp time.Time
}

// internal type for storing the hooks on a node instance
type EventHooks []func(e *Event) error

// Add a hook to an instance of logger. This is called with
// `log.Hooks.Add(new(MyHook))` where `MyHook` implements the `Hook` interface.
func (hooks *EventHooks) Add(hook func(e *Event) error) {
	*hooks = append(*hooks, hook)
}

// Fire all the hooks for the passed level. Used by `entry.log` to fire
// appropriate hooks for a log entry.
func (hooks EventHooks) Fire(t EventType, node Node) error {
	for _, hook := range hooks {
		if err := hook(&Event{Type: t, Node: node, Timestamp: time.Now()}); err != nil {
			return err
		}
	}

	return nil
}

var eventHooks = EventHooks{}

func AddHook(hook func(e *Event) error) {
	eventHooks.Add(hook)
}

func Fire(t EventType, node Node) {
	err := eventHooks.Fire(t, node)
	if err != nil {
		logrus.Errorln("could not fire event", t, node.id, err)
	}
}
