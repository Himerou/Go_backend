package repository

import (
	"GO_web/model"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-sql-driver/mysql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user model.User) (model.User, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO users(username, password_hash,profile)
		VALUES (?,?,?)`,
		user.Username,
		user.PasswordHash,
		user.Profile,
	)
	if err != nil {
		var mysqlerr *mysql.MySQLError
		if errors.As(err, &mysqlerr) && mysqlerr.Number == 1062 {
			return model.User{}, model.ErrDuplicateUsername
		}
		return model.User{}, fmt.Errorf("创建用户失败 %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.User{}, fmt.Errorf("获取用户ID失败%w", err)
	}
	user.ID = int(id)

	err = r.db.QueryRowContext(ctx, `
				SELECT created_at 
				FROM users
				WHERE id = ?`,
		user.ID).Scan(&user.CreatedAt)
	if err != nil {
		return model.User{}, fmt.Errorf("获取创建时间失败%w", err)
	}

	return user, nil
}
func (r *UserRepository) Get(ctx context.Context) (model.User, error) {
	userID, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return model.User{}, model.ErrUserNotFound
	}
	var user model.User
	err := r.db.QueryRowContext(ctx, `
SELECT id, username, password_hash, profile, created_at
FROM users
WHERE id = ?`,
		userID,
	).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Profile, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
	}
	return user, nil
}
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (model.User, error) {
	var user model.User
	err := r.db.QueryRowContext(ctx, `
		SELECT id,username,password_hash,profile,created_at
		FROM users 
		WHERE username =?`, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Profile, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, model.ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("查询用户失败 %w", err)
	}
	return user, nil

}

func (r *UserRepository) Delete(ctx context.Context) error {
	id, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return model.ErrUnauthorized
	}
	result, err := r.db.ExecContext(ctx, `
					DELETE FROM users
					WHERE id =?`,
		id,
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) Update(ctx context.Context, user model.User) error {
	id, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return model.ErrUnauthorized
	}
	result, err := r.db.ExecContext(ctx, `
UPDATE users
SET username=?,password_hash=?,profile=?,
WHERE id=?`,
		user.Username,
		user.PasswordHash,
		user.Profile,
		id,
	)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrUserNotFound
	}
	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, password string) error {
	userID, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return model.ErrUnauthorized
	}
	_, err := r.db.ExecContext(ctx, `
UPDATE users
SET password_hash=?
WHERE id =?`,
		password,
		userID)
	if err != nil {
		return err
	}
	return nil
}
func (r *UserRepository) UpdateProfile(ctx context.Context, msg string) error {
	userID, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return model.ErrUnauthorized
	}
	_, err := r.db.ExecContext(ctx, `
UPDATE users
SET profile=?
WHERE id =?`,
		msg,
		userID)
	if err != nil {
		return err
	}

	return nil
}
func (r *UserRepository) UpdateUsername(ctx context.Context, name string) error {
	userID, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return model.ErrUnauthorized
	}
	result, err := r.db.ExecContext(ctx, `
UPDATE users
SET username=?
WHERE id =?`,
		name,
		userID)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrUserNotFound
	}
	return nil
}
