package handler

import (
	"GO_web/model"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockTodoService struct {
	todo   model.Todo
	todos  []model.Todo
	err    error
	called int
	passID int
}

func (s *mockTodoService) List(ctx context.Context) ([]model.Todo, error) {
	return s.todos, s.err
}
func (s *mockTodoService) Create(ctx context.Context, todo model.Todo) (model.Todo, error) {
	return todo, s.err
}
func (s *mockTodoService) Get(ctx context.Context, id int) (model.Todo, error) {
	s.passID = id
	return s.todo, s.err
}
func (s *mockTodoService) Delete(ctx context.Context, id int) error {
	s.passID = id
	s.todo = model.Todo{}
	return s.err
}
func (s *mockTodoService) Update(ctx context.Context, id int, todo model.Todo) error {
	s.passID = id
	s.todo = todo
	return s.err
}

func TestListTodo(t *testing.T) {
	todos := []model.Todo{
		{
			ID:        1,
			Title:     "Title 1",
			Content:   "Content 1",
			Completed: false,
		},
		{
			ID:        2,
			Title:     "Title 2",
			Content:   "Content 2",
			Completed: true,
		},
		{
			ID:        3,
			Title:     "Title 3",
			Content:   "Content 3",
			Completed: true,
		},
	}
	mock := mockTodoService{
		todos: todos,
		err:   nil,
	}
	req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := ListTodos(&mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("期望得到code%v,实际得到code%v", http.StatusOK, rec.Code)
	}
}
func TestListTodo_ErrCases(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "内部错误",
			err:      errors.New("内部错误"),
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockTodoService{
				err: tt.err,
			}
			req := httptest.NewRequest(http.MethodGet, "/api/todos", nil)
			ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			h := ListTodos(&mock)
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantCode {
				t.Errorf("期望得到code %v,实际得到%v", tt.wantCode, rec.Code)
			}
		})
	}
}
func TestCreateTodo(t *testing.T) {
	mock := mockTodoService{
		todo: model.Todo{
			ID:        1,
			Title:     "TestTitle",
			Content:   "TestContent",
			Completed: false,
		},
	}
	todo := model.Todo{
		ID:        1,
		Title:     "TestTitle",
		Content:   "TestContent",
		Completed: false,
	}
	body, err := json.Marshal(todo)
	if err != nil {
		t.Errorf("json编码错误%v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/todos", bytes.NewReader(body))
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := CreateTodo(&mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("期望code %v,实际code %v", http.StatusCreated, rec.Code)
	}
	recBody := rec.Result().Body
	defer recBody.Close()
}
func TestCreateTodo_ErrCases(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "标题为空",
			err:      model.ErrEmptyContent,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "内部错误",
			err:      errors.New("内部错误"),
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockTodoService{
				err: tt.err,
			}
			body := `{}`
			req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(body))
			ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			h := CreateTodo(&mock)
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantCode {
				t.Errorf("期望得到code%v,实际得到code%v", tt.wantCode, rec.Code)
			}
			if mock.called != 0 {
				t.Errorf("期望调用service 0次 实际调用%v", mock.called)
			}
		})
	}
}
func TestCreateTodo_InvalidJson(t *testing.T) {
	mock := mockTodoService{}
	body := `{"title":}`
	req := httptest.NewRequest(http.MethodPost, "/api/todos", strings.NewReader(body))
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := CreateTodo(&mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("期望得到code%v,实际得到code%v", http.StatusBadRequest, rec.Code)
	}
}
func TestGetTodo(t *testing.T) {
	mock := mockTodoService{
		todo: model.Todo{
			ID:        1,
			Title:     "TestTitle",
			Content:   "TestContent",
			Completed: false,
		},
	}
	req := httptest.NewRequest(http.MethodGet, "/api/todos/1", nil)
	req.SetPathValue("id", "1")
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := GetTodo(&mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("期望得到code %v,实际得到code %v", http.StatusOK, rec.Code)
	}
	recBody := rec.Result().Body
	defer recBody.Close()
	var recTodo model.Todo
	if err := json.NewDecoder(recBody).Decode(&recTodo); err != nil {
		t.Errorf("解析json失败%v", err)
	}
	if mock.passID != 1 {
		t.Errorf("期望传递ID%v,实际传递%v", 1, mock.passID)
	}
	if recTodo.ID != 1 {
		t.Errorf("期望得到ID%v,实际得到%v", 1, recTodo.ID)
	}
	if recTodo.Title != "TestTitle" {
		t.Errorf("期望得到Title %v,实际得到%v", "TestTitle", recTodo.Title)
	}
	if recTodo.Content != "TestContent" {
		t.Errorf("期望得到Content %v,实际得到%v", "TestContent", recTodo.Content)
	}
	if recTodo.Completed != false {
		t.Errorf("期望Completed %v，实际得到%v", "false", recTodo.Completed)
	}
}
func TestGetTodo_ErrCases(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "未查询到信息",
			err:      model.ErrTodoNotFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "内部错误",
			err:      errors.New("内部错误"),
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockTodoService{
				err: tt.err,
			}
			req := httptest.NewRequest(http.MethodGet, "/api/todos/1", nil)
			req.SetPathValue("id", "1")
			ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			h := GetTodo(&mock)
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantCode {
				t.Errorf("期望得到code%v,实际得到code%v", tt.wantCode, rec.Code)
			}
		})
	}
}
func TestGetTodo_InvalidPathValue(t *testing.T) {
	mock := mockTodoService{}
	req := httptest.NewRequest(http.MethodGet, "/api/todos/1", nil)
	req.SetPathValue("id", "WrongValue")
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := GetTodo(&mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("期望得到code %v,实际得到code %v", http.StatusBadRequest, rec.Code)
	}
}
func TestDeleteTodo(t *testing.T) {
	mock := mockTodoService{
		todo: model.Todo{
			ID:        1,
			Title:     "TestTitle",
			Content:   "TestContent",
			Completed: false,
		},
	}
	req := httptest.NewRequest(http.MethodDelete, "/api/todos/1", nil)
	req.SetPathValue("id", "1")
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := DeleteTodo(&mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Errorf("期望得到code %v,实际得到code %v", http.StatusNoContent, rec.Code)
	}
	if mock.passID != 1 {
		t.Errorf("期望传递ID %v,实际传递%v", 1, mock.passID)
	}
}
func TestDeleteTodo_ErrCases(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "内容不存在",
			err:      model.ErrTodoNotFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "未认证",
			err:      model.ErrUnauthorized,
			wantCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockTodoService{
				err: tt.err,
			}
			req := httptest.NewRequest(http.MethodDelete, "/api/todos/1", nil)
			req.SetPathValue("id", "1")
			ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			h := DeleteTodo(&mock)
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantCode {
				t.Errorf("期望得到 code%v,实际得到%v", tt.wantCode, rec.Code)
			}

		})
	}
}
func TestDeleteTodo_InvalidPathValue(t *testing.T) {
	mock := mockTodoService{}
	req := httptest.NewRequest(http.MethodDelete, "/api/todos/1", nil)
	req.SetPathValue("id", "WrongValue")
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := DeleteTodo(&mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("期望得到 code%v,实际得到%v", http.StatusBadRequest, rec.Code)
	}
}
func TestUpdateTodo(t *testing.T) {
	mock := mockTodoService{
		todo: model.Todo{
			ID:        1,
			Title:     "TestTitle",
			Content:   "TestContent",
			Completed: false,
		},
	}
	body := `{"id":1,"title":"NewTestTitle","content":"NewTestContent","completed":true}`
	req := httptest.NewRequest(http.MethodPut, "/api/todos/1", strings.NewReader(body))
	req.SetPathValue("id", "1")
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := UpdateTodo(&mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Errorf("期望得到code%v,实际得到%v", http.StatusNoContent, rec.Code)
	}

	if mock.todo.ID != 1 {
		t.Errorf("期望得到ID %v,实际得到%v", 1, mock.todo.ID)
	}
	if mock.todo.Title != "NewTestTitle" {
		t.Errorf("期望得到NewTitle %v，实际得到%v", "NewTestTitle", mock.todo.Title)
	}
	if mock.todo.Content != "NewTestContent" {
		t.Errorf("期望得到NewContent %v,实际得到%v", "NewTestContent", mock.todo.Content)
	}
	if mock.todo.Completed != true {
		t.Errorf("期望得到NewCompleted %v，实际得到%v", "NewTestContent", mock.todo.Completed)
	}
}
func TestUpdateTodo_ErrCases(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "内容不存在",
			err:      model.ErrTodoNotFound,
			wantCode: http.StatusNotFound,
		},
		{
			name:     "内部错误",
			err:      errors.New("内部错误"),
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := mockTodoService{
				err: tt.err,
			}
			req := httptest.NewRequest(http.MethodPost, "/api/todos/1", nil)
			req.SetPathValue("id", "1")
			ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			h := DeleteTodo(&mock)
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantCode {
				t.Errorf("期望得到 code%v,实际得到%v", tt.wantCode, rec.Code)
			}

		})
	}
}
func TestUpdateTodo_InvalidJson(t *testing.T) {
	mock := mockTodoService{}
	body := `{"title":}`
	req := httptest.NewRequest(http.MethodPost, "/api/todos/1", strings.NewReader(body))
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := CreateTodo(&mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("期望得到code%v,实际得到code%v", http.StatusBadRequest, rec.Code)
	}
}
func TestUpdateTodo_InvalidPathValue(t *testing.T) {
	mock := mockTodoService{}
	req := httptest.NewRequest(http.MethodPost, "/api/todos/1", nil)
	req.SetPathValue("id", "WrongValue")
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := DeleteTodo(&mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("期望得到 code%v,实际得到%v", http.StatusBadRequest, rec.Code)
	}
}
