package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/blihor/todo-app/internal/config"
)

func main() {
	// logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// config
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load .env file", "error", err)
		os.Exit(1)
	}

	// db
	client, err := Connect(cfg)
	if err != nil {
		logger.Error("Failed to connect to db", "error", err)
		os.Exit(1)
	}
	defer client.Disconnect(context.Background())

	app := Application{
		config: cfg,
		logger: logger,
		client: client,
	}

	if err := app.run(nil); err != nil {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
