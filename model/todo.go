package model

import (
	"errors"
	"time"
)

type contextKey string

const UserIDKey contextKey = "userID"

type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Profile      string    `json:"profile"`
	CreatedAt    time.Time `json:"created_at"`
}
type PublicUserProfile struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Profile   string    `json:"profile"`
	CreatedAt time.Time `json:"created_at"`
}
type LoginAndRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Msg struct {
	Msg string `json:"msg"`
}

var ErrTodoNotFound = errors.New("未找到对应数据")
var ErrDuplicateUsername = errors.New("用户名重复")
var ErrEmptyUsername = errors.New("用户名为空")
var ErrEmptyPassword = errors.New("密码为空")
var ErrInvalidCredentials = errors.New("用户名或密码错误")
var ErrUnauthorized = errors.New("身份未验证")
var ErrUserNotFound = errors.New("未查询到用户数据")
var ErrEmptyContent = errors.New("提交为空")
var ErrDuplicatedPassword = errors.New("新密码不得与旧密码重复")
var ErrInternalServiceError = errors.New("内部错误")
