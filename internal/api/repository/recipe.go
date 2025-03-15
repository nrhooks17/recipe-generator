package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"recipe-generator/internal/api/models"
)

type RecipeRepository struct {
	ConnectionPool *pgxpool.Pool
}

// Constructor for the RecipeRepository
func NewRecipeRepository(pool *pgxpool.Pool) *RecipeRepository {
	return &RecipeRepository{ConnectionPool: pool}
}

// Insert a new recipe into the database
func (r *RecipeRepository) Insert(ctx context.Context, model *models.Recipe, tx pgx.Tx) (*models.Recipe, error) {
	log.Printf("Starting database insertion for recipe: %s", model.RecipeName)

	query := `
		INSERT INTO recipes (
			recipe_name, description, prep_time_minutes, 
			cook_time_minutes, servings, created_by,
			created_date, updated_by, updated_date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id`

	err := tx.QueryRow(
		ctx,
		query,
		model.RecipeName,
		model.Description,
		model.PrepTimeMinutes,
		model.CookTimeMinutes,
		model.Servings,
		model.CreatedBy,
		model.CreatedDate,
		model.UpdatedBy,
		model.UpdatedDate,
	).Scan(&model.ID)

	if err != nil {
		log.Printf("Error inserting recipe into database: %v", err)
		return nil, err
	}

	log.Printf("Successfully inserted recipe with ID: %d", model.ID)
	return model, nil
}
