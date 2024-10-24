package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	dbwrapper "github.com/jftrb/mugacke-backend/internal/dbWrapper"
	"github.com/jftrb/mugacke-backend/internal/middleware"
	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/rs/zerolog/log"
)

func RecipeRouter() chi.Router {
	recipeRouter := chi.NewRouter()

	recipeRouter.Get("/summaries", GetRecipeSummaries)
	recipeRouter.Get("/{recipeID:^[0-9]$}", GetRecipe)
	return recipeRouter
}

func GetRecipeSummaries(w http.ResponseWriter, r *http.Request) {
	db := dbwrapper.NewDbWrapper()
	defer db.Disconnect()

	userId := "18c47dfb-442f-423a-b0cd-70c8076cb7a9"
	summaries, err := db.GetRecipeSummaries(userId)
	if err != nil {
		log.Err(err).Msg("Error during Get Recipe Summaries operation.")
		api.RequestErrorHandler(w, err)
		return
	}
	
	response := api.GetRecipeSummariesResponse{
		Summaries: summaries,
	}
	
	middleware.EncodeResponse(w, response)
}


func GetRecipe(w http.ResponseWriter, r *http.Request) {
	db := dbwrapper.NewDbWrapper()
	defer db.Disconnect()

	sRecipeID := chi.URLParam(r, "recipeID")
	recipeID, err := strconv.Atoi(sRecipeID)
	if err != nil {
		log.Err(err).Msg("URL Param 'recipeID' is of an invalid format.")
		api.RequestErrorHandler(w, err)
		return
	}

	recipe, err := db.GetRecipe(recipeID)
	if err != nil {
		log.Err(err).Msg("Error during Get Recipe operation.")
		api.RequestErrorHandler(w, err)
		return
	}
	
	response := api.GetRecipeResponse{
		Recipe: recipe,
	}
	
	middleware.EncodeResponse(w, response)
}
