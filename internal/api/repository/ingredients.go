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
	ConnectionPool *pgxpool.Pool // Database connection pool
}

// NewIngredientsRepository creates a new instance of IngredientsRepository.
// It requires a database connection pool to perform database operations.
func NewIngredientsRepository(pool *pgxpool.Pool) *IngredientsRepository {
	return &IngredientsRepository{ConnectionPool: pool}
}

// Insert adds a new ingredient to the database within a transaction.
// It requires a context, the ingredient model, the associated recipe ID, and an active transaction.
// Returns an error if the insertion fails.
func (ir *IngredientsRepository) Insert(ctx context.Context, ingredient *model.Ingredient, recipeId int, tx pgx.Tx) error {
	log.Printf("Inside of IngredientsRepository.Insert")
	log.Printf("Inserting ingredient: %s", ingredient.IngredientName)

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

// GetIngredientsByRecipeId retrieves all ingredients for a specific recipe from the database.
// It requires a context and the ID of the recipe to retrieve.
// Returns a slice of ingredients and an error if the retrieval fails.
func (ir *IngredientsRepository) GetIngredientsByRecipeId(ctx context.Context, recipeID int) ([]model.Ingredient, error) {
	log.Printf("inside Get IngredientsRepository.Get")
	log.Printf("Retrieving ingredients with ID: %d from database.  ", recipeID)

	connection, err := ir.ConnectionPool.Acquire(ctx)

	if err != nil {
		log.Printf("Error getting a connection from the connection pool: %v", err)
		return nil, err
	}

	// release connection after done using it.
	defer connection.Release()

	// debug
	log.Printf("This is the recipeID: %v\n", recipeID)

	query := `SELECT ingredient_name, unit_of_measurement, unit_amount FROM ingredients WHERE recipe_id = $1`
	

	// execute the query
	result, err := connection.Query(ctx, query, recipeID)

	if err != nil {
		log.Printf("Something went wrong with the following query: %v\n", query)
		return nil, err
	}

	var ingredients []model.Ingredient

	// fill up ingredients array
	for result.Next() {

		var ingredient model.Ingredient

		err := result.Scan(&ingredient.IngredientName, &ingredient.UnitOfMeasurement, &ingredient.Amount)

		if err != nil {
			log.Printf("Error scanning ingredients: %v", err)
			return nil, err
		}

		ingredients = append(ingredients, ingredient)	

	}
	
	// pgx library says that I need to check this. so I'm checking it :S
	if result.Err() != nil {
		log.Printf("Error getting ingredients: %v", result.Err())
		return nil, result.Err()
	}
	
	return ingredients, nil
}

// TODO will work on this later when needed. 
func Get(ctx context.Context, id int) (*model.Ingredient, error) {
	return nil, nil
}
