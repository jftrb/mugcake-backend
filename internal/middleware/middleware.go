package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/rs/zerolog/log"
)

type ContextKey string

const (
	ContextKeyRecipeId ContextKey = "pageToken"
)

func EncodeResponse(w http.ResponseWriter, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Err(err).Msg("Error during response JSON encoding.")
		api.InternalErrorHandler(w)
	}
}