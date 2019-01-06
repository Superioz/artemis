package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/pkg/consoleprint"
	"github.com/superioz/artemis/pkg/rest"
	"github.com/superioz/artemis/raft"
	"sync"
	"time"
)

func main() {
	fmt.Println("Hello, World!")

	// ------
	// config
	// ------

	cfg, err := config.Load()
	if err != nil {
		logrus.Errorln("couldn't load config file :(", err)
	}
	consoleprint.PrintHeader("artemis")
	logrus.Infoln(" ")

	// -------------
	// internal rest TODO
	// -------------

	irest := rest.New(cfg.Rest)
	err = irest.Up()
	if err != nil {
		logrus.Fatalln("error while starting internal rest server", err)
	}

	// ---------
	// raft node
	// ---------

	timeStamp := time.Now()
	group := sync.WaitGroup{}

	logrus.Infoln("starting node..")
	node := raft.NewNode(cfg)

	group.Add(1)
	go node.Up(fmt.Sprintf("%s:%d", cfg.Broker.Host, cfg.Broker.Port), &group)
	group.Wait()

	logrus.WithFields(logrus.Fields{
		"runtime": time.Now().Sub(timeStamp),
	}).Infoln("node started.")

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
