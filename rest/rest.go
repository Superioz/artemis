package rest

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func lock() {
	fasthttp.AcquireArgs()
	fasthttprouter.New()
}
