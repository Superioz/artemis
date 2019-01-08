package cli

import (
	"github.com/superioz/artemis/cliexec"
	"github.com/superioz/artemis/config"
	"sync"
	"testing"
)

// makes sure the internal cli rest server starts
// properly
func TestStartup(t *testing.T) {
	group := sync.WaitGroup{}
	serv := Startup(config.DefaultCLIRestConfig(), &group, nil)

	if serv == nil {
		t.Fatal("server couldn't be created")
	}
}

// makes sure that get requests work
func TestGet(t *testing.T) {
	config.ApplyLoggingConfig(config.DefaultLoggingConfig())

	group := sync.WaitGroup{}
	serv := Startup(config.DefaultCLIRestConfig(), &group, nil)

	if serv == nil {
		t.Fatal("server couldn't be created")
	}

	data, _, _ := cliexec.Get([]byte{}, "/")
	if len(data) == 0 {
		t.Fatal("invalid data as response")
	}
}
