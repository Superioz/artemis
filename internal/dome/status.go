package dome

import (
	"github.com/superioz/artemis/pkg/uid"
	"github.com/superioz/artemis/raft/protocol"
	"time"
)

type Status struct {
	BrokerConnected bool      `json:"brokerConnected"`
	LastSent        time.Time `json:"lastSent"`
	LastReceived    time.Time `json:"lastReceived"`
	Id              uid.UID   `json:"id"`
	ClusterSize     int       `json:"clusterSize"`
	Runtime         time.Time `json:"runtime"`
	Ping            time.Time `json:"ping"`
	Raft            struct {
		State   string              `json:"state"`
		Term    string              `json:"term"`
		Log     []protocol.LogEntry `json:"logview"`
		Timeout time.Time           `json:"timeout"`
	} `json:"raft"`
}

func CurrentStatus() Status {
	// TODO
	return Status{}
}
