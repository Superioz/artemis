package rest

import (
	"github.com/superioz/artemis/config"
	"github.com/valyala/fasthttp"
	"sync"
	"testing"
)

// makes sure, that the rest server starts properly
func TestNew(t *testing.T) {
	handle := func(ctx *fasthttp.RequestCtx) {
		_, _ = ctx.Write([]byte("Hello World!"))
	}

	server := New(config.DefaultCLIRestConfig())
	server.router.GET("/", handle)

	group := sync.WaitGroup{}
	group.Add(1)
	err := server.Up(&group)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case e := <-server.ErrListen():
		t.Fatal(e)
	default:
		break
	}
}
