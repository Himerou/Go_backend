package handler

import (
	"GO_web/model"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type mockService struct {
	User           model.User
	PublicUser     model.PublicUserProfile
	Err            error
	session        string
	called         int
	calledUsername string
	calledPassword string
}

func (s *mockService) GetProfile(ctx context.Context) (model.PublicUserProfile, error) {
	s.called++
	return s.PublicUser, s.Err
}
func (s *mockService) GetPublicProfile(ctx context.Context, username string) (model.PublicUserProfile, error) {
	s.called++
	s.calledUsername = username
	return s.PublicUser, s.Err
}
func (s *mockService) UpdateProfile(ctx context.Context, msg string) (model.PublicUserProfile, error) {
	s.PublicUser.Profile = msg
	return s.PublicUser, s.Err
}
func (s *mockService) UpdateUsername(ctx context.Context, msg string) (model.PublicUserProfile, error) {
	s.PublicUser.Username = msg
	return s.PublicUser, s.Err
}
func (s *mockService) UpdatePassword(ctx context.Context, msg string) error {
	s.calledPassword = msg
	return s.Err
}
func (s *mockService) DeleteUser(ctx context.Context) error {
	return s.Err
}
func TestGetProfile(t *testing.T) {
	mock := &mockService{
		User: model.User{},
		PublicUser: model.PublicUserProfile{
			ID:        1,
			Username:  "test",
			Profile:   "test_profile",
			CreatedAt: time.Now(),
		},
		Err:     nil,
		session: "test_session",
	}
	req := httptest.NewRequest(http.MethodGet, "/api/profile/", nil)
	ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
	req = req.WithContext(ctx)
	rec := httptest.NewRecorder()
	h := GetProfile(mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("期望得到状态吗 %v,实际得到 %v", http.StatusOK, rec.Code)
	}
	var recData model.PublicUserProfile
	if err := json.NewDecoder(rec.Body).Decode(&recData); err != nil {
		t.Errorf("json 解析失败 %v", err)
	}
	if recData.Username != mock.PublicUser.Username {
		t.Errorf("期望用户名 %v,实际得到 %v", mock.User.Username, recData.Username)
	}
	if recData.Profile != mock.PublicUser.Profile {
		t.Errorf("期望得到简介内容 %v,实际得到 %v", mock.User.Profile, recData.Profile)
	}

}

func TestGetProfile_ErrCases(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "未认证",
			err:      model.ErrUnauthorized,
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "内部错误",
			err:      errors.New("内部错误"),
			wantCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockService{
				Err: tt.err,
			}
			req := httptest.NewRequest(http.MethodGet, "/api/profile/", nil)
			ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
			req.WithContext(ctx)
			rec := httptest.NewRecorder()
			h := GetProfile(mock)
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantCode {
				t.Errorf("期望code %v,实际code %v", tt.wantCode, rec.Code)
			}

		})
	}
}

func TestGetPublicProfile(t *testing.T) {
	mock := &mockService{
		User: model.User{},
		PublicUser: model.PublicUserProfile{
			ID:       1,
			Username: "testUsername",
			Profile:  "test_profile",
		},
		Err:     nil,
		session: "test_session",
	}
	req := httptest.NewRequest(http.MethodGet, "/api/users/testUsername/profile/", nil)
	req.SetPathValue("username", "testUsername")
	rec := httptest.NewRecorder()
	h := GetPublicProfile(mock)
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("期望得到 code %v 实际得到 code %v", http.StatusOK, rec.Code)
	}
	var recData model.PublicUserProfile
	if err := json.NewDecoder(rec.Body).Decode(&recData); err != nil {
		t.Errorf("解析json错误 %v", err)
	}
	if recData.Profile != mock.PublicUser.Profile {
		t.Errorf("期望 Profile %v,实际得到 %v", mock.PublicUser.Profile, recData.Profile)
	}
	if recData.Username != mock.PublicUser.Username {
		t.Errorf("期望Username %v,实际 %v", mock.PublicUser.Username, recData.Username)
	}
	if mock.calledUsername != mock.PublicUser.Username {
		t.Errorf("解析参数{username}期望得到%v,实际得到%v", mock.PublicUser.Username, mock.calledUsername)
	}
	if mock.called != 1 {
		t.Errorf("期望 Service 调用 1 次，实际调用 %d 次", mock.called)
	}
}
func TestGetPublicProfile_ErrCases(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
	}{
		{
			name:     "用户不存在",
			err:      model.ErrUserNotFound,
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
			mock := &mockService{
				PublicUser: model.PublicUserProfile{
					ID: 1,
				},
				Err: tt.err,
			}
			req := httptest.NewRequest(http.MethodGet, "/api/users/{testUsername}/profile/", nil)
			ctx := context.WithValue(req.Context(), model.UserIDKey, int64(1))
			req = req.WithContext(ctx)
			rec := httptest.NewRecorder()
			h := GetPublicProfile(mock)
			h.ServeHTTP(rec, req)
			if rec.Code != tt.wantCode {
				t.Errorf("期望code %v,实际得到code %v", tt.wantCode, rec.Code)
			}
		})
	}
}
