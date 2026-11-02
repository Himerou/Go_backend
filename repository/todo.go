package repository

import (
	"GO_web/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func OpenMySQL(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开数据库配置失败: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("连接 MySQL 失败: %w", err)
	}

	return db, nil
}

type TodoRepository struct {
	todos  []model.Todo
	nextID int
}

func NewTodoRepository() *TodoRepository {
	return &TodoRepository{
		nextID: 1,
	}
}

func (r *TodoRepository) Create(todo model.Todo) (model.Todo, error) {
	todo.ID = r.nextID
	r.nextID++
	r.todos = append(r.todos, todo)
	return todo, nil
}

func (r *TodoRepository) List() ([]model.Todo, error) {
	return r.todos, nil
}

func (r *TodoRepository) Get(ctx context.Context, id int) (model.Todo, error) {
	for _, todo := range r.todos {
		if todo.ID == id {
			return todo, nil
		}
	}
	var ErrTodoNotfound = errors.New("未找到对应数据")
	return model.Todo{}, ErrTodoNotfound
}

func (r *TodoRepository) Delete(id int) error {
	for i, todo := range r.todos {
		if todo.ID == id {
			r.todos = append(r.todos[:i], r.todos[i+1:]...)
			return nil
		}
	}
	var errTododelete = errors.New("删除失败")
	return errTododelete
}
