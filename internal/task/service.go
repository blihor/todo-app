package task

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (s *service) FindOne(ctx context.Context, key string, value any) (*Task, error) {
	var result *Task
	filter := bson.M{key: value}

	err := s.collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to find a task with %s = %s", key, value), "error", err)
	}

	return result, err
}

func (s *service) FindMany(ctx context.Context, key string, value any) ([]Task, error) {
	var result []Task
	filter := bson.M{key: value}

	cursor, err := s.collection.Find(ctx, filter)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to find a task with %s = %s", key, value), "error", err)
	}

	cursor.All(ctx, &result)
	return result, err
}

func (s *service) Create(ctx context.Context, dto *CreateTaskDTO) (*mongo.InsertOneResult, error) {
	task := &Task{
		ID:          bson.NewObjectID(),
		OwnerID:     dto.OwnerID,
		Title:       dto.Title,
		Description: dto.Description,
		Status:      StatusInProgress,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result, err := s.collection.InsertOne(ctx, task)
	if err != nil {
		s.logger.Error("Failed to insert a task", "error", err)
	}

	return result, err
}

func (s *service) Update(ctx context.Context, id bson.ObjectID, dto *UpdateTaskDTO) (*mongo.UpdateResult, error) {
	updateData := bson.M{
		"title":       dto.Title,
		"description": dto.Description,
		"status":      dto.Status,
		"updatedAt":   time.Now(),
	}

	result, err := s.collection.UpdateByID(ctx, id, bson.D{
		{Key: "$set", Value: updateData},
	})
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to update a task with id: %s",
			id.String()), "error", err)
	}

	if result.MatchedCount == 0 {
		s.logger.Error(fmt.Sprintf("Task with id = %s not found", id.String()))
		return nil, fmt.Errorf("Task not found")
	}

	return result, err
}

func (s *service) Delete(ctx context.Context, id bson.ObjectID) (*mongo.DeleteResult, error) {
	filter := bson.M{"_id": id}
	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		s.logger.Error(fmt.Sprintf("Failed to delete task with id: %s",
			id.String()), "error", err)
	}

	if result.DeletedCount == 0 {
		s.logger.Error(fmt.Sprintf("Task with id = %s not found", id.String()))
		return nil, fmt.Errorf("Task not found")
	}

	return result, err
}
