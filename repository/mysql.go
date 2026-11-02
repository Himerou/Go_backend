package repository

import (
	"GO_web/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type MySQLTodoRepository struct {
	db *sql.DB
}

func NewMySQLTodoRepository(db *sql.DB) *MySQLTodoRepository {
	return &MySQLTodoRepository{
		db: db,
	}
}

func (r *MySQLTodoRepository) Create(ctx context.Context, todo model.Todo) (model.Todo, error) {
	userID, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return model.Todo{}, fmt.Errorf("缺少用户名信息")
	}
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO todos(user_id,title,content,completed)
		VALUES (?,?,?,?)`,
		userID,
		todo.Title,
		todo.Content,
		todo.Completed,
	)
	if err != nil {
		return model.Todo{}, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return model.Todo{}, err
	}
	todo.ID = int(id)
	return todo, nil
}

func (r *MySQLTodoRepository) List(ctx context.Context) ([]model.Todo, error) {
	userID, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return []model.Todo{}, fmt.Errorf("缺少用户名信息")
	}

	var result []model.Todo
	rows, err := r.db.QueryContext(ctx,
		` SELECT id ,title,content, completed,created_at,updated_at
				FROM todos
				where user_id = ?`, userID)
	if err != nil {
		return []model.Todo{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Content,
			&todo.Completed,
			&todo.CreatedAt,
			&todo.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
func (r *MySQLTodoRepository) Get(ctx context.Context, id int) (model.Todo, error) {
	userID, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return model.Todo{}, model.ErrUnauthorized
	}
	var todo model.Todo
	err := r.db.QueryRowContext(ctx, `
				SELECT id, title,content, completed,created_at,updated_at 
				FROM todos 
				WHERE id =? AND user_id = ?`,
		id,
		userID,
	).Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Todo{}, model.ErrTodoNotFound
		}
		return model.Todo{}, err
	}
	return todo, nil
}

func (r *MySQLTodoRepository) Delete(ctx context.Context, id int) error {
	userID, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return model.ErrUnauthorized
	}
	result, err := r.db.ExecContext(ctx, `
			DELETE FROM todos
			WHERE id =? AND user_id = ?`,
		id,
		userID)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrTodoNotFound
	}
	return nil
}
func (r *MySQLTodoRepository) Update(ctx context.Context, id int, todo model.Todo) error {
	userID, ok := ctx.Value(model.UserIDKey).(int)
	if !ok {
		return fmt.Errorf("缺少用户名信息")
	}
	result, err := r.db.ExecContext(ctx,
		`UPDATE todos
				SET title=?,content=?,completed=?
				WHERE id=? AND user_id = ?`,
		todo.Title,
		todo.Content,
		todo.Completed,
		id, userID)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return model.ErrTodoNotFound
	}
	return nil
}
