package service

import (
	"GO_web/model"
	"context"
)

type TodoStore interface {
	Create(ctx context.Context, todo model.Todo) (model.Todo, error)
	List(ctx context.Context) ([]model.Todo, error)
	Get(ctx context.Context, id int) (model.Todo, error)
	Delete(ctx context.Context, id int) error
	Update(ctx context.Context, id int, todo model.Todo) error
}
type TodoService struct {
	repo TodoStore
}

func NewTodoService(repo TodoStore) *TodoService {
	return &TodoService{
		repo: repo,
	}
}

func (s *TodoService) Create(ctx context.Context, todo model.Todo) (model.Todo, error) {
	if todo.Title == "" {
		return model.Todo{}, model.ErrEmptyContent
	}

	return s.repo.Create(ctx, todo)
}

func (s *TodoService) List(ctx context.Context) ([]model.Todo, error) {
	return s.repo.List(ctx)
}

func (s *TodoService) Get(ctx context.Context, id int) (model.Todo, error) {
	return s.repo.Get(ctx, id)
}

func (s *TodoService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func (s *TodoService) Update(ctx context.Context, id int, todo model.Todo) error {
	return s.repo.Update(ctx, id, todo)
}
