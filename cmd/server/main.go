package main

import (
	"flag"
	"fmt"

	"github.com/k1ender/go-stash/internal/config"
)

func main() {
	filepath := flag.String("config", "", "Path to config file")
	flag.Parse()

	var cfg *config.Config

	if *filepath != "" {
		cfg = config.LoadConfig("config", config.WithConfigPath(*filepath))
	} else {
		cfg = config.LoadConfig("cli")
	}

	fmt.Printf("Server starting at %s:%d\n", cfg.Host, cfg.Port)
}
