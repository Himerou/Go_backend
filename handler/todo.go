package handler

import (
	"GO_web/model"
	"GO_web/service"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

func ListTodos(service service.ListTodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		content, err := service.List(req.Context())
		if err != nil {
			WriteError(w, "内部错误", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(content)
	}

}
func CreateTodo(service service.CreateTodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var todo model.Todo
		if err := json.NewDecoder(req.Body).Decode(&todo); err != nil {
			WriteError(w, "格式错误", http.StatusBadRequest)
			return
		}
		created, err := service.Create(req.Context(), todo)
		if err != nil {
			if errors.Is(err, model.ErrEmptyContent) {
				WriteError(w, "内容不得为空", http.StatusBadRequest)
				return
			}
			WriteError(w, "内部错误", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(created)
	}

}
func GetTodo(service service.GetTodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			WriteError(w, "无效的ID", http.StatusBadRequest)
			return
		}
		todo, err := service.Get(req.Context(), id)
		if err != nil {
			if errors.Is(err, model.ErrTodoNotFound) {
				WriteError(w, "未查询到信息", http.StatusNotFound)
				return
			}
			WriteError(w, "内部错误", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(todo)
	}

}
func DeleteTodo(service service.DeleteTodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			WriteError(w, "无效ID", http.StatusBadRequest)
			return
		}
		err = service.Delete(req.Context(), id)
		if err != nil {
			if errors.Is(err, model.ErrTodoNotFound) {
				WriteError(w, "内容不存在", http.StatusNotFound)
				return
			}
			if errors.Is(err, model.ErrUnauthorized) {
				WriteError(w, "未登录", http.StatusUnauthorized)
				return
			}
			WriteError(w, "内部错误", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}
func UpdateTodo(service service.UpdateTodoService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		id, err := strconv.Atoi(req.PathValue("id"))
		if err != nil {
			WriteError(w, "无效ID", http.StatusBadRequest)
			return
		}
		var todo model.Todo
		if err = json.NewDecoder(req.Body).Decode(&todo); err != nil {
			WriteError(w, "格式错误", http.StatusBadRequest)
			return
		}
		err = service.Update(req.Context(), id, todo)
		if err != nil {
			if errors.Is(err, model.ErrTodoNotFound) {
				WriteError(w, "内容不存在", http.StatusNotFound)
				return
			}
			WriteError(w, "内部错误", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)

	}
}
