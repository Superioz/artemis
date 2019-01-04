package rest

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func GetRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.GET("/", Index)

	return router
}

func Index(ctx *fasthttp.RequestCtx) {
	_, _ = ctx.Write([]byte("Welcome!"))
}
