package handler

import (
	"GO_web/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type MockService struct {
	User           model.User
	Session        string
	Err            error
	LoginCalled    int
	RegisterCalled int
	LogoutCalled   int
}

func (s *MockService) Login(ctx context.Context,
	username string,
	password string) (model.User, string, error) {
	s.LoginCalled++
	return s.User, s.Session, s.Err
}
func (s *MockService) Register(ctx context.Context,
	username string,
	password string) (model.User, string, error) {
	s.RegisterCalled++
	return s.User, s.Session, s.Err
}
func (s *MockService) Logout(ctx context.Context,
	session string) error {
	s.LogoutCalled++
	s.Session = session
	return s.Err
}
func TestLogin_ErrCases(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantStatusCode int
	}{
		{
			name:           "用户名为空",
			err:            model.ErrEmptyUsername,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "密码为空",
			err:            model.ErrEmptyPassword,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "用户名或密码错误",
			err:            model.ErrInvalidCredentials,
			wantStatusCode: http.StatusUnauthorized,
		},
		{
			name:           "内部错误",
			err:            errors.New("内部错误"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockService{
				Err: tt.err,
			}
			body := `{"username":"test","password":"test"}`
			req := httptest.NewRequest(
				http.MethodPost,
				"/api/login",
				strings.NewReader(body),
			)
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h := LoginTodo(mock)
			h.ServeHTTP(rec, req)

			if rec.Code != tt.wantStatusCode {
				t.Errorf("LoginTodo status code = %v, wantStatusCode %v", rec.Code, tt.wantStatusCode)
			}
			if len(rec.Result().Cookies()) != 0 {
				t.Error("登录失败时不应该设置 Cookie")
			}
		})
	}
}
func TestLogin_invalidJSON(t *testing.T) {
	mock := &MockService{}
	body := `{"username":"test","password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h := LoginTodo(mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("期望LoginTodo = %v,实际 %v", http.StatusBadRequest, rec.Code)
	}
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

func TestRegister_ErrCases(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		wantStatusCode int
	}{
		{
			name:           "用户名为空",
			err:            model.ErrEmptyUsername,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "密码为空",
			err:            model.ErrEmptyPassword,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "用户名重复",
			err:            model.ErrDuplicateUsername,
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "内部错误",
			err:            errors.New("错误"),
			wantStatusCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &MockService{
				Err: tt.err,
			}
			body := `{"username":"test","password":"test"}`
			req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			h := RegisterTodo(mock)
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantStatusCode {
				t.Errorf("期望 code = %v, 得到Code %v", tt.wantStatusCode, rec.Code)

			}
		})
	}
}
func TestRegister_invalidJSON(t *testing.T) {
	mock := &MockService{}
	body := `{"username":"test","password"}`
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h := RegisterTodo(mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("期望code %v，实际得到code %v", http.StatusBadRequest, rec.Code)
	}
}
func TestRegister(t *testing.T) {
	mock := &MockService{
		User: model.User{
			ID:           1,
			Username:     "test",
			PasswordHash: "this-is-a-secret-hash",
			Profile:      "",
			CreatedAt:    time.Now(),
		},
		Session: "test-session",
		Err:     nil,
	}
	body := `{"username":"test","password":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h := RegisterTodo(mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusCreated {
		t.Errorf("期望code %v,实际code %v", http.StatusCreated, rec.Code)
	}
	var sessionCookie *http.Cookie
	for _, cookie := range rec.Result().Cookies() {
		if cookie.Name == "session_id" {
			sessionCookie = cookie
			break
		}
	}
	if sessionCookie == nil {
		t.Errorf("注册后未设置Cookie")
	}
	if sessionCookie.Value != "test-session" {
		t.Errorf("Cookie 设置错误")
	}
	if !sessionCookie.HttpOnly {
		t.Error("session_id Cookie 应设置 HttpOnly")
	}
	if sessionCookie.MaxAge != 600 {
		t.Errorf("Cookie MaxAge 错误：%d", sessionCookie.MaxAge)
	}
	defer rec.Result().Body.Close()
	recBody := rec.Result().Body
	ResUser := model.User{}
	if err := json.NewDecoder(recBody).Decode(&ResUser); err != nil {
		t.Fatalf("解析响应JSON失败：%v", err)
	}
	if ResUser.Username != "test" {
		t.Errorf("期望用户名 test,实际用户名 %s ", ResUser.Username)
	}
	if ResUser.PasswordHash != "" {
		t.Error("密码哈希错误返回")
	}
	if ResUser.Profile != "" {
		t.Errorf("期望简介为空 实际得到%s", ResUser.Profile)
	}
}

func TestLogout(t *testing.T) {
	mock := &MockService{}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/logout",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session",
	})

	rec := httptest.NewRecorder()

	h := LogoutTodo(mock)
	h.ServeHTTP(rec, req)

	// 检查状态码
	if rec.Code != http.StatusNoContent {
		t.Errorf(
			"期望状态码 %d，实际得到 %d",
			http.StatusNoContent,
			rec.Code,
		)
	}

	// 检查 Service 是否被调用
	if mock.LogoutCalled != 1 {
		t.Errorf(
			"期望 Logout 调用 1 次，实际调用 %d 次",
			mock.LogoutCalled,
		)
	}
	if mock.Session != "test-session" {
		t.Errorf(
			"期望 Session ID 为 test-session，实际得到 %s",
			mock.Session,
		)
	}
	cookies := rec.Result().Cookies()

	var found bool
	for _, cookie := range cookies {
		if cookie.Name == "session_id" {
			found = true

			if cookie.MaxAge >= 0 {
				t.Error("退出登录后，Cookie 应该被删除")
			}
		}
	}

	if !found {
		t.Error("响应中没有清除 session_id 的 Cookie")
	}
}
func TestLogout_NoCookie(t *testing.T) {
	mock := &MockService{}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/logout",
		nil,
	)

	rec := httptest.NewRecorder()

	h := LogoutTodo(mock)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf(
			"期望状态码 %d，实际得到 %d",
			http.StatusUnauthorized,
			rec.Code,
		)
	}

	if mock.LogoutCalled != 0 {
		t.Error("没有 Cookie 时不应该调用 Logout Service")
	}
}
func TestLogout_ServiceError(t *testing.T) {
	mock := &MockService{
		Err: errors.New("session 删除失败"),
	}

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/logout",
		nil,
	)

	req.AddCookie(&http.Cookie{
		Name:  "session_id",
		Value: "test-session",
	})

	rec := httptest.NewRecorder()

	h := LogoutTodo(mock)
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf(
			"期望状态码 %d，实际得到 %d",
			http.StatusInternalServerError,
			rec.Code,
		)
	}
}
