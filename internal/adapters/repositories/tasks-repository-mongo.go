package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"headless-todo-tasks-service/internal/entities"
	"headless-todo-tasks-service/internal/services/repositories"
)

const TasksCollection = "tasks"

type TasksRepositoryMongo struct {
	db *mongo.Database
}

func NewTasksRepositoryMongo(db *mongo.Database) repositories.TasksRepository {
	return &TasksRepositoryMongo{db}
}

func (r *TasksRepositoryMongo) Create(ctx context.Context, name, description, userId string) (*entities.Task, error) {
	task := entities.NewTask(name, description, userId)
	result, err := r.db.Collection(TasksCollection).InsertOne(ctx, bson.M{"name": task.Name, "userId": task.UserId, "description": task.Description, "created": task.Created})
	if err != nil {
		return nil, err
	}
	task.ID = result.InsertedID.(primitive.ObjectID)
	return &task, nil
}
