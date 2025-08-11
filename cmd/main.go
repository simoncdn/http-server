package main

import (
	"log"

	"github.com/simoncdn/http-server/internal/config"
	"github.com/simoncdn/http-server/internal/server"
)

func main() {
	cfg := config.New()
	server := server.New(cfg)

	err := server.Start()
	if err != nil {
    log.Fatal("Failed to start server:", err)
	}
}
