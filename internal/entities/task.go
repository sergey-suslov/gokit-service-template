package entities

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID          primitive.ObjectID  `json:"_id" bson:"_id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	UserId      string              `json:"userId"`
	Created     primitive.Timestamp `json:"created"`
}

func NewTask(name, description, userId string) Task {
	return Task{Name: name, Description: description, UserId: userId, Created: primitive.Timestamp{T: uint32(time.Now().Unix())}}
}
