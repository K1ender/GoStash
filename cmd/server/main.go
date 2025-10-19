package main

import (
	"flag"

	"github.com/k1ender/go-stash/internal/config"
	"github.com/k1ender/go-stash/internal/server"
)

func main() {
	// profiler.StartPprof()
	filepath := flag.String("config", "", "Path to config file")
	flag.Parse()

	// logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	// 	AddSource: true,
	// 	Level:     slog.LevelDebug,
	// }))
	// slog.SetDefault(logger)

	var cfg *config.Config

	if *filepath != "" {
		cfg = config.LoadConfig("config", config.WithConfigPath(*filepath))
	} else {
		cfg = config.LoadConfig("cli")
	}

	srv := server.NewServer(cfg)

	srv.Start()
}
