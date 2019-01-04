package rest

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type Server struct {
	host string
	port int
}

func New(host string, port int) *Server {
	return &Server{
		host: host,
		port: port,
	}
}

func (s *Server) Up() error {
	router := GetRouter()

	err := fasthttp.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), router.Handler)
	if err != nil {
		return err
	}
	return nil
}
