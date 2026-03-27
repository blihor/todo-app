package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/blihor/todo-app/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	config *config.Config
	logger *slog.Logger
	client *mongo.Client
}

func (app *Application) run(h http.Handler) error {
	server := &http.Server{
		Addr:         app.config.Port,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Info("Starting server", "port", app.config.Port)

	return server.ListenAndServe()
}
