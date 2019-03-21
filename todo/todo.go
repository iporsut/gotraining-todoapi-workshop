package todo

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

var ErrNotFound = errors.New("not found")

type Service interface {
	Add(context.Context, *TodoList) (*TodoList, error)
}

type Handler struct {
	Service Service
}

func (h *Handler) wrapError(gh func(*gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := gh(c); err != nil {
			if err == ErrNotFound {
				c.AbortWithError(http.StatusNotFound, err)
			} else {
				c.AbortWithError(http.StatusInternalServerError, err)
			}
		}
	}
}

func (h *Handler) AddHandler() gin.HandlerFunc {
	return h.wrapError(func(c *gin.Context) error {
		todo := &TodoList{
			Tasks: []Task{},
		}
		if err := c.Bind(todo); err != nil {
			return err
		}
		todo, err := h.Service.Add(c.Request.Context(), todo)
		if err != nil {
			return err
		}
		c.JSON(http.StatusOK, todo)
		return nil
	})
}
