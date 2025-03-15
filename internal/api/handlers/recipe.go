package handlers

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
	"recipe-generator/internal/api/models"
	"recipe-generator/internal/api/repository"
)

type RecipeHandler struct {
	ConnectionPool        *pgxpool.Pool
	RecipeRepository      *repository.RecipeRepository
	IngredientsRepository *repository.IngredientsRepository
	ProcedureRepository   *repository.ProcedureRepository
	Config                *config.Config
}

func NewRecipeHandler(pool *pgxpool.Pool, config *config.Config) *RecipeHandler {
	return &RecipeHandler{
		ConnectionPool:        pool,
		RecipeRepository:      repository.NewRecipeRepository(pool),
		IngredientsRepository: repository.NewIngredientsRepository(pool),
		ProcedureRepository:   repository.NewProcedureRepository(pool),
		Config:                config,
	}
}

// Post() handles the submission of a new recipe to the database.
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

// handler functions for requests
func (rh *RecipeHandler) Get() http.HandlerFunc {
	// grab a recipe from the database and return it

	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]string{"recipe": "recipeObject"})
		if err != nil {
			return
		}
	}
}

func (rh *RecipeHandler) GetRandom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(map[string]string{"randomRecipe": "randomRecipe returned"})
		if err != nil {
			return
		}
	}
}

// private functions
func (rh *RecipeHandler) decodeRecipe(r *http.Request) (*models.Recipe, error) {
	var recipe models.Recipe
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

func (rh *RecipeHandler) validateRecipe(w http.ResponseWriter, recipe *models.Recipe) error {
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

func (rh *RecipeHandler) validateIngredients(w http.ResponseWriter, ingredients []models.Ingredient) error {
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

// Submits a recipe to the database.
func (rh *RecipeHandler) submitRecipe(ctx context.Context, recipe *models.Recipe, tx pgx.Tx) (*models.Recipe, error) {
	// insert the recipe into the database
	savedRecipe, err := rh.RecipeRepository.Insert(ctx, recipe, tx)
	if err != nil {
		log.Printf("Error when submitting a recipe to the database: %v", err)
		return nil, err
	}
	return savedRecipe, nil
}

// Handles errors when submitting a recipe to the database.
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

func (rh *RecipeHandler) submitIngredients(ctx context.Context, ingredients []models.Ingredient, recipeID int, tx pgx.Tx) error {
	for _, ingredient := range ingredients {
		err := rh.IngredientsRepository.Insert(ctx, &ingredient, recipeID, tx)
		if err != nil {
			log.Printf("Error inserting ingredient: %v", err)
			return err
		}
	}
	return nil
}

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
