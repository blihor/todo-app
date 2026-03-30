package main

import (
	"log/slog"
	"net/http"

	"github.com/blihor/todo-app/internal/auth"
	"github.com/blihor/todo-app/internal/config"
	"github.com/blihor/todo-app/internal/middleware"
	"github.com/blihor/todo-app/internal/task"
	"github.com/blihor/todo-app/internal/user"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func CreateRoutes(logger *slog.Logger, client *mongo.Client, cfg *config.Config) *http.ServeMux {
	mux := http.NewServeMux()
	userService := user.NewService(client.Database("todo").Collection("users"), logger)
	userHandler := user.NewHandler(userService, logger)

	authService := auth.NewService(userService, logger, cfg.SecretJwt)
	authHandler := auth.NewHandler(authService, logger)

	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/register", authHandler.Register)

	jwtMiddleware := middleware.NewMiddleware(authService, logger)

	mux.HandleFunc("GET /users/{id}", jwtMiddleware.Protect(userHandler.GetByID))
	mux.HandleFunc("POST /users", jwtMiddleware.Protect(userHandler.Create))
	mux.HandleFunc("DELETE /users/{id}", jwtMiddleware.Protect(userHandler.DeleteByID))
	mux.HandleFunc("PUT /users/{id}", jwtMiddleware.Protect(userHandler.UpdateByID))

	taskService := task.NewService(client.Database("todo").Collection("tasks"), logger)
	taskHandler := task.NewHandler(taskService, logger)

	mux.HandleFunc("GET /tasks/{id}", jwtMiddleware.Protect(taskHandler.GetByID))
	mux.HandleFunc("GET /tasks/title/{title}", jwtMiddleware.Protect(taskHandler.GetByTitle))
	mux.HandleFunc("GET /tasks/owner/{id}", jwtMiddleware.Protect(taskHandler.GetByOwnerID))
	mux.HandleFunc("POST /tasks", jwtMiddleware.Protect(taskHandler.Create))
	mux.HandleFunc("DELETE /tasks/{id}", jwtMiddleware.Protect(taskHandler.DeleteByID))
	mux.HandleFunc("PUT /tasks/{id}", jwtMiddleware.Protect(taskHandler.UpdateByID))

	return mux
}
