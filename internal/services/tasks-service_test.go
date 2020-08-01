package services

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"headless-todo-tasks-service/internal/entities"
	"testing"
)

type MockedTasksRepository struct {
	mock.Mock
}

func (m *MockedTasksRepository) Create(ctx context.Context, name, description, userId string) (*entities.Task, error) {
	args := m.Called(ctx, name, description, userId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entities.Task), args.Error(1)
}

func TestTasksService_Create_Error(t *testing.T) {
	t.Run("Empty name", func(t *testing.T) {
		mockedTaskRepository := new(MockedTasksRepository)
		testError := errors.New("name must be present")
		service := NewTasksService(mockedTaskRepository)
		task, err := service.Create(context.Background(), "", "description", "1")

		mockedTaskRepository.AssertNotCalled(t, "Create", context.Background(), "", "description", "1")
		assert.Nil(t, task)
		if assert.Error(t, err) {
			assert.EqualError(t, err, testError.Error())
		}
		mockedTaskRepository.AssertExpectations(t)
	})

	t.Run("Empty description", func(t *testing.T) {
		mockedTaskRepository := new(MockedTasksRepository)
		testError := errors.New("description must be present")
		service := NewTasksService(mockedTaskRepository)
		task, err := service.Create(context.Background(), "name", "", "1")

		mockedTaskRepository.AssertNotCalled(t, "Create", context.Background(), "name", "", "1")
		assert.Nil(t, task)
		if assert.Error(t, err) {
			assert.EqualError(t, err, testError.Error())
		}
		mockedTaskRepository.AssertExpectations(t)
	})

	t.Run("Empty userId", func(t *testing.T) {
		mockedTaskRepository := new(MockedTasksRepository)
		testError := errors.New("userId must be present")
		service := NewTasksService(mockedTaskRepository)
		task, err := service.Create(context.Background(), "name", "description", "")

		mockedTaskRepository.AssertNotCalled(t, "Create", context.Background(), "name", "description", "")
		assert.Nil(t, task)
		if assert.Error(t, err) {
			assert.EqualError(t, err, testError.Error())
		}
		mockedTaskRepository.AssertExpectations(t)
	})
}

func TestTasksService_Create_Success(t *testing.T) {
	mockedTaskRepository := new(MockedTasksRepository)
	testTask := &entities.Task{
		ID:          primitive.ObjectID{1},
		Name:        "task 1",
		Description: "description 1",
		UserId:      "2",
	}
	mockedTaskRepository.On("Create", context.Background(), testTask.Name, testTask.Description, testTask.UserId).Return(testTask, nil)
	service := NewTasksService(mockedTaskRepository)

	task, _ := service.Create(context.Background(), testTask.Name, testTask.Description, testTask.UserId)

	mockedTaskRepository.AssertCalled(t, "Create", context.Background(), testTask.Name, testTask.Description, testTask.UserId)
	mockedTaskRepository.AssertExpectations(t)
	assert.Equal(t, task, testTask)
}
