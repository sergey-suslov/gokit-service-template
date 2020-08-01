package repositories

import (
	"context"
	"headless-todo-tasks-service/internal/entities"
)

type TasksRepository interface {
	Create(context.Context, string, string, string) (*entities.Task, error)
}
