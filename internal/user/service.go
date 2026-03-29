package user

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	collection *mongo.Collection
	logger     *slog.Logger
}

func NewService(collection *mongo.Collection, logger *slog.Logger) *service {
	return &service{
		collection: collection,
		logger:     logger,
	}
}

func (s *service) FindOne(ctx context.Context, key string, value any) (*User, error) {
	var result *User
	filter := bson.M{key: value}

	err := s.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to find a user with %s = %s", key, value), "error", err)
	}

	return result, err
}

func (s *service) Create(ctx context.Context, dto *CreateUserDTO) (*mongo.InsertOneResult, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.MinCost)
	if err != nil {
		s.logger.Error("Failed to hash a password", "error", err)
		return nil, err
	}

	user := &User{
		ID:        bson.NewObjectID(),
		Email:     dto.Email,
		Password:  string(hashed_password),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, err := s.collection.InsertOne(ctx, user)
	if err != nil {
		s.logger.Error("Failed to insert a user", "error", err)
	}

	return result, err
}

func (s *service) Update(ctx context.Context, id bson.ObjectID, dto *UpdateUserDTO) (*mongo.UpdateResult, error) {
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(dto.Password), bcrypt.MinCost)
	if err != nil {
		s.logger.Error("Failed to hash a password", "error", err)
		return nil, err
	}

	updateData := bson.M{
		"email":     dto.Email,
		"password":  string(hashed_password),
		"updatedAt": time.Now(),
	}

	result, err := s.collection.UpdateByID(ctx, id, bson.D{
		{Key: "$set", Value: updateData},
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to update a user with id: %s",
			id.String()), "error", err)
	}

	if result.MatchedCount == 0 {
		s.logger.Error(fmt.Sprintf("User with id = %s not found", id.String()))
		return nil, fmt.Errorf("User not found")
	}

	return result, err
}

func (s *service) Delete(ctx context.Context, id bson.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to delete user with id: %s",
			id.String()), "error", err)
	}

	if result.DeletedCount == 0 {
		s.logger.Error(fmt.Sprintf("User with id = %s not found", id.String()))
		return nil, fmt.Errorf("User not found")
	}

	return result, err
}
