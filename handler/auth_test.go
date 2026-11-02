package handler

import (
	"GO_web/model"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockService struct {
	User    model.User
	Session string
	Err     error
}

func (s *MockService) Login(ctx context.Context,
	username string,
	password string) (model.User, string, error) {
	return s.User, s.Session, s.Err
}
func TestLogin(t *testing.T) {
	mock := &MockService{
		User: model.User{
			ID:           1,
			Username:     "test",
			PasswordHash: "",
			Profile:      "",
			CreatedAt:    time.Now(),
		},
		Session: "test-session",
		Err:     nil,
	}
	body := `{"username":"test",
			"password":"123456"
			}`
	req := httptest.NewRequest("POST", "/api/login", strings.NewReader(body))
	rec := httptest.NewRecorder()
	h := LoginTodo(mock)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Login failed, expected %d, got %d", http.StatusOK, rec.Code)
	}
	result := rec.Result()
	defer result.Body.Close()
	var sessionCookie *http.Cookie
	for _, cookie := range result.Cookies() {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}
	if sessionCookie == nil {
		t.Errorf("Login failed, expected cookie not found")
	}
	if sessionCookie.Value != "test-session" {
		t.Errorf("Login failed, expected session cookie value, got %s", sessionCookie.Value)
	}
	if !sessionCookie.HttpOnly {
		t.Error("session_id Cookie 应设置 HttpOnly")
	}
	if sessionCookie.MaxAge != 600 {
		t.Errorf("Cookie MaxAge 错误：%d", sessionCookie.MaxAge)
	}
}
