package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/jftrb/mugacke-backend/internal/dbWrapper"
	"github.com/jftrb/mugacke-backend/internal/encoders"
	"github.com/jftrb/mugacke-backend/internal/middleware"
	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/jftrb/mugacke-backend/src/api/models"
	"github.com/rs/zerolog/log"
)

func RecipeRouter() chi.Router {
	recipeRouter := chi.NewRouter()

	recipeRouter.Route("/summaries", func(r chi.Router) {
		r.Use(middleware.Paginate[api.RecipeSummaryPaginationRequest])
		r.Get("/", GetRecipeSummaries)
		r.Options("/", middleware.CorsPreflight)
	})

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
	db := dbWrapper.NewDbWrapper()
	defer db.Disconnect()

	userId := "18c47dfb-442f-423a-b0cd-70c8076cb7a9"
	pagination := r.Context().Value(middleware.ContextKeyPagination).(api.RecipeSummaryPaginationRequest)
	searchParams := r.Context().Value(middleware.ContextKeySearchParams).(api.RecipeSearchRequest)

	summaries, err := db.GetRecipeSummaries(userId, pagination, searchParams)
	if err != nil {
		log.Err(err).Msg("Error during Get Recipe Summaries operation.")
		api.RequestErrorHandler(w, err)
		return
	}
	
	encodedNextCursor := getNextCursor(len(summaries), pagination)
	response := api.GetRecipeSummariesResponse{
		Summaries: summaries,
		NextCursor: encodedNextCursor,
	}
	
	middleware.EncodeResponse(w, response)
}

func getNextCursor(resultsLength int, pagination api.RecipeSummaryPaginationRequest) string {
	encodedNextCursor := ""
	if resultsLength > pagination.Limit {
		nextOffset := pagination.Offset + pagination.Limit
		nextCursor := fmt.Sprintf("offset:%d", nextOffset)
		encodedNextCursor = encoders.EncodeToBase64(nextCursor)
	}
	return encodedNextCursor
}


func GetRecipe(w http.ResponseWriter, r *http.Request) {
	sRecipeID := chi.URLParam(r, "recipeID")
	recipeID, err := strconv.Atoi(sRecipeID)
	if err != nil {
		log.Err(err).Msg("URL Param 'recipeID' is of an invalid format.")
		api.RequestErrorHandler(w, err)
		return
	}

	db := dbWrapper.NewDbWrapper()
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

	db := dbWrapper.NewDbWrapper()
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

	db := dbWrapper.NewDbWrapper()
	defer db.Disconnect()
	if err := db.PutRecipe(recipeID, recipe); err != nil {
		log.Err(err).Msg("Error - unable to Put Recipe.")
		api.RequestErrorHandler(w, err)
	}

	w.WriteHeader(http.StatusOK)
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

	patchRequest , err := middleware.DecodeQueryParams[api.PatchRecipeRequest](r.URL.Query())
	if err != nil {
		log.Err(err).Msg("Invalid PATCH Recipe query params")
		api.RequestErrorHandler(w, err)
		return
	}

	db := dbWrapper.NewDbWrapper()
	defer db.Disconnect()
	if err := db.PatchRecipe(recipeID, patchRequest.Favorite); err != nil {
		log.Err(err).Msg("Error - unable to Put Recipe.")
		api.RequestErrorHandler(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	sRecipeID := chi.URLParam(r, "recipeID")
	recipeID, err := strconv.Atoi(sRecipeID)
	if err != nil {
		log.Err(err).Msg("URL Param 'recipeID' is of an invalid format.")
		api.RequestErrorHandler(w, err)
		return
	}

	db := dbWrapper.NewDbWrapper()
	defer db.Disconnect()
	if err := db.DeleteRecipe(recipeID); err != nil {
		log.Err(err).Msg("Error - unable to Delete Recipe.")
		api.RequestErrorHandler(w, err)
	}

	w.WriteHeader(http.StatusOK)
}
