// Package repository provides data access objects for interacting with the database.
package repository

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"recipe-generator/internal/api/model"
)

// RecipeRepository handles database operations related to recipes.
// It provides methods to create, read, update, and delete recipe records.
type RecipeRepository struct {
	ConnectionPool *pgxpool.Pool // Database connection pool
}

// NewRecipeRepository creates a new instance of RecipeRepository.
// It requires a database connection pool to perform database operations.
func NewRecipeRepository(pool *pgxpool.Pool) *RecipeRepository {
	return &RecipeRepository{ConnectionPool: pool}
}

// Insert adds a new recipe to the database within a transaction.
// It requires a context, the recipe model, and an active transaction.
// Returns the inserted recipe with its ID populated, or an error if the insertion fails.
func (r *RecipeRepository) Insert(ctx context.Context, model *model.Recipe,  transactionHandler pgx.Tx) (*model.Recipe, error) {
	log.Printf("Starting database insertion for recipe: %s", model.RecipeName)

	query := `
		INSERT INTO recipes (
			recipe_name, description, prep_time_minutes, 
			cook_time_minutes, servings, created_by,
			created_date, updated_by, updated_date
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9
		) RETURNING id`

	err := transactionHandler.QueryRow(
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

// Get retrieves a recipe from the database by its ID.
// It requires a context and the ID of the recipe to retrieve.
// Returns the recipe and an error if the retrieval fails.
// NOTE: This function's returned recipe does not have any ingredients or procedure steps set. 
func (r *RecipeRepository) Get(ctx context.Context, recipeID int) (*model.Recipe, error) {
	log.Printf("inside Get function of RecipeRepository")
	log.Printf("Retrieving recipe with ID: %d from recipe from database.  ", recipeID)

	// grab a connection from connection pool
	connection, err := r.ConnectionPool.Acquire(ctx)

	if err != nil {
		log.Printf("Error getting a connection from the connection pool: %v", err)
		return nil, err
	}

	// need to close connection when method returns.
	defer connection.Release()

	// take ID and retreive data from database.
	recipeQuery :=  `
		SELECT recipe_name, description, prep_time_minutes, cook_time_minutes, servings
		FROM recipes
		WHERE id = $1
	`
	result, err := connection.Query(ctx, recipeQuery, recipeID)

	if err != nil {
		log.Printf("Something went wrong with the following query: %v\n", recipeQuery)
		return nil, err
	}

	var recipe model.Recipe

	// scan the data into the recipe model
	if result.Next() {

		err = result.Scan(
			&recipe.RecipeName,
			&recipe.Description,
			&recipe.PrepTimeMinutes,
			&recipe.CookTimeMinutes,
			&recipe.Servings,
		)

		if err != nil {
			log.Printf("Error happened when retrieving the random ID to generate from database: ")
			return nil, err
		}
	}

	return &recipe, nil
}


// GetRandomRecipeId retrieves a random recipe ID from the database.
// It requires a context.
// Returns the random recipe ID and an error if the retrieval fails.
func (r *RecipeRepository) GetRandomRecipeId(ctx context.Context) (int, error) {

	// grab connection from connection pool
	connection, err := r.ConnectionPool.Acquire(ctx)

	if err != nil {
		log.Printf("Error getting a connection from the connection pool: %v", err)

		return 0, err
	}

	// release connection
	defer connection.Release()

	// first get a random ID from the database. If no recipes exist, then return an error.
	// then grab that ID from the database. 
	selectRecipeQuery := `
		WITH random_row AS (
		SELECT id FROM recipes 
		OFFSET floor(random() * (SELECT COUNT(*) FROM recipes))
		LIMIT 1)
		SELECT id FROM random_row
	`

	result, err := connection.Query(
		ctx,
		selectRecipeQuery,
	)

	// check if something went wrong with the query 
	if err != nil {
		log.Printf("Something went wrong with the following query: %v\n", selectRecipeQuery)

		return 0, err
	}
	
	// grab random ID from db
	var randomID int

	if result.Next() {

		// randomID is set by Scan method via a reference to the randomID variable
		err := result.Scan(&randomID)

		if err != nil {
			
			log.Printf("Error happened when retrieving the random ID to generate from database: ")
			return 0, err
		}
	}

	return randomID, nil
}
