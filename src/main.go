package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/jftrb/mugacke-backend/internal/handlers"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	
  err := godotenv.Load() 
  if err != nil {
    log.Fatal().Msg("Error loading .env file")
		return
  }

	var r *chi.Mux = chi.NewRouter()
	handlers.MainHandler(r)

	log.Info().
		Msg("Started App")

	err = http.ListenAndServe(":1128", r)
	if err != nil {
		log.Err(err).Msg("Error listening http port")
	}

	log.Info().
		Msg("Exiting App")
}
