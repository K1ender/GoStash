package server

import (
	"fmt"
	"net"

	"github.com/k1ender/go-stash/internal/config"
	"github.com/k1ender/go-stash/internal/handler"
)

type Server struct {
	cfg *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start() {
	conn, err := net.Listen(
		"tcp",
		net.JoinHostPort(
			s.cfg.Host,
			fmt.Sprintf("%d", s.cfg.Port),
		),
	)
	if err != nil {
		// FIXME: proper logging
		panic(err)
	}
	defer conn.Close()

	handler := handler.NewHandler()

	fmt.Printf("Server started on %s:%d\n", s.cfg.Host, s.cfg.Port)

	for {
		client, err := conn.Accept()
		if err != nil {
			continue
		}

		go handler.Handle(client)
	}
}
