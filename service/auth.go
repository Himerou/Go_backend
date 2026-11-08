package service

import (
	"GO_web/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	Get(ctx context.Context) (model.User, error)
	GetByUsername(ctx context.Context, username string) (model.User, error)
	Delete(ctx context.Context) error
	Update(ctx context.Context, user model.User) error
	UpdatePassword(ctx context.Context, password string) error
	UpdateProfile(ctx context.Context, msg string) error
	UpdateUsername(ctx context.Context, name string) error
}
type SessionRepository interface {
	Create(userID int) (sessionID string, err error)
	Delete(sessionID string) error
	Get(sessionID string) (userID int, ok bool)
}
type AuthService struct {
	user    UserRepository
	session SessionRepository
}

func NewAuthService(user UserRepository, session SessionRepository) *AuthService {
	return &AuthService{
		user:    user,
		session: session}
}

type LoginService interface {
	Login(ctx context.Context,
		username string,
		password string) (model.User, string, error)
}
type RegisterService interface {
	Register(ctx context.Context,
		username string,
		password string) (model.User, string, error)
}

type LogoutService interface {
	Logout(ctx context.Context, sessionID string) error
}

func (s *AuthService) Register(
	ctx context.Context,
	username string,
	password string,
) (model.User, string, error) {
	var user model.User
	username = strings.TrimSpace(username)
	if username == "" {
		return model.User{}, "", model.ErrEmptyUsername
	}
	if strings.TrimSpace(password) == "" {
		return model.User{}, "", model.ErrEmptyPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, "", fmt.Errorf("哈希加密失败 %w", err)
	}
	user.PasswordHash = string(hash)
	user.Username = username

	user, err = s.user.Create(ctx, user)
	if err != nil {
		return model.User{}, "", err
	}
	session, err := s.session.Create(user.ID)
	if err != nil {
		return model.User{}, "", err
	}
	return user, session, nil
}

func (s *AuthService) Login(ctx context.Context,
	username string,
	password string) (model.User, string, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return model.User{}, "", model.ErrEmptyUsername
	}
	if strings.TrimSpace(password) == "" {
		return model.User{}, "", model.ErrEmptyPassword
	}
	user, err := s.user.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, "", model.ErrInvalidCredentials
		}
		return model.User{}, "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return model.User{}, "", model.ErrInvalidCredentials
		}
		return model.User{}, "", err
	}
	session, err := s.session.Create(user.ID)
	if err != nil {
		return model.User{}, "", err
	}
	return user, session, nil
}
func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	err := s.session.Delete(sessionID)
	if err != nil {
		return err
	}
	return nil
}
