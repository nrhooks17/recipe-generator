package repository

import (
	"context"
	"log"
	"recipe-generator/internal/api/models"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type IngredientsRepository struct {
	db *pgxpool.Pool
}

func NewIngredientsRepository(db *pgxpool.Pool) *IngredientsRepository {
	return &IngredientsRepository{db: db}
}

func (ir *IngredientsRepository) Insert(ctx context.Context, ingredient *models.Ingredient, recipeId int, tx pgx.Tx) error {
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
