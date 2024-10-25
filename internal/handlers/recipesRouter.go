package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	dbwrapper "github.com/jftrb/mugacke-backend/internal/dbWrapper"
	"github.com/jftrb/mugacke-backend/internal/middleware"
	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/jftrb/mugacke-backend/src/api/models"
	"github.com/rs/zerolog/log"
)

func RecipeRouter() chi.Router {
	recipeRouter := chi.NewRouter()

	recipeRouter.Get("/summaries", GetRecipeSummaries)
	recipeRouter.Get("/{recipeID:^[0-9]$}", GetRecipe)
	recipeRouter.Put("/{recipeID:^[0-9]$}", PutRecipe)
	recipeRouter.Delete("/{recipeID:^[0-9]$}", DeleteRecipe)
	recipeRouter.Post("/", PostRecipe)
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
	sRecipeID := chi.URLParam(r, "recipeID")
	recipeID, err := strconv.Atoi(sRecipeID)
	if err != nil {
		log.Err(err).Msg("URL Param 'recipeID' is of an invalid format.")
		api.RequestErrorHandler(w, err)
		return
	}

	db := dbwrapper.NewDbWrapper()
	defer db.Disconnect()

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

func PostRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		log.Err(err).Msg("Error - unable to parse request body into a Recipe.")
		api.RequestErrorHandler(w, err)
		return
	}

	log.Debug().Str("Recipe Title", recipe.Title).Msg("Posting Recipe")

	db := dbwrapper.NewDbWrapper()
	defer db.Disconnect()

	userId := "18c47dfb-442f-423a-b0cd-70c8076cb7a9"
	if err := db.AddRecipe(userId, recipe); err != nil {
		log.Err(err).Msg("Error - unable to Post Recipe.")
		api.RequestErrorHandler(w, err)
	}
}

func PutRecipe(w http.ResponseWriter, r *http.Request) {
	sRecipeID := chi.URLParam(r, "recipeID")
	recipeID, err := strconv.Atoi(sRecipeID)
	if err != nil {
		log.Err(err).Msg("URL Param 'recipeID' is of an invalid format.")
		api.RequestErrorHandler(w, err)
		return
	}

	var recipe models.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		log.Err(err).Msg("Error - unable to parse request body into a Recipe.")
		api.RequestErrorHandler(w, err)
		return
	}

	db := dbwrapper.NewDbWrapper()
	defer db.Disconnect()
	if err := db.UpdateRecipe(recipeID, recipe); err != nil {
		log.Err(err).Msg("Error - unable to Put Recipe.")
		api.RequestErrorHandler(w, err)
	}
}

func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	sRecipeID := chi.URLParam(r, "recipeID")
	recipeID, err := strconv.Atoi(sRecipeID)
	if err != nil {
		log.Err(err).Msg("URL Param 'recipeID' is of an invalid format.")
		api.RequestErrorHandler(w, err)
		return
	}

	db := dbwrapper.NewDbWrapper()
	defer db.Disconnect()
	if err := db.DeleteRecipe(recipeID); err != nil {
		log.Err(err).Msg("Error - unable to Delete Recipe.")
		api.RequestErrorHandler(w, err)
	}
}
