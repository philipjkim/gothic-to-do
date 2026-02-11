package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/philipjkim/gothic-to-do/internal/model"
)

type TodoStore struct {
	mu     sync.RWMutex
	todos  map[string]*model.Todo
	nextID int
}

func NewTodoStore() *TodoStore {
	return &TodoStore{
		todos: make(map[string]*model.Todo),
	}
}

func (s *TodoStore) Create(title string) *model.Todo {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.nextID++
	todo := &model.Todo{
		ID:        fmt.Sprintf("todo-%06d", s.nextID),
		Title:     title,
		Completed: false,
		CreatedAt: time.Now(),
	}
	s.todos[todo.ID] = todo
	return todo
}

func (s *TodoStore) GetAll() []*model.Todo {
	s.mu.RLock()
	defer s.mu.RUnlock()

	todos := make([]*model.Todo, 0, len(s.todos))
	for _, t := range s.todos {
		todos = append(todos, t)
	}
	return todos
}

func (s *TodoStore) Toggle(id string) (*model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found: %s", id)
	}
	todo.Completed = !todo.Completed
	return todo, nil
}

func (s *TodoStore) Delete(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.todos[id]; !ok {
		return fmt.Errorf("todo not found: %s", id)
	}
	delete(s.todos, id)
	return nil
}

func (s *TodoStore) Update(id, title string) (*model.Todo, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, fmt.Errorf("todo not found: %s", id)
	}
	todo.Title = title
	return todo, nil
}
