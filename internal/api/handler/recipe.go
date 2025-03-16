// Package handlers provides HTTP handlers for the recipe API.
// It contains handlers for creating, retrieving, and managing recipes.
package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"recipe-generator/internal/api/config"
	"recipe-generator/internal/api/model"
	"recipe-generator/internal/api/repository"
)

// RecipeHandler manages HTTP requests related to recipes.
// It provides methods for creating, retrieving, and managing recipes
// and their associated ingredients and procedures.
type RecipeHandler struct {
	// ConnectionPool is the PostgreSQL connection pool
	ConnectionPool        *pgxpool.Pool
	// RecipeRepository handles database operations for recipes
	RecipeRepository      *repository.RecipeRepository
	// IngredientsRepository handles database operations for ingredients
	IngredientsRepository *repository.IngredientsRepository
	// ProcedureRepository handles database operations for procedures
	ProcedureRepository   *repository.ProcedureRepository
	// Config contains application configuration
	Config                *config.Config
}

// NewRecipeHandler creates a new RecipeHandler instance with the provided database connection pool and configuration.
// It initializes all required repositories for recipe management.
//
// Parameters:
//   - pool: PostgreSQL connection pool for database operations
//   - config: Application configuration
//
// Returns:
//   - *RecipeHandler: A new recipe handler instance
func NewRecipeHandler(pool *pgxpool.Pool, config *config.Config) *RecipeHandler {
	return &RecipeHandler{
		ConnectionPool:        pool,
		RecipeRepository:      repository.NewRecipeRepository(pool),
		IngredientsRepository: repository.NewIngredientsRepository(pool),
		ProcedureRepository:   repository.NewProcedureRepository(pool),
		Config:                config,
	}
}

// Post returns an HTTP handler function that processes POST requests for creating new recipes.
// The handler validates the recipe data, creates a database transaction, and inserts the recipe,
// its ingredients, and procedure steps into the database.
//
// Returns:
//   - http.HandlerFunc: A handler function that processes recipe creation requests
func (rh *RecipeHandler) Post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Method not allowed",
			})
			return
		}

		log.Printf("Received new recipe submission request")

		recipe, err := rh.decodeRecipe(r)
		if err != nil {
			log.Printf("Error decoding request body: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid request body",
			})
			return
		}

		ingredients := recipe.Ingredients
		procedure := recipe.Procedure

		log.Printf("Successfully decoded recipe: %s with %d ingredients and %d procedure steps",
			recipe.RecipeName,
			len(ingredients),
			len(procedure))

		if err := rh.validateRecipe(w, recipe); err != nil {
			return
		}

		if err := rh.validateIngredients(w, ingredients); err != nil {
			return
		}

		// Start a new transaction
		tx, err := rh.ConnectionPool.Begin(r.Context())
		if err != nil {
			log.Printf("Error starting transaction: %v", err)
			rh.handleRecipeSubmissionError(w, err)
			return
		}
		defer tx.Rollback(r.Context()) // Rollback if we don't commit

		// Insert the recipe into the database using the transaction
		savedRecipe, err := rh.submitRecipe(r.Context(), recipe, tx)
		if err != nil {
			rh.handleRecipeSubmissionError(w, err)
			return
		}

		// Insert ingredients into the database using the transaction
		err = rh.submitIngredients(r.Context(), ingredients, savedRecipe.ID, tx)
		if err != nil {
			rh.handleIngredientSubmissionError(w, err)
			return
		}

		err = rh.submitProcedure(r.Context(), procedure, savedRecipe.ID, tx)
		if err != nil {
			rh.handleProcedureSubmissionError(w, err)
			return
		}

		// Commit the transaction
		if err := tx.Commit(r.Context()); err != nil {
			log.Printf("Error committing transaction: %v", err)
			rh.handleRecipeSubmissionError(w, err)
			return
		}

		log.Printf("Successfully inserted recipe: %s with ID: %d", recipe.RecipeName, recipe.ID)
		json.NewEncoder(w).Encode(recipe)
	}
}

// Get returns an HTTP handler function that processes GET requests for retrieving recipes.
// It fetches a recipe from the database and returns it as JSON.
//
// Returns:
//   - http.HandlerFunc: A handler function that processes recipe retrieval requests
func (rh *RecipeHandler) Get() http.HandlerFunc {
	// grab a recipe from the database and return it

	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]string{"recipe": "recipeObject"})
		if err != nil {
			return
		}
	}
}

// GetRandom returns an HTTP handler function that selects and returns a random recipe from the database.
// The handler responds with a randomly selected recipe in JSON format.
//
// Returns:
//   - http.HandlerFunc: A handler function that processes random recipe retrieval requests
func (rh *RecipeHandler) GetRandom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			
		}
		err := json.NewEncoder(w).Encode(map[string]string{"randomRecipe": "randomRecipe returned"})
		if err != nil {
			return
		}
	}
}

// private functions

// decodeRecipe parses the HTTP request body into a Recipe struct.
// It also sets default values for creation and update metadata.
//
// Parameters:
//   - r: The HTTP request containing the recipe data
//
// Returns:
//   - *model.Recipe: The decoded recipe with metadata
//   - error: An error if decoding fails
func (rh *RecipeHandler) decodeRecipe(r *http.Request) (*model.Recipe, error) {
	var recipe model.Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		return nil, fmt.Errorf("error decoding request body: %v", err)
	}

	recipe.CreatedBy = 1 // Dummy user ID
	recipe.UpdatedBy = 1 // Same dummy user ID
	now := time.Now()
	recipe.CreatedDate = now
	recipe.UpdatedDate = now

	return &recipe, nil
}

