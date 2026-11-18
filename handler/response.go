package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteError(w http.ResponseWriter, msg string, code int) {
	log.Printf("状态 %v 错误 %v", code, msg)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": msg,
	})

}
