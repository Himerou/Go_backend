package handler

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, msg string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})

}
