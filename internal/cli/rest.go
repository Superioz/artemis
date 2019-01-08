package cli

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/pkg/rest"
	"sync"
)

// starts the internal rest server
func Startup(cfg config.CLIRest, group *sync.WaitGroup, hook func(router *fasthttprouter.Router)) *rest.Server {
	irest := rest.New(cfg)

	// add handler
	hook(irest.Router())

	group.Add(1)
	err := irest.Up(group)
	if err != nil {
		logrus.Fatalln("error while starting internal rest server", err)
	}
	group.Wait()

	// waiting for error to occur
	go func() {
		for {
			select {
			case err := <-irest.ErrListen():
				logrus.Fatalln("internal rest server threw an error", err)
			}
		}
	}()
	return irest
}
