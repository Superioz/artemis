package dome

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/internal/clirest"
	"github.com/superioz/artemis/pkg/consoleprint"
	"github.com/superioz/artemis/pkg/rest"
	"github.com/superioz/artemis/raft"
	"sync"
	"time"
)

var bubble *Dome

type Dome struct {
	config  *config.NodeConfig
	raft    *raft.Node
	clirest *rest.Server

	started time.Time
}

func (b *Dome) Node() *raft.Node {
	return b.raft
}

func (b *Dome) Config() *config.NodeConfig {
	return b.config
}

func (b *Dome) InternalRest() *rest.Server {
	return b.clirest
}

func (b *Dome) GetRuntime() time.Duration {
	return time.Now().Sub(b.started)
}

func newBubble(config *config.NodeConfig, raft *raft.Node, clirest *rest.Server) *Dome {
	bubble = &Dome{
		config:  config,
		raft:    raft,
		clirest: clirest,
		started: time.Now(),
	}

	logrus.Infoln("artemis has been set up")
	return bubble
}

func Startup() Dome {
	group := sync.WaitGroup{}
	timeStamp := time.Now()

	// ------
	// config
	// ------

	cfg, err := config.Load()

	// header
	consoleprint.PrintHeader("artemis")
	logrus.Infoln(" ")

	if err != nil {
		logrus.Errorln("couldn't load config file :( Falling back to default config..", err)
	} else {
		logrus.Infoln("config successfully loaded")
	}

	// -------------
	// internal rest
	// -------------

	logrus.Infoln("starting internal rest server..")
	irest := clirest.Startup(cfg.CLIRest, &group)
	group.Wait()

	logrus.WithFields(logrus.Fields{
		"address": irest.Address(),
	}).Infoln("internal rest server started")

	// ---------
	// raft node
	// ---------

	logrus.Infoln("starting raft..")
	node := raft.NewNode(cfg)

	group.Add(1)
	go node.Up(fmt.Sprintf("%s:%d", cfg.Broker.Host, cfg.Broker.Port), &group)
	group.Wait()

	logrus.Infoln("raft started.")

	// ------
	// finish
	// ------

	logrus.WithFields(logrus.Fields{
		"runtime": time.Now().Sub(timeStamp),
	}).Infoln("dome successfully built")
	return *newBubble(&cfg, &node, irest)
}
