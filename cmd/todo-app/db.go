package main

import (
	"context"
	"fmt"
	"time"

	"github.com/blihor/todo-app/internal/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

func Connect(cfg *config.Config) (*mongo.Client, error) {
	opts := options.Client().ApplyURI(cfg.DBConnStr)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("Couldn't ping db: %w", err)
	}

	return client, nil
}
