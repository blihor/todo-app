package main

import (
	"log/slog"
	"net/http"

	"github.com/blihor/todo-app/internal/auth"
	"github.com/blihor/todo-app/internal/config"
	"github.com/blihor/todo-app/internal/user"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateRoutes(logger *slog.Logger, client *mongo.Client, cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	userService := user.NewService(client.Database("todo").Collection("users"), logger)
	userHandler := user.NewHandler(userService, logger)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("POST /users", userHandler.CreateUser)
	mux.HandleFunc("DELETE /users/{id}", userHandler.DeleteUserByID)
	mux.HandleFunc("PUT /users/{id}", userHandler.UpdateUserByID)

	authService := auth.NewService(userService, logger, cfg.SecretJwt)
	authHandler := auth.NewHandler(authService, logger)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/register", authHandler.Register)
	return mux
}
