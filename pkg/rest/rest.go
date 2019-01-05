package rest

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
	"github.com/superioz/artemis/config"
	"github.com/superioz/artemis/pkg/rest/signature"
	"github.com/valyala/fasthttp"
	"time"
)

// represents a rest server
type Server struct {
	// the config for the host, port, ...
	config config.Rest

	// the host, probably `0.0.0.0` or `localhost`
	host string

	// the port of the server
	port int

	// the internal http server
	server *fasthttp.Server

	// router to add routes to before `Up()`
	router *fasthttprouter.Router

	// authentication for requests
	auth *signature.Authenticator

	// channel for error handling
	err chan error
}

// creates a new server instance with given configuration.
// default is `config.DefaultRestConfig()`
func New(config config.Rest) *Server {
	router := fasthttprouter.New()
	s := &Server{
		config: config,
		router: router,
		host:   config.Host,
	}
	s.server = &fasthttp.Server{
		ReadTimeout:          time.Duration(config.ReadTimeout) * time.Second,
		WriteTimeout:         time.Duration(config.WriteTimeout) * time.Second,
		MaxConnsPerIP:        config.MaxConnsPerIP,
		MaxKeepaliveDuration: time.Duration(config.MaxKeepaliveDuration) * time.Second,
		Logger:               logrus.StandardLogger(),
	}
	return s
}

// starts the server. returns an error if the port couldn't be
// found or if the server can't bind to that address.
func (s *Server) Up() error {
	s.server.Handler = s.router.Handler

	// get free port in config range
	// default `2310`-`2315`
	port, err := GetFreePortInRange(s.config.Host, s.config.MinPort, s.config.MaxPort)
	if err != nil {
		return err
	}
	s.port = port

	go func() {
		err := s.server.ListenAndServe(fmt.Sprintf("%s:%d", s.config.Host, port))
		if err != nil {
			select {
			case s.err <- err:
			}
		}
	}()
	return nil
}

// channel for error listening. mainly if the server stops/shutdowns.
func (s *Server) ErrListen() <-chan error {
	return s.err
}

// the address of the server as combination of
// host and port.
func (s *Server) Address() string {
	return fmt.Sprintf("%s:%d", s.host, s.port)
}
