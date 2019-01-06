package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/internal/artemis"
	"github.com/superioz/artemis/internal/clirest"
	"github.com/superioz/artemis/pkg/consoleprint"
	"github.com/superioz/artemis/raft"
	"sync"
	"time"
)

func main() {

	fmt.Println("Hello, World! We need to load the config (ᴗ˳ᴗ)")
	group := sync.WaitGroup{}

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
	irest := clirest.Startup(cfg.Rest, &group)
	group.Wait()

	logrus.WithFields(logrus.Fields{
		"address": irest.Address(),
	}).Infoln("internal rest server started")

	// ---------
	// raft node
	// ---------

	timeStamp := time.Now()

	logrus.Infoln("starting node..")
	node := raft.NewNode(cfg)

	group.Add(1)
	go node.Up(fmt.Sprintf("%s:%d", cfg.Broker.Host, cfg.Broker.Port), &group)
	group.Wait()

	logrus.WithFields(logrus.Fields{
		"runtime": time.Now().Sub(timeStamp),
	}).Infoln("node started.")

	// ------
	// finish
	// ------

	artemis.New(&cfg, &node, irest)

	// --------
	// run loop
	// --------

	for node.BrokerConnected() {
		// wait for node to finish
		select {}
	}
	logrus.WithFields(logrus.Fields{
		"runtime": time.Now().Sub(timeStamp),
	}).Errorln("program exited.")
}
