package dbWrapper_test

import (
	"testing"

	"github.com/jftrb/mugacke-backend/internal/dbWrapper"
	"github.com/jftrb/mugacke-backend/src/api"
	"github.com/stretchr/testify/assert"
)

func Test_ApplySort(t *testing.T) {
	type args struct {
		sort  api.Sort
	}
	tests := []struct {
		name          string
		args          args
		expectedString string
	}{
		{name: "Favorite", args: args{sort: api.Favorite}, expectedString: "favorite DESC"},
		{name: "Newest", args: args{sort: api.Newest}, expectedString: "created DESC"},
		{name: "Oldest", args: args{sort: -api.Newest}, expectedString: "created"},
		{name: "Alphabetical", args: args{sort: api.Alphabetical}, expectedString: "title"},
		{name: "Alphabetical (reversed)", args: args{sort: -api.Alphabetical}, expectedString: "title DESC"},
		{name: "Most Recent", args: args{sort: api.RecentlyUsed}, expectedString: "last_viewed DESC"},
		{name: "Least Recent", args: args{sort: -api.RecentlyUsed}, expectedString: "last_viewed"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortString := dbWrapper.ApplySort(tt.args.sort)
			assert.Equal(t, tt.expectedString, sortString)
		})
	}
}

func Test_BuildSortQuery(t *testing.T) {
	sortArray := []api.Sort{api.Favorite, -api.Newest}

	expectedSortQuery := "ORDER BY favorite DESC,created"
	sortQuery := dbWrapper.BuildSortQuery(sortArray)
	assert.Equal(t, expectedSortQuery, sortQuery)
}

func Test_BuildSearchQuery(t *testing.T) {
	type args struct {
		query  string
		tags []string
	}
	tests := []struct {
		name        string
		args        args
		expectedSubstring 	string
	}{
		{name: "No Tags", args: args{query: "title", tags: []string{}}, expectedSubstring: "position(LOWER('title') in LOWER(title)) > 0 AND ARRAY[LOWER('title')]"},
		{name: "With Tags", args: args{query: "title", tags: []string{"tag 1", "tag 2"}}, 
			expectedSubstring: `position(LOWER('title') in LOWER(title)) > 0 AND ARRAY[LOWER('tag 1'), LOWER('tag 2'), LOWER('title')]`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			searchParams := api.RecipeSearchRequest{Query: tt.args.query, Tags: tt.args.tags}
			sortString := dbWrapper.BuildSearchQuery(searchParams)

			assert.Contains(t, sortString, tt.expectedSubstring)
		})
	}
}