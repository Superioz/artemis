package main

import (
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config/logc"
)

func main() {
	logc.ApplyConfig(logc.DefaultConfig)
	logrus.Info("Hello, World!")
}
