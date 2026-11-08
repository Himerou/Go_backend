package service

import (
	"GO_web/model"
	"context"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type mockUserRepo struct {
	user   model.User
	err    error
	called int
}

func (m *mockUserRepo) Create(ctx context.Context, user model.User) (model.User, error) {
	m.called++
	return user, m.err
}

func (m *mockUserRepo) Get(ctx context.Context) (model.User, error) {
	m.called++
	return m.user, m.err
}

func (m *mockUserRepo) GetByUsername(ctx context.Context, username string) (model.User, error) {
	m.called++
	return m.user, m.err
}

func (m *mockUserRepo) Delete(ctx context.Context) error {
	m.called++
	return m.err
}

func (m *mockUserRepo) Update(ctx context.Context, user model.User) error {
	m.called++
	return m.err
}

func (m *mockUserRepo) UpdatePassword(ctx context.Context, password string) error {
	m.called++
	return m.err
}

func (m *mockUserRepo) UpdateProfile(ctx context.Context, msg string) error {
	m.called++
	return m.err
}

func (m *mockUserRepo) UpdateUsername(ctx context.Context, name string) error {
	m.called++
	return m.err
}

type mockSessionRepo struct {
	sessionID string
	userID    int
	ok        bool
	err       error
	called    int
}

func (m *mockSessionRepo) Create(userID int) (string, error) {
	m.called++
	return m.sessionID, m.err
}

func (m *mockSessionRepo) Delete(sessionID string) error {
	m.called++
	return m.err
}

func (m *mockSessionRepo) Get(sessionID string) (int, bool) {
	m.called++
	return m.userID, m.ok
}
func TestRegister(t *testing.T) {
	username := "TestUsername"
	password := "TestPassword"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	mockUser := &mockUserRepo{
		user: model.User{
			Username:     "TestUsername",
			PasswordHash: string(hash),
		},
	}
	mockSession := &mockSessionRepo{
		sessionID: "test_session_id",
	}
	service := NewAuthService(mockUser, mockSession)
	ctx := context.Background()

	recUser, session, err := service.Login(ctx, username, password)
	if err != nil {
		t.Errorf("服务错误 %v", err)
	}
	if recUser.Username != username {
		t.Errorf("期望用户名 %v,实际得到 %v", username, recUser.Username)
	}
	if session != mockSession.sessionID {
		t.Errorf("期望得到session_id %v,实际得到 %v", mockSession.sessionID, session)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(recUser.PasswordHash), []byte(password)); err != nil {
		t.Errorf("hash加密错误 %v", err)
	}
}
func TestRegister_ErrCases(t *testing.T) {
	var (
		errUserRepo    = errors.New("userRepoErr")
		errSessionRepo = errors.New("sessionRepoErr")
	)

	tests := []struct {
		name       string
		username   string
		password   string
		userErr    error
		sessionErr error
		wantErr    error
	}{
		{
			name:     "空名称",
			username: "",
			password: "TestPassword",
			wantErr:  model.ErrEmptyUsername,
		},
		{
			name:     "空密码",
			username: "TestUsername",
			password: "",
			wantErr:  model.ErrEmptyPassword,
		},
		{
			name:     "userRepo错误",
			username: "TestUsername",
			password: "TestPassword",
			userErr:  errUserRepo,
			wantErr:  errUserRepo,
		},
		{
			name:     "用户名重复",
			username: "TestDuplicatedUsername",
			password: "TestPassword",
			userErr:  model.ErrDuplicateUsername,
			wantErr:  model.ErrDuplicateUsername,
		},
		{
			name:       "sessionRepo错误",
			username:   "TestUsername",
			password:   "TestPassword",
			sessionErr: errSessionRepo,
			wantErr:    errSessionRepo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := &mockUserRepo{
				err: tt.userErr,
			}
			mockSession := &mockSessionRepo{
				sessionID: "test_session_id",
				err:       tt.sessionErr,
			}
			service := NewAuthService(mockUser, mockSession)
			ctx := context.Background()
			_, _, err := service.Login(ctx, tt.username, tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("期望得到报错 %v,实际得到%v", tt.wantErr, err)
			}
			if tt.username != "" && tt.password != "" {
				if mockUser.called != 1 {
					t.Errorf("期望调用 UserRepository 一次，实际 %d 次", mockUser.called)
				}
			}
			if tt.userErr != nil && mockSession.called != 0 {
				t.Errorf("UserRepository 出错后不应该创建 Session")
			}
			if tt.sessionErr != nil {
				if mockSession.called != 0 {
					t.Errorf("期望调用 SessionRepository 0次，实际 %d 次", mockSession.called)
				}
			}
		})
	}
}
func TestLogin(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("TestPassword"), bcrypt.DefaultCost)
	mockUser := &mockUserRepo{
		user: model.User{
			Username:     "TestUsername",
			PasswordHash: string(hash),
		},
	}
	mockSession := &mockSessionRepo{
		sessionID: "test_session_id",
	}
	service := NewAuthService(mockUser, mockSession)
	ctx := context.Background()
	recUser, session, err := service.Login(ctx, "TestUsername", "TestPassword")
	if err != nil {
		t.Errorf("得到错误 %v", err)
	}
	if recUser.Username != mockUser.user.Username {
		t.Errorf("期望用户名 %v,实际得到 %v", recUser.Username, recUser.Username)
	}
	if session != mockSession.sessionID {
		t.Errorf("期望得到session_id %v 实际得到 %v", mockSession.sessionID, session)
	}
}
func TestLogin_ErrCases(t *testing.T) {
	var (
		errUserRepo    = errors.New("userRepoErr")
		errSessionRepo = errors.New("sessionRepoErr")
	)
	hash, err := bcrypt.GenerateFromPassword([]byte("TestPassword"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash化失败 %v", err)
	}
	tests := []struct {
		name       string
		username   string
		password   string
		userErr    error
		sessionErr error
		wantErr    error
	}{
		{
			name:     "空用户名",
			username: "",
			password: "TestPassword",
			wantErr:  model.ErrEmptyUsername,
		},
		{
			name:     "空密码",
			username: "TestUsername",
			password: "",
			wantErr:  model.ErrEmptyPassword,
		},
		{
			name:     "密码错误",
			username: "TestUsername",
			password: "TestWrongPassword",
			wantErr:  model.ErrInvalidCredentials,
		},
		{
			name:     "userRepoErr",
			username: "TestUsername",
			password: "TestPassword",
			userErr:  errUserRepo,
			wantErr:  errUserRepo,
		},
		{
			name:       "sessionRepoErr",
			username:   "TestUsername",
			password:   "TestPassword",
			sessionErr: errSessionRepo,
			wantErr:    errSessionRepo,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := &mockUserRepo{
				user: model.User{
					Username:     "TestUsername",
					PasswordHash: string(hash),
				},
				err: tt.userErr,
			}
			mockSession := &mockSessionRepo{
				sessionID: "test_session_id",
				err:       tt.sessionErr,
			}
			service := NewAuthService(mockUser, mockSession)
			ctx := context.Background()
			_, _, err := service.Login(ctx, tt.username, tt.password)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("期望得到错误 %v,实际得到 %v", tt.wantErr, err)
			}
			if tt.username == "" || tt.password == "" {
				if mockUser.called != 0 {
					t.Errorf("参数不合法时不应该调用 UserRepository")
				}
			} else if mockUser.called != 1 {
				t.Errorf("期望调用 UserRepository 1 次，实际调用 %d 次", mockUser.called)
			}
			if tt.sessionErr != nil && mockSession.called != 1 {
				t.Errorf("期望调用 1次,实际调用 %v", mockSession.called)
			}
			if tt.wantErr != nil && mockSession.called != 0 {
				t.Errorf("登录失败时不应该创建 Session")
			}
		})
	}
}
func TestLogout(t *testing.T) {
	mockUser := &mockUserRepo{}
	mockSession := &mockSessionRepo{
		sessionID: "test_session_id",
	}
	ctx := context.Background()
	service := NewAuthService(mockUser, mockSession)
	err := service.Logout(ctx, "test_session_id")
	if err != nil {
		t.Errorf("意外的错误 %v", err)
	}
}
func TestLogout_ErrCases(t *testing.T) {
	var sessionRepoErr = errors.New("sessionRepoErr")
	tests := []struct {
		name      string
		sessionID string
		err       error
		wantErr   error
	}{
		{
			name:      "sessionRepoErr",
			sessionID: "test_session_id",
			err:       sessionRepoErr,
			wantErr:   sessionRepoErr,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := &mockUserRepo{}
			mockSession := &mockSessionRepo{
				sessionID: "test_session_id",
				err:       tt.err,
			}
			service := NewAuthService(mockUser, mockSession)
			ctx := context.Background()
			err := service.Logout(ctx, "test_session_id")
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("期望得到错误 %v，实际得到 %v", tt.wantErr, err)
			}
			if mockSession.called != 1 {
				t.Errorf("期望调用 SessionRepository 1 次，实际 %d 次", mockSession.called)
			}
		})
	}
}