// validateRecipe checks if the recipe data is valid.
// If validation fails, it writes an appropriate error response to the HTTP response writer.
//
// Parameters:
//   - w: The HTTP response writer
//   - recipe: The recipe to validate
//
// Returns:
//   - error: An error if validation fails, nil otherwise
func (rh *RecipeHandler) validateRecipe(w http.ResponseWriter, recipe *model.Recipe) error {
	if err := recipe.Validate(); err != nil {
		log.Printf("Recipe validation failed: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": fmt.Sprintf("Recipe validation failed: %v", err),
		})
		return err
	}
	return nil
}

// validateIngredients checks if all ingredients in the recipe are valid.
// If validation fails for any ingredient, it writes an appropriate error response to the HTTP response writer.
//
// Parameters:
//   - w: The HTTP response writer
//   - ingredients: The list of ingredients to validate
//
// Returns:
//   - error: An error if validation fails for any ingredient, nil otherwise
func (rh *RecipeHandler) validateIngredients(w http.ResponseWriter, ingredients []model.Ingredient) error {
	for _, ingredient := range ingredients {
		if err := ingredient.Validate(); err != nil {
			log.Printf("Ingredient validation failed: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{
				"error": fmt.Sprintf("Ingredient validation failed: %v", err),
			})
			return err
		}
	}
	return nil
}

// submitRecipe inserts a recipe into the database using the provided transaction.
//
// Parameters:
//   - ctx: The context for database operations
//   - recipe: The recipe to insert
//   - tx: The database transaction
//
// Returns:
//   - *model.Recipe: The saved recipe with its database ID
//   - error: An error if the insertion fails
func (rh *RecipeHandler) submitRecipe(ctx context.Context, recipe *model.Recipe, tx pgx.Tx) (*model.Recipe, error) {
	// insert the recipe into the database
	savedRecipe, err := rh.RecipeRepository.Insert(ctx, recipe, tx)
	if err != nil {
		log.Printf("Error when submitting a recipe to the database: %v", err)
		return nil, err
	}
	return savedRecipe, nil
}

// handleRecipeSubmissionError writes an appropriate error response when recipe submission fails.
// In development mode, it includes detailed error information.
//
// Parameters:
//   - w: The HTTP response writer
//   - err: The error that occurred during recipe submission
func (rh *RecipeHandler) handleRecipeSubmissionError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	var response map[string]string
	environment := strings.ToLower(rh.Config.Environment)

	if environment == "development" {
		response = map[string]string{
			"error":   "Error when submitting a recipe to the database:",
			"details": err.Error(),
		}
	} else {
		response = map[string]string{
			"error": "Internal server error",
		}
	}

	json.NewEncoder(w).Encode(response)
}

// submitIngredients inserts all ingredients for a recipe into the database using the provided transaction.
//
// Parameters:
//   - ctx: The context for database operations
//   - ingredients: The list of ingredients to insert
//   - recipeID: The ID of the recipe to which the ingredients belong
//   - tx: The database transaction
//
// Returns:
//   - error: An error if any ingredient insertion fails, nil otherwise
func (rh *RecipeHandler) submitIngredients(ctx context.Context, ingredients []model.Ingredient, recipeID int, tx pgx.Tx) error {
	for _, ingredient := range ingredients {
		err := rh.IngredientsRepository.Insert(ctx, &ingredient, recipeID, tx)
		if err != nil {
			log.Printf("Error inserting ingredient: %v", err)
			return err
		}
	}
	return nil
}

// handleIngredientSubmissionError writes an appropriate error response when ingredient submission fails.
// In development mode, it includes detailed error information.
//
// Parameters:
//   - w: The HTTP response writer
//   - err: The error that occurred during ingredient submission
func (rh *RecipeHandler) handleIngredientSubmissionError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	var response map[string]string
	environment := strings.ToLower(rh.Config.Environment)

	if environment == "development" {
		response = map[string]string{
			"error":   "Error inserting ingredient",
			"details": err.Error(),
		}
	} else {
		response = map[string]string{
			"error": "Internal server error",
		}
	}
	json.NewEncoder(w).Encode(response)
}

// submitProcedure inserts all procedure steps for a recipe into the database using the provided transaction.
//
// Parameters:
//   - ctx: The context for database operations
//   - procedure: The list of procedure steps to insert
//   - recipeID: The ID of the recipe to which the procedure steps belong
//   - tx: The database transaction
//
// Returns:
//   - error: An error if any procedure step insertion fails, nil otherwise
func (rh *RecipeHandler) submitProcedure(ctx context.Context, procedure []string, recipeID int, tx pgx.Tx) error {

	for _, step := range procedure {
		err := rh.ProcedureRepository.Insert(ctx, step, recipeID, tx)
		if err != nil {
			log.Printf("Error inserting procedure: %v", err)
			return err
		}
	}
	return nil
}

// handleProcedureSubmissionError writes an appropriate error response when procedure submission fails.
// In development mode, it includes detailed error information.
//
// Parameters:
//   - w: The HTTP response writer
//   - err: The error that occurred during procedure submission
func (rh *RecipeHandler) handleProcedureSubmissionError(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)

	var response map[string]string

	environment := strings.ToLower(rh.Config.Environment)

	if environment == "development" {
		response = map[string]string{
			"error":   "Error inserting procedure",
			"details": err.Error(),
		}
	} else {
		response = map[string]string{
			"error": "Internal server error",
		}
	}

	json.NewEncoder(w).Encode(response)
}
