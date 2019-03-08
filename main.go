package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://usko7imkex32llnhtmyz:AvsbbVN7Dzo9jp58pPoT@bazk4bps4tqprxp-mongodb.services.clever-cloud.com:27017/bazk4bps4tqprxp"))
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	if err := client.Ping(context.Background(), nil); err != nil {
		log.Println(err)
		os.Exit(1)
	}
	coll := client.Database("bazk4bps4tqprxp").Collection("todos")
	r := gin.Default()
	r.GET("/todos", func(c *gin.Context) {
		ctx := c.Request.Context()
		cur, err := coll.Find(ctx, bson.D{})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		defer cur.Close(ctx)
		var todos []TodoList
		for cur.Next(ctx) {
			todo := TodoList{
				Tasks: []Task{},
			}
			if err := cur.Decode(&todo); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}
			todos = append(todos, todo)
		}
		if err := cur.Err(); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		c.JSON(http.StatusOK, todos)
	})

	r.POST("/todos", func(c *gin.Context) {
		ctx := c.Request.Context()
		todo := TodoList{
			Tasks: []Task{},
		}
		if err := c.Bind(&todo); err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		} // id, _ := primitive.ObjectIDFromHex("aaa")
		result, err := coll.InsertOne(ctx, todo)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		todo.ID = result.InsertedID.(primitive.ObjectID)
		c.JSON(http.StatusOK, todo)
	})

	r.GET("/todos/:id", func(c *gin.Context) {
		ctx := c.Request.Context()
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		todo := TodoList{
			Tasks: []Task{},
		}
		err := coll.FindOne(ctx, bson.D{{"_id", id}}).Decode(&todo)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
			return
		}
		c.JSON(http.StatusOK, todo)
	})

	r.PUT("/todos/:id", func(c *gin.Context) {
		var param struct {
			Title string `json:"title" bson:"title"`
		}
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		if err := c.Bind(&param); err != nil {
			return
		}
		_, err := coll.UpdateOne(
			c.Request.Context(),
			bson.D{{"_id", id}},
			bson.D{{"$set", param}})

		if err != nil {

			c.AbortWithError(http.StatusInternalServerError, err)
		}
	})

	r.DELETE("/todos/:id", func(c *gin.Context) {
		id, _ := primitive.ObjectIDFromHex(c.Param("id"))
		result, err := coll.DeleteOne(c.Request.Context(), bson.D{{"_id", id}})
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if result.DeletedCount == 0 {
			c.Status(http.StatusNotFound)
			return
		}
	})
	r.Run(":8000")
}
