package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/internal/dome"
)

func main() {

	fmt.Println("Hello, World! We need to build the dome (ᴗ˳ᴗ)")

	d := dome.Startup()

	// --------
	// run loop
	// --------

	for d.Node().BrokerConnected() {
		// wait for node to finish
		select {}
	}
	logrus.WithFields(logrus.Fields{
		"runtime": d.GetRuntime(),
	}).Errorln("program exited.")
}
