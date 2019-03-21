package main

import (
	"context"
	"log"
	"os"
	"todoapi/todo"
	"todoapi/todo/mongoservice"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DATABASE_URL")))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	coll := client.Database(os.Getenv("DATABASE_NAME")).Collection("todos")
	handler := &todo.Handler{
		Service: &mongoservice.TodoList{
			C: coll,
		},
	}

	r := gin.Default()

	r.POST("/todos", handler.AddHandler())
	r.Run(":8000")
}
