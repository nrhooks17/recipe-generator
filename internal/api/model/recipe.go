// Package model provides data structures and error types for the recipe generator application.
package model

import (
	"time"
	"net/http"
	"encoding/json"
	"fmt"
)

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



// decodeRecipe parses the HTTP request body into a Recipe struct.
// It also sets default values for creation and update metadata.
//
// Parameters:
//   - r: The HTTP request containing the recipe data
//
// Returns:
//   - *model.Recipe: The decoded recipe with metadata
//   - error: An error if decoding fails
func (r *Recipe) DecodeRecipe(request *http.Request) (*Recipe, error) {
	var recipe Recipe
	if err := json.NewDecoder(request.Body).Decode(&recipe); err != nil {
		return nil, fmt.Errorf("error decoding request body: %v", err)
	}

	recipe.CreatedBy = 1 // Dummy user ID
	recipe.UpdatedBy = 1 // Same dummy user ID
	now := time.Now()
	recipe.CreatedDate = now
	recipe.UpdatedDate = now

	return &recipe, nil
}
