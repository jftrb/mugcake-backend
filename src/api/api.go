package api

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrBadQuery error = errors.New("invalid query parameter values")

type Error struct {
	Code    int
	Message string
}

func writeError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Message: message,
		Code:    code,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(resp)
}

func RequestErrorHandler(w http.ResponseWriter, err error) {
	writeError(w, err.Error(), http.StatusBadRequest)
}

func InternalErrorHandler(w http.ResponseWriter) {
	writeError(w, "An unexpected error occured", http.StatusInternalServerError)
}
