package main

import (
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/pkg/logc"
	"github.com/superioz/artemis/raft"
)

func main() {
	logc.ApplyConfig(logc.DefaultConfig)
	logrus.Info("Hello, World!")

	node := raft.NewNode()
	go node.Up("amqp://guest:guest@localhost:5672")

	for {
		select {}
	}
}
