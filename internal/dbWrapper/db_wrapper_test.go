package dbWrapper_test

import (
	"os"
	"testing"

	"github.com/jftrb/mugacke-backend/internal/dbWrapper"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func Test_dbWrapper_GetRecipe(t *testing.T) {
  godotenv.Load() 
	// os.Setenv("POSTGRES_URL", "FILL OUT HERE")
	log.Info().Msg(os.Getenv("POSTGRES_URL"))

	db := dbWrapper.NewDbWrapper()
	defer db.Disconnect()

	recipe, err := db.GetRecipe(2)
	log.Err(err).Send()

	assert.NotEmpty(t, recipe.Title)
}
