package main

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"headless-todo-tasks-service/internal/adapters/repositories"
	"log"
	"time"
)

func ConnectMongo() (*mongo.Client, func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	connectionString := "mongodb://" + viper.GetString("DB_USER") + ":" + viper.GetString("DB_PASSWORD") + "@localhost:27017/" + viper.GetString("DB_NAME")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database(viper.GetString("DB_NAME"))

	userIdIndex := mongo.IndexModel{
		Keys: bson.M{
			"userId": 1,
		},
		Options: nil,
	}
	_, err = database.Collection(repositories.TasksCollection).Indexes().CreateOne(ctx, userIdIndex)
	if err != nil {
		log.Fatal(err)
	}

	return client, func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}
}
