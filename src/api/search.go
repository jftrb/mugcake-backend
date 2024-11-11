package api

type Sort int

const (
	Favorite Sort = iota
	Newest
	RecentlyUsed
	Alphabetical
)

type RecipeSearchRequest struct {
	Query string
	Tags []string
	SortBy []Sort
}

