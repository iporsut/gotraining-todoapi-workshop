package todo

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/go-playground/assert.v1"
)

type TestService struct {
	oid primitive.ObjectID
	err error
}

func (s *TestService) Add(ctx context.Context, todo *TodoList) (*TodoList, error) {
	todo.ID = s.oid
	return todo, s.err
}

func TestAddHandler(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	oid := primitive.NewObjectID()
	h := &Handler{
		Service: &TestService{
			oid: oid,
		},
	}
	r := gin.New()
	r.POST("/todos", h.AddHandler())
	req := httptest.NewRequest(http.MethodPost, "https://test.com/todos",
		strings.NewReader(`{"title": "List 1"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	resp := w.Result()
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Header.Get("Content-Type"))
	assert.Equal(t, fmt.Sprintf(`{"id":"%s","title":"List 1","tasks":[]}`, oid.Hex()), string(b))
}
