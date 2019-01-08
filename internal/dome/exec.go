package dome

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/internal/cli"
	"github.com/superioz/artemis/pkg/consoleprint"
	"github.com/superioz/artemis/pkg/uid"
	"github.com/superioz/artemis/raft"
	"sync"
	"time"
)

var instance *Dome

func Instance() *Dome {
	return instance
}

// builds the dome with the important features such as
// the cli rest server, the broker connection, ...
//
// returns the built dome.
func Build(hook func(router *fasthttprouter.Router)) *Dome {
	id := uid.NewUID()
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
	irest := cli.Startup(cfg.CLIRest, &group, hook)
	group.Wait()

	logrus.WithFields(logrus.Fields{
		"address": irest.Address(),
	}).Infoln("internal rest server started")

	// ---------
	// raft node
	// ---------

	logrus.Infoln("starting raft..")
	node := raft.NewNode(cfg, id)

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

	instance = Create(id, &cfg, &node, irest)

	logrus.Infoln("artemis has been set up")
	return instance
}
