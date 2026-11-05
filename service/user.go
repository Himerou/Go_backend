package service

import (
	"GO_web/model"
	"GO_web/repository"
	"context"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	user *repository.UserRepository
}

func NewUserService(user *repository.UserRepository) *UserService {
	return &UserService{user: user}
}

type GetProfileService interface {
	GetProfile(ctx context.Context) (model.PublicUserProfile, error)
}

func (s *UserService) GetProfile(
	ctx context.Context) (model.PublicUserProfile, error) {
	user, err := s.user.Get(ctx)
	if err != nil {
		return model.PublicUserProfile{}, err
	}
	profile := model.PublicUserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
	}
	return profile, nil
}

type GetPublicProfileService interface {
	GetPublicProfile(ctx context.Context, username string) (model.PublicUserProfile, error)
}

func (s *UserService) GetPublicProfile(
	ctx context.Context,
	username string) (model.PublicUserProfile, error) {
	if username == "" {
		return model.PublicUserProfile{}, model.ErrEmptyUsername
	}
	user, err := s.user.GetByUsername(ctx, username)
	if err != nil {
		return model.PublicUserProfile{}, err
	}
	profile := model.PublicUserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
	}
	return profile, nil
}

type UpdateProfileService interface {
	UpdateProfile(ctx context.Context, msg string) (model.PublicUserProfile, error)
}

func (s *UserService) UpdateProfile(
	ctx context.Context,
	msg string,
) (model.PublicUserProfile, error) {
	err := s.user.UpdateProfile(ctx, msg)
	if err != nil {
		return model.PublicUserProfile{}, err
	}
	user, err := s.user.Get(ctx)
	if err != nil {
		return model.PublicUserProfile{}, err
	}
	profile := model.PublicUserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
	}
	return profile, nil
}

type UpdateUsernameService interface {
	UpdateUsername(ctx context.Context, msg string) (model.PublicUserProfile, error)
}

func (s *UserService) UpdateUsername(
	ctx context.Context,
	msg string) (model.PublicUserProfile, error) {
	if msg == "" {
		return model.PublicUserProfile{}, model.ErrEmptyContent
	}
	_, err := s.user.GetByUsername(ctx, msg)
	if err == nil {
		return model.PublicUserProfile{}, model.ErrDuplicateUsername
	}
	err = s.user.UpdateUsername(ctx, msg)
	if err != nil {
		return model.PublicUserProfile{}, err
	}
	user, err := s.user.GetByUsername(ctx, msg)
	if err != nil {
		return model.PublicUserProfile{}, err
	}
	profile := model.PublicUserProfile{
		ID:        user.ID,
		Username:  user.Username,
		Profile:   user.Profile,
		CreatedAt: user.CreatedAt,
	}
	return profile, nil
}

type UpdatePasswordService interface {
	UpdatePassword(ctx context.Context, msg string) error
}

func (s *UserService) UpdatePassword(
	ctx context.Context,
	msg string) error {
	if strings.TrimSpace(msg) == "" {
		return model.ErrEmptyPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(msg), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.user.UpdatePassword(ctx, string(hash))
	if err != nil {
		return err
	}
	return nil
}

type DeleteUserService interface {
	DeleteUser(ctx context.Context) error
}

func (s *UserService) DeleteUser(ctx context.Context) error {
	err := s.user.Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}
