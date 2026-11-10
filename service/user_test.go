package service

import (
	"GO_web/model"
	"context"
	"errors"
	"testing"
	"time"
)

func TestGetProfile(t *testing.T) {
	mock := &mockUserRepo{
		user: model.User{
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
			CreatedAt:    time.Now(),
		},
	}
	service := NewUserService(mock)
	ctx := context.Background()
	recUser, err := service.GetProfile(ctx)
	if err != nil {
		t.Errorf("意外错误 %v", err)
	}
	if recUser.Username != mock.user.Username {
		t.Errorf("期望Username %v，实际得到%v", "TestUsername", recUser.Username)
	}
	if recUser.Profile != mock.user.Profile {
		t.Errorf("期望得到 %v ,实际得到%v", "TestProfile", recUser.Profile)
	}
}
func TestGetProfile_ErrCases(t *testing.T) {
	tests := []struct {
		name         string
		Username     string
		PasswordHash string
		Profile      string
		err          error
		wantErr      error
	}{
		{
			name:         "未认证",
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
			err:          model.ErrUserNotFound,
			wantErr:      model.ErrUserNotFound,
		},
		{
			name:         "未查询到数据",
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
			err:          model.ErrUserNotFound,
			wantErr:      model.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{
				user: model.User{
					Username:     tt.Username,
					PasswordHash: tt.PasswordHash,
					Profile:      tt.Profile,
				},
				err: tt.err,
			}
			service := NewUserService(mock)
			ctx := context.Background()
			_, err := service.GetProfile(ctx)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("期望得到错误 %v,实际得到 %v", tt.wantErr, err)
			}

		})
	}
}
func TestGetPublicProfile(t *testing.T) {
	mock := &mockUserRepo{
		user: model.User{
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
		},
	}
	service := NewUserService(mock)
	ctx := context.Background()
	recProfile, err := service.GetPublicProfile(ctx, "TestUsername")
	if err != nil {
		t.Errorf("意外的错误 %v", err)
	}
	if recProfile.Username != mock.user.Username {
		t.Errorf("期望得到 Username %v,实际得到 %v", mock.user.Username, recProfile.Username)
	}
	if recProfile.Profile != mock.user.Profile {
		t.Errorf("期望得到 Profile %v,实际得到 %v", mock.user.Profile, recProfile.Profile)
	}
}
func TestGetPublicProfile_ErrCases(t *testing.T) {
	var errInternalError = errors.New("内部错误")
	tests := []struct {
		name         string
		pathValue    string
		Username     string
		PasswordHash string
		Profile      string
		err          error
		wantErr      error
	}{
		{
			name:         "参数为空",
			pathValue:    "",
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
			err:          model.ErrEmptyUsername,
			wantErr:      model.ErrEmptyUsername,
		},
		{
			name:         "内部错误",
			pathValue:    "TestUsername",
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
			err:          errInternalError,
			wantErr:      errInternalError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{
				user: model.User{
					Username:     tt.Username,
					PasswordHash: tt.PasswordHash,
					Profile:      tt.Profile,
				},
				err: tt.err,
			}
			service := NewUserService(mock)
			ctx := context.Background()
			_, err := service.GetPublicProfile(ctx, tt.pathValue)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("期望得到错误 %v,实际得到 %v", tt.wantErr, err)
			}
		})
	}
}
func TestUpdateProfile(t *testing.T) {
	mock := &mockUserRepo{
		user: model.User{
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
		},
	}
	msg := "NewTestProfile"
	service := NewUserService(mock)
	ctx := context.Background()
	recData, err := service.UpdateProfile(ctx, msg)
	if err != nil {
		t.Errorf("意外的错误 %v", err)
	}
	if recData.Username != mock.user.Username {
		t.Errorf("期望得到Username %v,实际得到 %v", mock.user.Username, recData.Username)
	}
	if recData.Profile != msg {
		t.Errorf("期望更新Profile %v,实际得到 %v", msg, recData.Profile)
	}
}
func TestUpdateProfile_ErrCases(t *testing.T) {
	var ErrInternalServiceError = errors.New("内部错误")
	tests := []struct {
		name         string
		msg          string
		Username     string
		PasswordHash string
		Profile      string
		err          error
		wantErr      error
	}{
		{
			name:         "更新内容为空",
			msg:          "",
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
			err:          model.ErrEmptyContent,
			wantErr:      model.ErrEmptyContent,
		},
		{
			name:         "内部错误",
			msg:          "NewTestProfile",
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
			err:          ErrInternalServiceError,
			wantErr:      ErrInternalServiceError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{
				user: model.User{
					Username:     tt.Username,
					PasswordHash: tt.PasswordHash,
					Profile:      tt.Profile,
				},
				err: tt.err,
			}
			service := NewUserService(mock)
			ctx := context.Background()
			_, err := service.UpdateProfile(ctx, tt.msg)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("期望得到错误 %v,实际得到 %v", tt.wantErr, err)
			}

		})
	}
}
func TestUpdateUsername(t *testing.T) {
	mock := &mockUserRepo{
		user: model.User{
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
		},
	}
	msg := "NewTestUsername"
	service := NewUserService(mock)
	ctx := context.Background()
	recData, err := service.UpdateUsername(ctx, msg)
	if err != nil {
		t.Errorf("意外的错误%v", err)
	}
	if recData.Username != msg {
		t.Errorf("期望更新Username %v，实际得到 %v", msg, recData.Username)
	}
	if recData.Profile != mock.user.Profile {
		t.Errorf("期望得到Profile %v,实际得到%v", mock.user.Profile, recData.Profile)
	}
}
func TestUpdateUsername_ErrCases(t *testing.T) {
	var ErrInternalServiceError = errors.New("内部错误")
	tests := []struct {
		name     string
		msg      string
		Username string
		Password string
		Profile  string
		err      error
		wantErr  error
	}{
		{
			name:     "空参数",
			msg:      "",
			Username: "TestUsername",
			Password: "TestPasswordHash",
			Profile:  "TestProfile",
			err:      model.ErrEmptyContent,
			wantErr:  model.ErrEmptyContent,
		},
		{
			name:     "内部错误",
			msg:      "NewTestUsername",
			Username: "TestUsername",
			Password: "TestPasswordHash",
			Profile:  "TestProfile",
			err:      ErrInternalServiceError,
			wantErr:  ErrInternalServiceError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{
				user: model.User{
					Username:     tt.Username,
					PasswordHash: tt.Password,
					Profile:      tt.Profile,
				},
				err: tt.err,
			}
			service := NewUserService(mock)
			ctx := context.Background()
			_, err := service.UpdateUsername(ctx, tt.msg)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("期望得到错误 %v,实际得到 %v", tt.wantErr, err)
			}
		})
	}
}
func TestUpdatePassword(t *testing.T) {
	mock := &mockUserRepo{
		user: model.User{
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
		},
	}
	msg := "NewTestPassword"
	service := NewUserService(mock)
	ctx := context.Background()
	err := service.UpdatePassword(ctx, msg)
	if err != nil {
		t.Errorf("意外的错误%v", err)
	}
}
func TestUpdatePassword_ErrCases(t *testing.T) {
	tests := []struct {
		name     string
		msg      string
		Username string
		Password string
		Profile  string
		err      error
		wantErr  error
	}{
		{
			name:     "重复密码",
			msg:      "TestPassword",
			Username: "TestUsername",
			Password: "TestPasswordHash",
			Profile:  "TestProfile",
			err:      model.ErrDuplicatedPassword,
			wantErr:  model.ErrDuplicatedPassword,
		},
		{
			name:     "内部错误",
			msg:      "NewTestPassword",
			Username: "TestUsername",
			Password: "TestPasswordHash",
			Profile:  "TestProfile",
			err:      model.ErrInternalServiceError,
			wantErr:  model.ErrInternalServiceError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mock := &mockUserRepo{
				user: model.User{
					Username:     "TestUsername",
					PasswordHash: "TestPasswordHash",
					Profile:      "TestProfile",
				},
				err: tt.err,
			}
			service := NewUserService(mock)
			ctx := context.Background()
			err := service.UpdatePassword(ctx, tt.msg)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("期望得到错误 %v,实际得到%v", tt.wantErr, err)
			}
		})
	}
}
func TestDeleteUser(t *testing.T) {
	mock := &mockUserRepo{
		user: model.User{
			Username:     "TestUsername",
			PasswordHash: "TestPasswordHash",
			Profile:      "TestProfile",
		},
	}
	ctx := context.Background()
	service := NewUserService(mock)
	err := service.DeleteUser(ctx)
	if err != nil {
		t.Errorf("意外的错误%v", err)
	}
}
