package handler

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/internal/dome"
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

func Status(ctx *fasthttp.RequestCtx) {
	status := dome.CurrentStatus()
	data, err := json.Marshal(status)
	if err != nil {
		ctx.Error("couldn't fetch current status", 500)
		return
	}
	ctx.Success("application/json", data)
}
