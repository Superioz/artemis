package dome

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/internal/clirest"
	"github.com/superioz/artemis/pkg/consoleprint"
	"github.com/superioz/artemis/pkg/rest"
	"github.com/superioz/artemis/pkg/uid"
	"github.com/superioz/artemis/raft"
	"sync"
	"time"
)

var instance *Dome

type Dome struct {
	id uid.UID

	config  *config.NodeConfig
	raft    *raft.Node
	clirest *rest.Server

	started time.Time
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

func Instance() *Dome {
	return instance
}

func Startup() Dome {
	instance = &Dome{
		id: uid.NewUID(),
	}
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
	node := raft.NewNode(cfg, instance.id)

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

	instance.config = &cfg
	instance.clirest = irest
	instance.raft = &node
	instance.started = time.Now()

	logrus.Infoln("artemis has been set up")
	return *instance
}
