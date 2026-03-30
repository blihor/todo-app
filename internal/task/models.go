package task

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Status int

const (
	StatusInProgress Status = iota
	StatusDone
)

type Task struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	OwnerID     bson.ObjectID `bson:"owner_id,omitempty" json:"owner_id"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Status      Status        `bson:"status" json:"status"`
	CreatedAt   time.Time     `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time     `bson:"updatedAt" json:"updatedAt"`
}

type CreateTaskDTO struct {
	OwnerID     bson.ObjectID `bson:"owner_id,omitempty" json:"owner_id"`
	Title       string        `bson:"title" json:"title"`
	Description string        `bson:"description" json:"description"`
	Status      Status        `bson:"status" json:"status"`
}

type UpdateTaskDTO struct {
	Title       string `bson:"title" json:"title"`
	Description string `bson:"description" json:"description"`
	Status      Status `bson:"status" json:"status"`
}
