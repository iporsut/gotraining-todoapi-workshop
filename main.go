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

type UpdateTodoListParam struct {
	Title string `json:"title"`
}

type UpdateTaskParam struct {
	Done bool `json:"done"`
}

func wrapError(coll *mongo.Collection, h func(context.Context, *gin.Context, *mongo.Collection) error) func(*gin.Context) {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		err := h(ctx, c, coll)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func newTodoList() *TodoList {
	return &TodoList{
		Tasks: []Task{},
	}
}

func listTodoList(ctx context.Context, coll *mongo.Collection) ([]*TodoList, error) {
	cur, err := coll.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	todos := []*TodoList{}
	for cur.Next(ctx) {
		todo := newTodoList()
		if err := cur.Decode(todo); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := cur.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func addTodoList(ctx context.Context, coll *mongo.Collection, todo *TodoList) (*TodoList, error) {
	result, err := coll.InsertOne(ctx, todo)
	if err != nil {
		return nil, err
	}
	todo.ID = result.InsertedID.(primitive.ObjectID)
	return todo, nil
}

func findTodoList(ctx context.Context, coll *mongo.Collection, id string) (*TodoList, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	todo := &TodoList{
		Tasks: []Task{},
	}
	err = coll.FindOne(ctx, bson.D{{"_id", oid}}).Decode(todo)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func updateTodoList(ctx context.Context, coll *mongo.Collection, id string, param *UpdateTodoListParam) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = coll.UpdateOne(
		ctx,
		bson.D{{"_id", oid}},
		bson.D{{"$set", param}})
	return err
}

func deleteTodoList(ctx context.Context, coll *mongo.Collection, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = coll.DeleteOne(ctx, bson.D{{"_id", oid}})
	return err
}

func addTask(ctx context.Context, coll *mongo.Collection, id string, task *Task) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if _, err := coll.UpdateOne(
		ctx,
		bson.D{{"_id", oid}},
		bson.D{{"$push",
			bson.D{{"tasks", task}},
		}},
	); err != nil {
		return err
	}
	return nil
}

func updateTask(ctx context.Context, coll *mongo.Collection, id string, taskID string, param *UpdateTaskParam) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	if _, err := coll.UpdateOne(
		ctx,
		bson.D{{"_id", oid}},
		bson.D{{"$set", bson.D{{"tasks." + taskID + ".done", param.Done}}}},
	); err != nil {
		return err
	}
	return nil
}

func listTodoListHandler(ctx context.Context, c *gin.Context, coll *mongo.Collection) error {
	todos, err := listTodoList(ctx, coll)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, todos)
	return nil
}

func addTodoListHandler(ctx context.Context, c *gin.Context, coll *mongo.Collection) error {
	todo := &TodoList{
		Tasks: []Task{},
	}
	if err := c.Bind(todo); err != nil {
		return err
	}
	todo, err := addTodoList(ctx, coll, todo)
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, todo)
	return nil
}

func getTodoListHandler(ctx context.Context, c *gin.Context, coll *mongo.Collection) error {
	todo, err := findTodoList(ctx, coll, c.Param("id"))
	if err != nil {
		return err
	}
	c.JSON(http.StatusOK, todo)
	return nil
}

func updateTodoListHandler(ctx context.Context, c *gin.Context, coll *mongo.Collection) error {
	var param UpdateTodoListParam
	if err := c.Bind(&param); err != nil {
		return err
	}
	return updateTodoList(ctx, coll, c.Param("id"), &param)
}

func deleteTodoListHandler(ctx context.Context, c *gin.Context, coll *mongo.Collection) error {
	return deleteTodoList(ctx, coll, c.Param("id"))
}

func addTaskHandler(ctx context.Context, c *gin.Context, coll *mongo.Collection) error {
	var t Task
	if err := c.Bind(&t); err != nil {
		return err
	}
	return addTask(ctx, coll, c.Param("id"), &t)
}

func updateTaskHanlder(ctx context.Context, c *gin.Context, coll *mongo.Collection) error {
	var param UpdateTaskParam
	if err := c.Bind(&param); err != nil {
		return err
	}
	return updateTask(ctx, coll, c.Param("id"), c.Param("task_id"), &param)
}

func main() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://db:27017"))
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

	r.GET("/todos", wrapError(coll, listTodoListHandler))
	r.POST("/todos", wrapError(coll, addTodoListHandler))
	r.GET("/todos/:id", wrapError(coll, getTodoListHandler))
	r.PUT("/todos/:id", wrapError(coll, updateTodoListHandler))
	r.DELETE("/todos/:id", wrapError(coll, deleteTodoListHandler))
	r.POST("/todos/:id/tasks", wrapError(coll, addTaskHandler))
	r.PUT("/todos/:id/tasks/:task_id", wrapError(coll, updateTaskHanlder))

	r.Run(":8000")
}
