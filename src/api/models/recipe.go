package models

type PrepInfo struct {
	PrepTime  string `json:"prepTime,omitempty"`
	CookTime  string `json:"cookTime,omitempty"`
	TotalTime string `json:"totalTime,omitempty"`
	Yield     string `json:"yield,omitempty"`
}

type Ingredient struct {
	Quantity   float32 `json:"quantity,omitempty"`
	Unit       string  `json:"unit,omitempty"`
	Ingredient string  `json:"ingredient,omitempty"`
	Other      string  `json:"other,omitempty"`
}

type IngredientSection struct {
	Header      string       `json:"header,omitempty"`
	Ingredients []Ingredient `json:"ingredients,omitempty"`
}

type Recipe struct {
	Favorite           bool                `json:"favorite,omitempty"`
	Title              string              `json:"title,omitempty"`
	URL                string              `json:"url,omitempty"`
	ImageSource        string              `json:"imageSource,omitempty"`
	PrepInfo           PrepInfo            `json:"prepInfo,omitempty"`
	Tags               []string            `json:"tags,omitempty"`
	IngredientSections []IngredientSection `json:"ingredientSections,omitempty"`
	Directions         []string            `json:"directions,omitempty"`
	Notes              []string            `json:"notes,omitempty"`
}
