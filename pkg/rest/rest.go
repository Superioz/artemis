package rest

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/superioz/artemis/pkg/rest/signature"
	"github.com/valyala/fasthttp"
)

type Server struct {
	host string
	port int

	router *fasthttprouter.Router
	auth   *signature.Authenticator
}

func New(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

func (s *Server) Up() error {
	router := GetRouter()
	s.router = router

	err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), router.Handler)
	if err != nil {
		return err
	}
	return nil
}

func GetRouter() *fasthttprouter.Router {
	router := fasthttprouter.New()
	return router
}
