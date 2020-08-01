package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"headless-todo-tasks-service/internal/adapters/endpoints"
	"log"
	"net/http"
)

func main() {
	initConfig()

	client, closeConnection := ConnectMongo()
	defer func() { closeConnection() }()

	c := Init(client)

	http.Handle("/create-task", endpoints.CreateTaskHandler(c))
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
