package endpoints

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"go.uber.org/dig"
	"headless-todo-tasks-service/internal/adapters/middlewares"
	"headless-todo-tasks-service/internal/services"
	"log"
	"net/http"
)

type createTaskRequest struct {
	UserClaim
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *createTaskRequest) SetUserClaim(claim UserClaim) {
	c.UserClaim = claim
}

func makeCreateTaskEndpoint(service services.TasksService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*createTaskRequest)
		task, err := service.Create(ctx, req.Name, req.Description, req.UserClaim.ID)
		if err != nil {
			return nil, err
		}
		return task, nil
	}
}

func CreateTaskHandler(c *dig.Container) http.Handler {
	var service services.TasksService
	if err := c.Invoke(func(s services.TasksService) {
		service = s
	}); err != nil {
		log.Fatal(err)
	}

	var logger kitlog.Logger

	if err := c.Invoke(func(log kitlog.Logger) {
		logger = log
	}); err != nil {
		log.Fatal(err)
	}

	service = &middlewares.LoggerMiddleware{Logger: kitlog.With(logger), Next: service}
	taskEndpoint := makeCreateTaskEndpoint(service)

	return httptransport.NewServer(
		taskEndpoint,
		DefaultRequestDecoder(func(r *http.Request) (UserClaimable, error) {
			var request createTaskRequest
			if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
				return nil, err
			}
			return &request, nil
		}),
		DefaultRequestEncoder,
	)
}
