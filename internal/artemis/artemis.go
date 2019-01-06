package artemis

import (
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/pkg/rest"
	"github.com/superioz/artemis/raft"
)

var bubble *Bubble

type Bubble struct {
	config  *config.NodeConfig
	raft    *raft.Node
	clirest *rest.Server
}

func New(config *config.NodeConfig, raft *raft.Node, clirest *rest.Server) *Bubble {
	bubble = &Bubble{
		config:  config,
		raft:    raft,
		clirest: clirest,
	}

	logrus.Infoln("artemis has been set up")
	return bubble
}
