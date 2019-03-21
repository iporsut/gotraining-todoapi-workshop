package mongoservice

import (
	"context"
	"todoapi/todo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoList struct {
	C *mongo.Collection
}

func (s *TodoList) Add(ctx context.Context, todo *todo.TodoList) (*todo.TodoList, error) {
	result, err := s.C.InsertOne(ctx, todo)
	if err != nil {
		return nil, err
	}
	todo.ID = result.InsertedID.(primitive.ObjectID)
	return todo, nil
}
