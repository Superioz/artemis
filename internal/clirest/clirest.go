package clirest

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/pkg/rest"
	"github.com/valyala/fasthttp"
	"sync"
)

func Address() string {
	return fmt.Sprintf("http://%s:%d", config.DefaultRestHost, config.DefaultRestPort)
}

// starts the internal rest server
func Startup(cfg config.CLIRest, group *sync.WaitGroup) *rest.Server {
	irest := rest.New(cfg)

	// add handler
	irest.Router().GET("/", Index)
	irest.Router().GET("/status", Status)

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

func Get(header []byte, route string) ([]byte, int, error) {
	code, d, err := fasthttp.Get(header, fmt.Sprintf("%s/%s", Address(), route))
	return d, code, err
}

func Post(header []byte, route string, args *fasthttp.Args) ([]byte, int, error) {
	code, d, err := fasthttp.Post(header, fmt.Sprintf("%s/%s", Address(), route), args)
	return d, code, err
}
