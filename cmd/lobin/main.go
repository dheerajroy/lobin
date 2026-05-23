package main

import (
	"context"
	"flag"

	"github.com/dheerajroy/lobin/internal/config"
	"github.com/dheerajroy/lobin/internal/health"
	"github.com/dheerajroy/lobin/internal/logger"
	"github.com/dheerajroy/lobin/internal/server"
	"github.com/dheerajroy/lobin/internal/strategy"
	"github.com/dheerajroy/lobin/internal/upstream"
)

func main() {

	configPath := flag.String(
		"config",
		"./config.yaml",
		"path to config file",
	)

	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		panic(err)
	}

	logger.Init(cfg.Logging.Level)
	logger.Info.Printf("Logger initialized with level: %s\n", cfg.Logging.Level)

	err = config.Validate(cfg)
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	upstreams := upstream.BuildFromConfig(cfg)
	logger.Info.Printf("Loaded %d upstream(s) from config\n", len(upstreams))

	for _, upstream := range upstreams {
		logger.Debug.Printf("Starting health check for: %s\n", upstream.URL.String())
		go health.Start(ctx, upstream)
	}

	strategy := strategy.BuildStrategy(cfg.Strategy)

	server := server.NewServer(
		upstreams,
		strategy,
		cfg,
	)

	server.Serve()
}
