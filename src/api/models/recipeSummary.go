package models

type RecipeSummary struct {
	ID          int
	Favorite    bool
	Title       string
	ImageSource string
	TotalTime   string
	Tags        []string
}