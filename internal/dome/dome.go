package dome

import (
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/pkg/rest"
	"github.com/superioz/artemis/pkg/uid"
	"github.com/superioz/artemis/raft"
	"time"
)

// represents an artemis daemon as a dome
// as the artemis setting is in space, dome seems
// fitting for containing all important stuff
type Dome struct {
	id uid.UID

	config  *config.NodeConfig
	raft    *raft.Node
	clirest *rest.Server

	started time.Time
}

// creates a new dome object
func Create(id uid.UID, cfg *config.NodeConfig, node *raft.Node, cliRest *rest.Server) *Dome {
	instance := Dome{
		id:      id,
		config:  cfg,
		clirest: cliRest,
		raft:    node,
		started: time.Now(),
	}
	return &instance
}

func (d *Dome) Node() *raft.Node {
	return d.raft
}

func (d *Dome) Config() *config.NodeConfig {
	return d.config
}

func (d *Dome) InternalRest() *rest.Server {
	return d.clirest
}

func (d *Dome) GetRuntime() time.Duration {
	return time.Now().Sub(d.started)
}

func (d *Dome) Id() uid.UID {
	return d.id
}
