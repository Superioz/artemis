package clirest

import (
	"github.com/superioz/artemis/config"
	"sync"
	"testing"
)

// makes sure the internal cli rest server starts
// properly
func TestStartup(t *testing.T) {
	group := sync.WaitGroup{}
	serv := Startup(config.DefaultCLIRestConfig(), &group)

	if serv == nil {
		t.Fatal("server couldn't be created")
	}
}

// makes sure that get requests work
func TestGet(t *testing.T) {
	config.ApplyLoggingConfig(config.DefaultLoggingConfig())

	group := sync.WaitGroup{}
	serv := Startup(config.DefaultCLIRestConfig(), &group)

	if serv == nil {
		t.Fatal("server couldn't be created")
	}

	data, _, _ := Get([]byte{}, "/", "localhost:2310")
	if len(data) == 0 {
		t.Fatal("invalid data as response")
	}
}
