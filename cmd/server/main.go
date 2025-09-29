package main

import (
	"flag"

	"github.com/k1ender/go-stash/internal/config"
	"github.com/k1ender/go-stash/internal/server"
)

func main() {
	filepath := flag.String("config", "", "Path to config file")
	flag.Parse()
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	var cfg *config.Config

	if *filepath != "" {
		cfg = config.LoadConfig("config", config.WithConfigPath(*filepath))
	} else {
		cfg = config.LoadConfig("cli")
	}

	srv := server.NewServer(cfg)

	srv.Start()
}
