package main

import (
	kitlog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/dig"
	"headless-todo-tasks-service/internal/adapters/middlewares"
	"headless-todo-tasks-service/internal/adapters/repositories"
	"headless-todo-tasks-service/internal/services"
	"log"
	"os"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Init(client *mongo.Client) *dig.Container {
	c := dig.New()

	err := c.Provide(func() *mongo.Client {
		return client
	})
	handleError(err)

	err = c.Provide(func() *mongo.Database {
		return client.Database(viper.GetString("DB_NAME"))
	})
	handleError(err)

	err = c.Provide(repositories.NewTasksRepositoryMongo)
	handleError(err)

	err = c.Provide(func() *middlewares.PrometheusMetrics {
		fieldKeys := []string{"method", "error"}
		requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: viper.GetString("METRICS_NAMESPACE"),
			Subsystem: viper.GetString("METRICS_SUBSYSTEM"),
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys)
		requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: viper.GetString("METRICS_NAMESPACE"),
			Subsystem: viper.GetString("METRICS_SUBSYSTEM"),
			Name:      "request_latency_microseconds",
			Help:      "Total duration of requests in microseconds.",
		}, fieldKeys)
		countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: viper.GetString("METRICS_NAMESPACE"),
			Subsystem: viper.GetString("METRICS_SUBSYSTEM"),
			Name:      "count_result",
			Help:      "The result of each count method.",
		}, []string{}) // no fields here
		return middlewares.NewPrometheusMetrics(requestCount, requestLatency, countResult)
	})

	err = c.Provide(services.NewTasksService)
	handleError(err)

	err = c.Provide(func() kitlog.Logger {
		logger := kitlog.NewLogfmtLogger(os.Stderr)
		return logger
	})
	handleError(err)

	return c
}
