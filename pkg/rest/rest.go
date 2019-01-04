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
	router := fasthttprouter.New()
	return &Server{
		host:   host,
		port:   port,
		router: router,
	}
}

func (s *Server) Up() error {
	err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), s.router.Handler)
	if err != nil {
		return err
	}
	return nil
}
