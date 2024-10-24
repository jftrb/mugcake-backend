package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	dbwrapper "github.com/jftrb/mugacke-backend/internal/dbWrapper"
	"github.com/jftrb/mugacke-backend/internal/middleware"
	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/rs/zerolog/log"
)

func UserRouter() chi.Router {
	userRouter := chi.NewRouter()
	userRouter.Get("/", GetUsers)

	// userRouter.Get("/{userID}", GetUser)
	return userRouter
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	db := dbwrapper.NewDbWrapper()
	defer db.Disconnect()

	users, err := db.GetUsers()
	if err != nil {
		log.Err(err).Msg("Error during Fetch Users operation.")
		api.RequestErrorHandler(w, err)
		return
	}

	response := api.GetUsersResponse{
		Users: users,
	}

	middleware.EncodeResponse(w, response)
}