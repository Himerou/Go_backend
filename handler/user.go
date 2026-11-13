package handler

import (
	"GO_web/model"
	"GO_web/service"
	"encoding/json"
	"errors"
	"net/http"
)

func GetProfile(service service.GetProfileService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		profile, err := service.GetProfile(ctx)
		if err != nil {
			if errors.Is(err, model.ErrUnauthorized) {
				WriteError(w, "请求错误", http.StatusUnauthorized)
				return
			}
			WriteError(w, "内部错误", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(profile)
	}
}
func GetPublicProfile(service service.GetPublicProfileService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		username := req.PathValue("username")
		ctx := req.Context()
		profile, err := service.GetPublicProfile(ctx, username)
		if err != nil {
			if errors.Is(err, model.ErrUserNotFound) {
				WriteError(w, err.Error(), http.StatusNotFound)
				return
			}
			WriteError(w, "内部错误", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(profile)
	}
}

func UpdateProfile(service service.UpdateProfileService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var msg model.Msg
		if err := json.NewDecoder(req.Body).Decode(&msg); err != nil {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		profile, err := service.UpdateProfile(ctx, msg.Msg)
		if err != nil {
			if errors.Is(err, model.ErrUserNotFound) {
				WriteError(w, err.Error(), http.StatusNotFound)
				return
			}
			if errors.Is(err, model.ErrEmptyContent) {
				WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}
			WriteError(w, "内部错误", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(profile)
	}
}

func UpdatePassword(service service.UpdatePasswordService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var msg model.Msg
		if err := json.NewDecoder(req.Body).Decode(&msg); err != nil {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		err := service.UpdatePassword(ctx, msg.Msg)
		if err != nil {
			if errors.Is(err, model.ErrUserNotFound) {
				WriteError(w, err.Error(), http.StatusNotFound)
				return
			}
			if errors.Is(err, model.ErrDuplicatedPassword) {
				WriteError(w, err.Error(), http.StatusConflict)
				return
			}
			if errors.Is(err, model.ErrEmptyContent) {
				WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}
			WriteError(w, "内部错误", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateUsername(service service.UpdateUsernameService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		var msg model.Msg
		if err := json.NewDecoder(req.Body).Decode(&msg); err != nil {
			WriteError(w, err.Error(), http.StatusBadRequest)
			return
		}
		profile, err := service.UpdateUsername(ctx, msg.Msg)
		if err != nil {
			if errors.Is(err, model.ErrDuplicateUsername) {
				WriteError(w, err.Error(), http.StatusConflict)
				return
			}
			if errors.Is(err, model.ErrEmptyContent) {
				WriteError(w, err.Error(), http.StatusBadRequest)
				return
			}
			WriteError(w, "内部错误", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(profile)
	}
}

func DeleteUser(service service.DeleteUserService) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		err := service.DeleteUser(ctx)
		if err != nil {
			WriteError(w, "内部错误", http.StatusInternalServerError)
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
