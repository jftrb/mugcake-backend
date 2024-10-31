package handlers

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/jftrb/mugacke-backend/internal/middleware"
	"github.com/jftrb/mugacke-backend/src/api"
)

func MainHandler(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)
	r.Use(chimiddle.Logger)
	r.Use(middleware.CorsAllowOrigin)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Mount("/api", ApiRouter())
}

func ApiRouter() chi.Router {
	router := chi.NewRouter()

	router.Mount("/recipes", RecipeRouter())
	router.Mount("/users", UserRouter())
	router.Get("/extractor/key", GetRecipeExtractorKey)
	return router
}

func GetRecipeExtractorKey(w http.ResponseWriter, r *http.Request) {
	response := api.GetExtractorKeyResponse{Key: os.Getenv("GEMINI_API_KEY")}
	middleware.EncodeResponse(w, response)
}