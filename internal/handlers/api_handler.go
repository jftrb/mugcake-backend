package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
	"github.com/jftrb/mugacke-backend/internal/middleware"
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
	return router
}
