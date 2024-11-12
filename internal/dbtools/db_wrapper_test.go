package dbtools_test

import (
	"os"
	"testing"

	"github.com/jftrb/mugacke-backend/internal/dbtools"
	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	log.Info().Msg("Setting up test")
  godotenv.Load("../../.env") 
	m.Run()
}

func Test_dbWrapper_GetRecipe(t *testing.T) {
	assert.NotEmpty(t, os.Getenv("POSTGRES_URL"))

	db := dbtools.NewDbWrapper()
	defer db.Disconnect()

	recipe, err := db.GetRecipe(2)
	log.Err(err).Send()

	assert.NotEmpty(t, recipe.Title)
}

func Test_dbWrapper_GetRecipeSummaries(t *testing.T) {
	assert.NotEmpty(t, os.Getenv("POSTGRES_URL"))

	db := dbtools.NewDbWrapper()
	defer db.Disconnect()

	
	type args struct {
		query string
		tags 	[]string
	}
	tests := []struct {
		name          string
		args          args
		expectedResults int
	}{
		{name: "Only query", args: args{query: "vietnamese tomato tofu", tags: []string{}}, expectedResults: 1},
		{name: "Only tag", args: args{query: "", tags: []string{"dessert", "tag 2"}}, expectedResults: 1},
		{name: "Tag and query", args: args{query: "vietnamese tomato tofu", tags: []string{"tag 2"}}, expectedResults: 1},
		{name: "Incompatible tag and query", args: args{query: "vietnamese tomato tofu", tags: []string{"dessert"}}, expectedResults: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userId := "18c47dfb-442f-423a-b0cd-70c8076cb7a9"
			ctx := dbtools.GetSummariesContext{
				Limit: 10,
				Offset: 0,
				SearchParams: api.RecipeSearchRequest{Query: tt.args.query, Tags: tt.args.tags, SortBy: []api.Sort{}},
			}
			summaries, err := db.GetRecipeSummaries(userId, ctx)
			log.Err(err).Send()
		
			assert.Equal(t, tt.expectedResults, len(summaries))
		})
	}

}
