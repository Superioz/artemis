package rest

import (
	"github.com/buaazp/fasthttprouter"
)

func GetRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	return router
}
