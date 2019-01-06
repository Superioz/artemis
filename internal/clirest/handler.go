package clirest

import (
	"github.com/sirupsen/logrus"
	"github.com/valyala/fasthttp"
)

func Index(ctx *fasthttp.RequestCtx) {
	_, err := ctx.Write([]byte("artemis internal rest api"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"prefix": "clirest",
		}).Errorln("couldn't handle", ctx.Request.String())
	}
}
