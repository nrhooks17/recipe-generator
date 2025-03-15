package models

import "time"

type Recipe struct {
	ID              int          `json:"id"`
	RecipeName      string       `json:"recipeName"`
	Description     string       `json:"description,omitempty"`
	PrepTimeMinutes int          `json:"prepTimeMinutes,omitempty"`
	CookTimeMinutes int          `json:"cookTimeMinutes,omitempty"`
	Ingredients     []Ingredient `json:"ingredients"`
	Procedure       []string     `json:"procedure"`
	Servings        int          `json:"servings,omitempty"`
	CreatedBy       int          `json:"createdBy"`
	CreatedDate     time.Time    `json:"createdDate"`
	UpdatedBy       int          `json:"updatedBy"`
	UpdatedDate     time.Time    `json:"updatedDate"`
}

// NewRecipe creates a new Recipe instance with required fields
func NewRecipe(name string, createdBy int) *Recipe {
	now := time.Now()
	return &Recipe{
		RecipeName:  name,
		CreatedBy:   createdBy,
		CreatedDate: now,
		UpdatedBy:   createdBy,
		UpdatedDate: now,
	}
}

// Validate ensures all required fields are present
func (r *Recipe) Validate() error {
	if r.RecipeName == "" {
		return ErrMissingRequiredField("recipeName")
	}
	if r.CreatedBy == 0 {
		return ErrMissingRequiredField("createdBy")
	}
	if r.UpdatedBy == 0 {
		return ErrMissingRequiredField("updatedBy")
	}
	return nil
}
