package service

import (
    "time"
    "fmt"
    "io"
    "encoding/json"


    "recipe-generator/internal/api/model"
)

type RecipeService struct {}


func NewRecipeService() *RecipeService{
    return &RecipeService{}
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
func (rs *RecipeService) DecodeRecipe(recipeBody io.Reader) (*model.Recipe, error) {
	var recipe model.Recipe
    
	if err := json.NewDecoder(recipeBody).Decode(&recipe); err != nil {
		return nil, fmt.Errorf("error decoding request body: %v", err)
	}

	recipe.CreatedBy = 1 // Dummy user ID
	recipe.UpdatedBy = 1 // Same dummy user ID
	now := time.Now()
	recipe.CreatedDate = now
	recipe.UpdatedDate = now

	return &recipe, nil
}
