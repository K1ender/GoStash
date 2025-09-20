package main

import (
	"fmt"

	"github.com/k1ender/go-stash/internal/config"
)

func main() {
	cfg := config.LoadConfig("cli")
	fmt.Println(cfg)
}
