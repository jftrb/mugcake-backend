package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jftrb/mugacke-backend/src/api/models"
)

var ErrBadQuery error = errors.New("invalid query parameter values")

type Error struct {
	Code    int
	Message string
}

type GetExtractorKeyResponse struct {
	Key string
}

type GetUsersResponse struct {
	Users 	[]models.User
}

type GetRecipeSummariesResponse struct {
	Summaries 	[]models.RecipeSummary
}

type GetRecipeResponse struct {
	Recipe 	models.Recipe
}

type PostRecipeResponse struct {
	ID int
}

type PatchRecipeRequest struct {
	Favorite bool
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
