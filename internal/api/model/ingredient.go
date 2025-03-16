// Package model provides data structures and error types for the recipe generator application.
package model

import (
	"time"
)

// Ingredient represents an ingredient in a recipe in the database.
// It contains all the necessary information about an ingredient including
// its amount, unit of measurement, and relationship to a recipe.
type Ingredient struct {
	ID                int       `json:"id"`                  // Unique identifier for the ingredient
	Amount            float64   `json:"amount"`              // Quantity of the ingredient
	UnitOfMeasurement string    `json:"unitOfMeasurement"`   // Unit of measurement (e.g., cup, tablespoon)
	IngredientName    string    `json:"ingredientName"`      // Name of the ingredient
	RecipeId          int       `json:"recipeId"`            // Foreign key to the recipe this ingredient belongs to
	CreatedBy         int       `json:"createdBy"`           // User ID who created this ingredient
	CreatedDate       time.Time `json:"createdDate"`         // Timestamp when the ingredient was created
	UpdatedBy         int       `json:"updatedBy"`           // User ID who last updated this ingredient
	UpdatedDate       time.Time `json:"updatedDate"`         // Timestamp when the ingredient was last updated
}

// NewIngredient creates a new Ingredient instance with required fields.
// It automatically sets the creation and update timestamps to the current time.
func NewIngredient(amount float64, unitOfMeasurement string, ingredientName string, recipeId int, createdBy int) *Ingredient {
	now := time.Now()
	return &Ingredient{
		Amount:            amount,
		UnitOfMeasurement: unitOfMeasurement,
		IngredientName:    ingredientName,
		RecipeId:          recipeId,
		CreatedBy:         createdBy,
		CreatedDate:       now,
		UpdatedBy:         createdBy,
		UpdatedDate:       now,
	}
}

// Validate checks if the Ingredient instance has all required fields properly set.
// It returns an error if any required field is missing or invalid.
func (i *Ingredient) Validate() error {
	if i.Amount <= 0 {
		return ErrMissingRequiredField("amount")
	}

	if i.UnitOfMeasurement == "" {
		return ErrMissingRequiredField("unitOfMeasurement")
	}

	if i.IngredientName == "" {
		return ErrMissingRequiredField("ingredientName")
	}

	return nil
}
