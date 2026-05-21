package main

import (
	"github.com/Goga211/go-auth/internal/config"
	"github.com/Goga211/go-auth/internal/logger"
)

func main() {
	cfg := config.Load()
	log := logger.SetupLogger(cfg.Env)

	log.Info("Starting app...")
	
}
