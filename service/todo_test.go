package service

import (
	"GO_web/model"
	"context"
	"errors"
	"reflect"
	"testing"
)

type mockTodoRepo struct {
	todo       model.Todo
	todos      []model.Todo
	createTodo model.Todo
	err        error
	called     int
	gotID      int
	gotTodo    model.Todo
}

func (m *mockTodoRepo) Create(ctx context.Context, todo model.Todo) (model.Todo, error) {
	m.called++
	m.gotTodo = todo
	return m.createTodo, m.err
}

func (m *mockTodoRepo) List(ctx context.Context) ([]model.Todo, error) {
	m.called++
	return m.todos, m.err
}

func (m *mockTodoRepo) Get(ctx context.Context, id int) (model.Todo, error) {
	m.called++
	m.gotID = id
	return m.todo, m.err
}

func (m *mockTodoRepo) Delete(ctx context.Context, id int) error {
	m.called++
	m.gotID = id
	return m.err
}

func (m *mockTodoRepo) Update(ctx context.Context, id int, todo model.Todo) error {
	m.called++
	m.gotID = id
	m.gotTodo = todo
	return m.err
}

func TestCreateTodo(t *testing.T) {
	input := model.Todo{
		Title:   "Test Todo",
		Content: "Test Content",
	}

	expected := model.Todo{
		ID:      1,
		Title:   "Test Todo",
		Content: "Test Content",
	}

	mock := &mockTodoRepo{
		createTodo: expected,
	}

	service := NewTodoService(mock)
	ctx := context.Background()

	recTodo, err := service.Create(ctx, input)

	if err != nil {
		t.Fatalf("意外的错误 %v", err)
	}

	if !reflect.DeepEqual(recTodo, expected) {
		t.Errorf("期望得到 %v，实际得到 %v", expected, recTodo)
	}

	if !reflect.DeepEqual(mock.gotTodo, input) {
		t.Errorf("传入 Repository 的 Todo 不正确")
	}

	if mock.called != 1 {
		t.Errorf("期望调用 Repository 1 次，实际调用 %d 次", mock.called)
	}
}

func TestCreateTodo_ErrCases(t *testing.T) {
	errRepo := errors.New("todoRepoErr")

	tests := []struct {
		name    string
		todo    model.Todo
		repoErr error
		wantErr error
	}{
		{
			name: "空标题",
			todo: model.Todo{
				Title: "",
			},
			wantErr: model.ErrEmptyContent,
		},
		{
			name: "Repository错误",
			todo: model.Todo{
				Title: "Test Todo",
			},
			repoErr: errRepo,
			wantErr: errRepo,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockTodoRepo{
				err: tt.repoErr,
			}

			service := NewTodoService(mock)
			ctx := context.Background()

			_, err := service.Create(ctx, tt.todo)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("期望得到错误 %v，实际得到 %v", tt.wantErr, err)
			}

			if tt.todo.Title == "" {
				if mock.called != 0 {
					t.Errorf("标题为空时不应该调用 Repository")
				}
			} else if mock.called != 1 {
				t.Errorf("期望调用 Repository 1 次，实际调用 %d 次", mock.called)
			}
		})
	}
}

func TestListTodo(t *testing.T) {
	expected := []model.Todo{
		{ID: 1, Title: "Todo1"},
		{ID: 2, Title: "Todo2"},
	}

	mock := &mockTodoRepo{
		todos: expected,
	}

	service := NewTodoService(mock)
	ctx := context.Background()

	recTodos, err := service.List(ctx)

	if err != nil {
		t.Fatalf("意外的错误 %v", err)
	}

	if !reflect.DeepEqual(recTodos, expected) {
		t.Errorf("期望得到 %v，实际得到 %v", expected, recTodos)
	}

	if mock.called != 1 {
		t.Errorf("期望调用 Repository 1 次，实际调用 %d 次", mock.called)
	}
}

func TestListTodo_ErrCases(t *testing.T) {
	errRepo := errors.New("todoRepoErr")

	mock := &mockTodoRepo{
		err: errRepo,
	}

	service := NewTodoService(mock)
	ctx := context.Background()

	_, err := service.List(ctx)

	if !errors.Is(err, errRepo) {
		t.Errorf("期望得到错误 %v，实际得到 %v", errRepo, err)
	}

	if mock.called != 1 {
		t.Errorf("期望调用 Repository 1 次，实际调用 %d 次", mock.called)
	}
}

