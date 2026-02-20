package handler

import (
	"net/http"
	"sort"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
	"github.com/philipjkim/gothic-to-do/internal/model"
	"github.com/philipjkim/gothic-to-do/internal/store"
	"github.com/philipjkim/gothic-to-do/templates"
)

type TodoHandler struct {
	store *store.TodoStore
}

func NewTodoHandler(store *store.TodoStore) *TodoHandler {
	return &TodoHandler{store: store}
}

// render 는 templ 컴포넌트를 gin.Context 에 렌더링하는 헬퍼입니다.
func render(c *gin.Context, component templ.Component) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := component.Render(c.Request.Context(), c.Writer); err != nil {
		c.String(http.StatusInternalServerError, "render error: %v", err)
	}
}

func (h *TodoHandler) Index(c *gin.Context) {
	todos := h.sortedTodos()
	render(c, templates.Index(todos))
}

func (h *TodoHandler) Create(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	h.store.Create(title)
	render(c, templates.TodoList(h.sortedTodos()))
}

func (h *TodoHandler) Toggle(c *gin.Context) {
	id := c.Param("id")
	if _, err := h.store.Toggle(id); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	render(c, templates.TodoList(h.sortedTodos()))
}

func (h *TodoHandler) Update(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	if title == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	if _, err := h.store.Update(id, title); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	render(c, templates.TodoList(h.sortedTodos()))
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.store.Delete(id); err != nil {
		c.Status(http.StatusNotFound)
		return
	}
	render(c, templates.TodoList(h.sortedTodos()))
}

func (h *TodoHandler) sortedTodos() []*model.Todo {
	todos := h.store.GetAll()
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].CreatedAt.After(todos[j].CreatedAt)
	})
	return todos
}
