package dome

import (
	"fmt"
	"github.com/superioz/artemis/appversion"
	"github.com/superioz/artemis/pkg/uid"
	"github.com/superioz/artemis/raft/protocol"
	"time"
)

// represents a current status of the dome
// containing all important information.
// can be used to e.g. send to the cli.
type Status struct {
	Version         string              `json:"version"`
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

// fetches the current status. simple as that.
func CurrentStatus() Status {
	d := Instance()
	lastEntries, _ := d.Node().Log().LastEntries(10)

	return Status{
		Version:         fmt.Sprintf("%s, build %s", appversion.Version, appversion.Build),
		BrokerConnected: d.Node().BrokerConnected(),
		LastSent:        d.Node().TransportState().LastSent(),
		LastReceived:    d.Node().TransportState().LastReceived(),
		Id:              d.Id(),
		ClusterSize:     d.Config().ClusterSize,
		Runtime:         d.GetRuntime(),
		Ping:            time.Now(),
		State:           string(d.Node().State()),
		Term:            d.Node().CurrentTerm(),
		Log:             lastEntries,
	}
}
