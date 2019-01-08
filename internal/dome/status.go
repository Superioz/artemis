package dome

import (
	"github.com/superioz/artemis/pkg/uid"
	"github.com/superioz/artemis/raft/protocol"
	"time"
)

type Status struct {
	BrokerConnected bool                `json:"brokerConnected"`
	LastSent        time.Time           `json:"lastSent"`
	LastReceived    time.Time           `json:"lastReceived"`
	Id              uid.UID             `json:"id"`
	ClusterSize     int                 `json:"clusterSize"`
	Runtime         time.Duration       `json:"runtime"`
	Ping            time.Time           `json:"ping"`
	State           string              `json:"state"`
	Term            uint64              `json:"term"`
	Log             []protocol.LogEntry `json:"logview"`
}

func CurrentStatus() Status {
	d := Instance()
	lastEntries, _ := d.Node().Log().LastEntries(10)

	return Status{
		BrokerConnected: d.Node().BrokerConnected(),
		LastSent:        d.Node().TransportState().LastSent(),
		LastReceived:    d.Node().TransportState().LastReceived(),
		Id:              d.Id(),
		ClusterSize:     d.Config().ClusterSize,
		Runtime:         d.GetRuntime(),
		State:           string(d.Node().State()),
		Term:            d.Node().CurrentTerm(),
		Log:             lastEntries,
	}
}
