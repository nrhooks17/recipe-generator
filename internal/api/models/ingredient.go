package models

import (
	"time"
)

// Ingredient represents an ingredient in a recipe in the database.
type Ingredient struct {
	ID                int       `json:"id"`
	Amount            float64   `json:"amount"`
	UnitOfMeasurement string    `json:"unitOfMeasurement"`
	IngredientName    string    `json:"ingredientName"`
	RecipeId          int       `json:"recipeId"`
	CreatedBy         int       `json:"createdBy"`
	CreatedDate       time.Time `json:"createdDate"`
	UpdatedBy         int       `json:"updatedBy"`
	UpdatedDate       time.Time `json:"updatedDate"`
}

// NewIngredient creates a new Ingredient instance with required fields
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

// Validate() validates the Ingredient instance
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
