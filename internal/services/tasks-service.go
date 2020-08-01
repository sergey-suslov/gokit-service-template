package services

import (
	"context"
	"errors"
	"headless-todo-tasks-service/internal/entities"
	"headless-todo-tasks-service/internal/services/repositories"
)

type TasksService interface {
	Create(ctx context.Context, name, description, userId string) (*entities.Task, error)
}

type tasksService struct {
	tasksRepository repositories.TasksRepository
}

func NewTasksService(tasksRepository repositories.TasksRepository) TasksService {
	return &tasksService{tasksRepository}
}

func (service *tasksService) Create(ctx context.Context, name, description, userId string) (*entities.Task, error) {
	if name == "" {
		return nil, errors.New("name must be present")
	}
	if description == "" {
		return nil, errors.New("description must be present")
	}
	if userId == "" {
		return nil, errors.New("userId must be present")
	}
	return service.tasksRepository.Create(ctx, name, description, userId)
}
