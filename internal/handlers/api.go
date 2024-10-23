package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
	chimiddle "github.com/go-chi/chi/middleware"
)

func MainHandler(r *chi.Mux) {
	r.Use(chimiddle.StripSlashes)
	r.Use(chimiddle.Logger)

	r.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
}