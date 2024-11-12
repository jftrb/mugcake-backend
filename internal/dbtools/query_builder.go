package dbtools

import (
	"fmt"
	"math"
	"slices"
	"strings"

	"github.com/jftrb/mugacke-backend/src/api"
)

// Search behaviour :
// Only Query present, search in Title or in Tags (TODO: or in Ingredients)
// Only Tags present, search in tags and all tags must be present
// Query + Tags, search a title that has all tags present
func BuildSearchQuery(searchParams api.RecipeSearchRequest) string {
	// Query is substring of Title
	out := fmt.Sprintf(`position(LOWER('%s') in LOWER(title)) > 0 `, searchParams.Query)

	tagsToSearch := slices.DeleteFunc(searchParams.Tags, func(s string) bool {return s == ""})
	joinMethod := "AND "
	if len(tagsToSearch) == 0 {
		joinMethod = "OR "
		tagsToSearch = append(searchParams.Tags, searchParams.Query)
	}

	// AND ARRAY[LOWER('Tag 1'), LOWER('Tag 2')] <@ tags
	out += joinMethod
	out += fmt.Sprintf("ARRAY[LOWER('%s')] ", strings.Join(tagsToSearch, "'), LOWER('"))
	out += `<@ ARRAY (
					SELECT LOWER(tags.name)
					FROM   unnest(r.tags) AS a(tag_id)
					JOIN   tags ON tags.id = a.tag_id
					)`

	return out
}

func BuildSortQuery(s []api.Sort) string {
	out := ""
	for _, val := range s {
		out += ApplySort(val) + ","
	}
	out = strings.Trim(out, ",")

	if len(out) == 0 { return out } 
	return "ORDER BY " + out
}

func ApplySort(s api.Sort) string {
	absoluteSort := math.Abs(float64(s))
	negative := math.Signbit(float64(s))
	switch sort := api.Sort(absoluteSort); sort {
		case api.Favorite:
			if negative { return "favorite" }
			return "favorite DESC"
		case api.Newest:
			if negative { return "created" }
			return "created DESC"
		case api.RecentlyUsed:
			if negative { return "last_viewed" }
			return "last_viewed DESC"
		case api.Alphabetical:
			if negative { return "title DESC" }
			return "title"
	}
	return ""
}