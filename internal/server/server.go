package server

import (
	"fmt"
	"log/slog"
	"net"

	"github.com/k1ender/go-stash/internal/config"
	"github.com/k1ender/go-stash/internal/handler"
	"github.com/k1ender/go-stash/internal/store"
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
	defer func() {
		if r := recover(); r != nil {
			slog.Error("server crashed", "error", r)
		}
	}()

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

	shardedStore := store.NewShardedStore(32)

	newHandler := handler.NewHandler(shardedStore)

	fmt.Printf("Server started on %s:%d\n", s.cfg.Host, s.cfg.Port)

	for {
		client, err := conn.Accept()
		if err != nil {
			continue
		}

		go func(client net.Conn) {
			defer client.Close()
			for {
				isFatal, err := newHandler.Handle(client)
				if err != nil {
					slog.Error("error handling client request", "error", err)
					if isFatal {
						return
					}
				}
			}
		}(client)
	}
}
