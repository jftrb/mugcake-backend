package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jftrb/mugacke-backend/internal/dbWrapper"
	"github.com/jftrb/mugacke-backend/internal/middleware"
	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/jftrb/mugacke-backend/src/api/models"
	"github.com/rs/zerolog/log"
)

func RecipeRouter() chi.Router {
	recipeRouter := chi.NewRouter()

	recipeRouter.Get("/summaries", GetRecipeSummaries)
	recipeRouter.Options("/summaries", middleware.CorsPreflight)

	recipeIdRoute := "/{recipeID:^[0-9]+$}"
	recipeRouter.Get(recipeIdRoute, GetRecipe)
	recipeRouter.Put(recipeIdRoute, PutRecipe)
	recipeRouter.Patch(recipeIdRoute, PatchRecipe)
	recipeRouter.Delete(recipeIdRoute, DeleteRecipe)
	recipeRouter.Options(recipeIdRoute, middleware.CorsPreflight)

	recipeRouter.Post("/", PostRecipe)
	recipeRouter.Options("/", middleware.CorsPreflight)
	return recipeRouter
}

// TODO : Paginate summaries
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

// TODO : create new tags if not exist
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
	recipeId, err := db.AddRecipe(userId, recipe); 
	if err != nil {
		log.Err(err).Msg("Error - unable to Post Recipe.")
		api.RequestErrorHandler(w, err)
	}

	response := api.PostRecipeResponse{
		ID: recipeId,
	}
	
	middleware.EncodeResponse(w, response)
}

// TODO : create new tags if not exist
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
	if err := db.PutRecipe(recipeID, recipe); err != nil {
		log.Err(err).Msg("Error - unable to Put Recipe.")
		api.RequestErrorHandler(w, err)
	}
}

// TODO : create new tags if not exist
func PatchRecipe(w http.ResponseWriter, r *http.Request) {
	sRecipeID := chi.URLParam(r, "recipeID")

	recipeID, err := strconv.Atoi(sRecipeID)
	if err != nil {
		log.Err(err).Msg("URL Param 'recipeID' is of an invalid format.")
		api.RequestErrorHandler(w, err)
		return
	}

	patchRequest , err := middleware.DecodeQueryParams[api.PatchRecipeRequest](w, r)
	if err != nil {
		log.Err(err).Msg("Invalid PATCH Recipe query params")
		api.RequestErrorHandler(w, err)
		return
	}

	db := dbwrapper.NewDbWrapper()
	defer db.Disconnect()
	if err := db.PatchRecipe(recipeID, patchRequest.Favorite); err != nil {
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
