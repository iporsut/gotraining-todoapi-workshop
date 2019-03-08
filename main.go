package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Task struct {
	Desc string `bson:"desc" json:"desc"`
	Done bool   `bson:"done" json:"done"`
}

type TodoList struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title string             `bson:"title" json:"title"`
	Tasks []Task             `bson:"tasks" json:"tasks"`
}

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	collection := client.Database("test").Collection("todolist")
	r := gin.Default()
	r.GET("/todos", func(c *gin.Context) {
		// LAB
	})

	r.POST("/todos", func(c *gin.Context) {
		// LAB
	})

	r.GET("/todos/:id", func(c *gin.Context) {
		// LAB
	})

	r.PUT("/todos/:id", func(c *gin.Context) {
		// LAB
	})

	r.DELETE("/todos/:id", func(c *gin.Context) {
		// LAB
	})
	r.Run(":8080")
}
