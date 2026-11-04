package handler

import (
	"GO_web/model"
	"GO_web/service"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func Middleware(s *service.SessionStore) func(http.Handler) http.HandlerFunc {
	return func(next http.Handler) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			// 统一处理未登录的请求
			unauthorized := func() {
				if strings.HasPrefix(req.URL.Path, "/api/") {
					WriteError(w, "当前未登录，请登录或注册使用服务", http.StatusUnauthorized)
					return
				}

				http.Redirect(
					w,
					req,
					"/?error=unauthorized",
					http.StatusSeeOther,
				)
			}

			sessionCookie, err := req.Cookie("session_id")
			if err != nil {
				unauthorized()
				return
			}

			sessionID := sessionCookie.Value
			userID, ok := s.Get(sessionID)
			if !ok {
				unauthorized()
				return
			}

			ctx := context.WithValue(
				req.Context(),
				model.UserIDKey,
				userID,
			)
			next.ServeHTTP(w, req.WithContext(ctx))
		}
	}
}

func LoginTodo(service service.LoginService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var reqdata model.LoginAndRegisterRequest
		if err := json.NewDecoder(req.Body).Decode(&reqdata); err != nil {
			WriteError(w, "格式错误", http.StatusBadRequest)
			return
		}
		user, session, err := service.Login(ctx, reqdata.Username, reqdata.Password)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrEmptyUsername):
				WriteError(w, err.Error(), http.StatusBadRequest)
				return

			case errors.Is(err, model.ErrEmptyPassword):
				WriteError(w, err.Error(), http.StatusBadRequest)
				return
			case errors.Is(err, model.ErrInvalidCredentials):
				WriteError(w, err.Error(), http.StatusUnauthorized)
				return
			default:
				WriteError(w, "", http.StatusInternalServerError)
				return
			}
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session,
			Path:     "/",
			MaxAge:   600,
			HttpOnly: true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(user)
	}
}

func RegisterTodo(service service.RegisterService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var reqdata model.LoginAndRegisterRequest
		if err := json.NewDecoder(req.Body).Decode(&reqdata); err != nil {
			WriteError(w, "请求错误", http.StatusBadRequest)
			return
		}
		user, session, err := service.Register(ctx, reqdata.Username, reqdata.Password)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrEmptyUsername):
				WriteError(w, err.Error(), http.StatusBadRequest)
				return
			case errors.Is(err, model.ErrEmptyPassword):
				WriteError(w, err.Error(), http.StatusBadRequest)
				return
			case errors.Is(err, model.ErrDuplicateUsername):
				WriteError(w, err.Error(), http.StatusBadRequest)
				return
			default:
				WriteError(w, "内部错误", http.StatusInternalServerError)
				return
			}
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    session,
			Path:     "/",
			MaxAge:   600,
			HttpOnly: true,
		})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(user)
	}
}
func LogoutTodo(service service.LogoutService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		session, err := req.Cookie("session_id")
		if err != nil {
			WriteError(w, err.Error(), http.StatusUnauthorized)
			return
		}
		err = service.Logout(req.Context(), session.Value)
		if err != nil {
			WriteError(w, "退出登录失败", http.StatusInternalServerError)
			return
		}
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1,
			HttpOnly: true,
		})
		w.WriteHeader(http.StatusNoContent)
	}
}
