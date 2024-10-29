package models

type RecipeSummary struct {
	RecipeID    int      `json:"recipeId,omitempty"`
	Favorite    bool     `json:"favorite,omitempty"`
	Title       string   `json:"title,omitempty"`
	TotalTime   string   `json:"totalTime,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	ImageSource string   `json:"imageSource,omitempty"`
}