// Package repository provides data access objects for interacting with the database.
package repository

import (
	"context"
	"log"
	"recipe-generator/internal/api/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// IngredientsRepository handles database operations related to recipe ingredients.
// It provides methods to create, read, update, and delete ingredient records.
type IngredientsRepository struct {
	db *pgxpool.Pool // Database connection pool
}

// NewIngredientsRepository creates a new instance of IngredientsRepository.
// It requires a database connection pool to perform database operations.
func NewIngredientsRepository(db *pgxpool.Pool) *IngredientsRepository {
	return &IngredientsRepository{db: db}
}

// Insert adds a new ingredient to the database within a transaction.
// It requires a context, the ingredient model, the associated recipe ID, and an active transaction.
// Returns an error if the insertion fails.
func (ir *IngredientsRepository) Insert(ctx context.Context, ingredient *model.Ingredient, recipeId int, tx pgx.Tx) error {
	log.Printf("Inserting ingredient: %s", ingredient.IngredientName)

	//debug
	log.Printf("ingredient unit of measurement: %s", ingredient.UnitOfMeasurement)
	log.Printf("ingredient name: %s", ingredient.IngredientName)
	log.Printf("ingredient amount: %f", ingredient.Amount)
	log.Printf("recipe id: %d", recipeId)

	query := `
		INSERT INTO ingredients (
			unit_of_measurement,
			ingredient_name,
			unit_amount,
			recipe_id,
			created_by,
			created_date,
			updated_by,
			updated_date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8
		)
		`

	_, err := tx.Exec(ctx, query, ingredient.UnitOfMeasurement, ingredient.IngredientName, ingredient.Amount, recipeId, 1, time.Now(), 1, time.Now())
	if err != nil {
		log.Printf("Error inserting ingredient: %v", err)
		return err
	}

	log.Printf("Successfully inserted ingredient: %s", ingredient.IngredientName)
	return nil
}