func TestGetTodo(t *testing.T) {
	expected := model.Todo{
		ID:      1,
		Title:   "Test Todo",
		Content: "Test Content",
	}

	mock := &mockTodoRepo{
		todo: expected,
	}

	service := NewTodoService(mock)
	ctx := context.Background()

	recTodo, err := service.Get(ctx, 1)

	if err != nil {
		t.Fatalf("意外的错误 %v", err)
	}

	if !reflect.DeepEqual(recTodo, expected) {
		t.Errorf("期望得到 %v，实际得到 %v", expected, recTodo)
	}

	if mock.gotID != 1 {
		t.Errorf("期望 ID 为 1，实际得到 %d", mock.gotID)
	}

	if mock.called != 1 {
		t.Errorf("期望调用 Repository 1 次，实际调用 %d 次", mock.called)
	}
}

func TestGetTodo_ErrCases(t *testing.T) {
	errRepo := errors.New("todoRepoErr")

	mock := &mockTodoRepo{
		err: errRepo,
	}

	service := NewTodoService(mock)
	ctx := context.Background()

	_, err := service.Get(ctx, 1)

	if !errors.Is(err, errRepo) {
		t.Errorf("期望得到错误 %v，实际得到 %v", errRepo, err)
	}

	if mock.called != 1 {
		t.Errorf("期望调用 Repository 1 次，实际调用 %d 次", mock.called)
	}
}

// Delete

func TestDeleteTodo(t *testing.T) {
	mock := &mockTodoRepo{}

	service := NewTodoService(mock)
	ctx := context.Background()

	err := service.Delete(ctx, 1)

	if err != nil {
		t.Fatalf("意外的错误 %v", err)
	}

	if mock.gotID != 1 {
		t.Errorf("期望 ID 为 1，实际得到 %d", mock.gotID)
	}

	if mock.called != 1 {
		t.Errorf("期望调用 Repository 1 次，实际调用 %d 次", mock.called)
	}
}

func TestDeleteTodo_ErrCases(t *testing.T) {
	errRepo := errors.New("todoRepoErr")

	mock := &mockTodoRepo{
		err: errRepo,
	}

	service := NewTodoService(mock)
	ctx := context.Background()

	err := service.Delete(ctx, 1)

	if !errors.Is(err, errRepo) {
		t.Errorf("期望得到错误 %v，实际得到 %v", errRepo, err)
	}

	if mock.called != 1 {
		t.Errorf("期望调用 Repository 1 次，实际调用 %d 次", mock.called)
	}
}

func TestUpdateTodo(t *testing.T) {
	input := model.Todo{
		Title:   "Updated Todo",
		Content: "Updated Content",
	}

	mock := &mockTodoRepo{}

	service := NewTodoService(mock)
	ctx := context.Background()

	err := service.Update(ctx, 1, input)

	if err != nil {
		t.Fatalf("意外的错误 %v", err)
	}

	if mock.gotID != 1 {
		t.Errorf("期望 ID 为 1，实际得到 %d", mock.gotID)
	}

	if !reflect.DeepEqual(mock.gotTodo, input) {
		t.Errorf("传入 Repository 的 Todo 不正确")
	}

	if mock.called != 1 {
		t.Errorf("期望调用 Repository 1 次，实际调用 %d 次", mock.called)
	}
}

func TestUpdateTodo_ErrCases(t *testing.T) {
	errRepo := errors.New("todoRepoErr")

	mock := &mockTodoRepo{
		err: errRepo,
	}

	service := NewTodoService(mock)
	ctx := context.Background()

	err := service.Update(ctx, 1, model.Todo{
		Title: "Updated Todo",
	})

	if !errors.Is(err, errRepo) {
		t.Errorf("期望得到错误 %v，实际得到 %v", errRepo, err)
	}

	if mock.called != 1 {
		t.Errorf("期望调用 Repository 1 次，实际调用 %d 次", mock.called)
	}
}
