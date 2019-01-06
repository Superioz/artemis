package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/internal/dome"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	fmt.Println("Hello, World! We need to build the dome (ᴗ˳ᴗ)")

	d := dome.Startup()

	// -------------
	// shutdown hook
	// -------------

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		logrus.WithFields(logrus.Fields{
			"runtime": d.GetRuntime(),
		}).Fatalln("program exited.")
	}()

	// --------
	// run loop
	// --------

	for d.Node().BrokerConnected() {
		// wait for node to finish
	}
	logrus.WithFields(logrus.Fields{
		"runtime": d.GetRuntime(),
	}).Fatalln("program exited.")
}
