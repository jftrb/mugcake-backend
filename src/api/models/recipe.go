package models

type PrepInfo struct {
	PrepTime  string
	CookTime  string
	TotalTime string
	Yield     string
}

type Ingredient struct {
	Quantity   float32
	Unit       string
	Ingredient string
	Other      string
}

type Recipe struct {
	ID          int
	Favorite    bool
	Title       string
	URL         string
	ImageSource string
	PrepInfo    PrepInfo
	Tags        []string
	Ingredients []Ingredient
	Directions  []string
	Notes       []string
}
