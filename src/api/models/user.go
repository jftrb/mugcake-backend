package models

type User struct {
	ID    string `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
}