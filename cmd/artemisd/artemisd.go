package main

import (
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/raft"
)

func main() {
	logrus.Info("Hello, World!")

	cfg, err := config.Load()
	if err != nil {
		logrus.Fatalln("couldn't load config file :(")
	}

	node := raft.NewNode(cfg)
	go node.Up(cfg.Broker.Host + ":" + cfg.Broker.Port)

	for node.BrokerConnected() {
		// wait for node to finish
		select {}
	}
}
