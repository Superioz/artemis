package dome

import (
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/pkg/rest"
	"github.com/superioz/artemis/pkg/uid"
	"github.com/superioz/artemis/raft"
	"time"
)

type Dome struct {
	id uid.UID

	config  *config.NodeConfig
	raft    *raft.Node
	clirest *rest.Server

	started time.Time
}

func Create(id uid.UID, cfg *config.NodeConfig, node *raft.Node, cliRest *rest.Server) *Dome {
	instance := Dome{}
	instance.id = id
	instance.config = cfg
	instance.clirest = cliRest
	instance.raft = node
	instance.started = time.Now()
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
