// Package model provides data structures and error types for the recipe generator application.
package model

import "time"

// Recipe represents a complete recipe in the database.
// It contains all the necessary information about a recipe including
// its name, description, preparation details, ingredients, and cooking procedure.
type Recipe struct {
	ID              int          `json:"id"`                        // Unique identifier for the recipe
	RecipeName      string       `json:"recipeName"`                // Name of the recipe
	Description     string       `json:"description,omitempty"`     // Description of the recipe
	PrepTimeMinutes int          `json:"prepTimeMinutes,omitempty"` // Time required for preparation in minutes
	CookTimeMinutes int          `json:"cookTimeMinutes,omitempty"` // Time required for cooking in minutes
	Ingredients     []Ingredient `json:"ingredients"`               // List of ingredients required for the recipe
	Procedure       []string     `json:"procedure"`                 // Step-by-step cooking instructions
	Servings        int          `json:"servings,omitempty"`        // Number of servings the recipe yields
	CreatedBy       int          `json:"createdBy"`                 // User ID who created this recipe
	CreatedDate     time.Time    `json:"createdDate"`               // Timestamp when the recipe was created
	UpdatedBy       int          `json:"updatedBy"`                 // User ID who last updated this recipe
	UpdatedDate     time.Time    `json:"updatedDate"`               // Timestamp when the recipe was last updated
}

// NewRecipe creates a new Recipe instance with required fields.
// It automatically sets the creation and update timestamps to the current time.
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

// Validate checks if the Recipe instance has all required fields properly set.
// It returns an error if any required field is missing or invalid.
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
