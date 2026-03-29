package auth

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/blihor/todo-app/internal/user"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	FindOne(ctx context.Context, key string, value any) (*user.User, error)
	Create(ctx context.Context, dto *user.CreateUserDTO) (*mongo.InsertOneResult, error)
	Delete(ctx context.Context, id bson.ObjectID) (*mongo.DeleteResult, error)
	Update(ctx context.Context, id bson.ObjectID, dto *user.UpdateUserDTO) (*mongo.UpdateResult, error)
}

type service struct {
	userService UserService
	logger      *slog.Logger
	sercetJwt   string
}

func NewService(userService UserService, logger *slog.Logger, secretJwt string) *service {
	return &service{
		userService: userService,
		logger:      logger,
		sercetJwt:   secretJwt,
	}
}

func (s *service) Login(ctx context.Context, dto *UserLoginDTO) (string, error, int) {
	user, err := s.userService.FindOne(ctx, "email", dto.Email)
	if err != nil {
		s.logger.Error(fmt.Sprintf("User with emali %s not found", dto.Email), "error", err)
		return "", err, http.StatusUnauthorized
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password)); err != nil {
		s.logger.Error("Passwords don't match", "error", err)
		return "", err, http.StatusUnauthorized
	}

	token, err := GenerateToken(user.ID, []byte(s.sercetJwt))
	if err != nil {
		s.logger.Error("Failed to signed jwt token", "error", err)
		return "", nil, http.StatusInternalServerError
	}

	return token, nil, http.StatusOK
}

func (s *service) Register(ctx context.Context, dto *UserRegisterDTO) (any, error, int) {
	_, err := s.userService.FindOne(ctx, "email", dto.Email)
	if err == nil {
		s.logger.Error(fmt.Sprintf("User with emali %s already exists", dto.Email), "error", err)
		return bson.NilObjectID, fmt.Errorf("Email already exists"), http.StatusUnauthorized
	}

	createDTO := &user.CreateUserDTO{
		Email:    dto.Email,
		Password: dto.Password,
	}

	result, err := s.userService.Create(ctx, createDTO)
	if err != nil {
		s.logger.Error("Failed to create user", "error", err)
		return bson.NilObjectID, err, http.StatusInternalServerError
	}

	return result.InsertedID, nil, http.StatusCreated
}
