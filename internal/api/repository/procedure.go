// Package repository provides data access objects for interacting with the database.
package repository

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ProcedureRepository handles database operations related to recipe procedure steps.
// It provides methods to create, read, update, and delete procedure step records.
type ProcedureRepository struct {
	ConnectionPool *pgxpool.Pool // Database connection pool
}

// NewProcedureRepository creates a new instance of ProcedureRepository.
// It requires a database connection pool to perform database operations.
func NewProcedureRepository(pool *pgxpool.Pool) *ProcedureRepository {
	return &ProcedureRepository{ConnectionPool: pool}
}

// Insert adds a new procedure step to the database within a transaction.
// It requires a context, the procedure step text, the associated recipe ID, and an active transaction.
// Returns an error if the insertion fails.
func (pr *ProcedureRepository) Insert(ctx context.Context, procedureStep string, recipeID int, tx pgx.Tx) error {
	log.Printf("Inside of ProcedureRepository.Insert")
	log.Printf("Inserting procedure step: %v", procedureStep)

	query := `
		INSERT INTO procedure_steps (
			step, recipe_id, created_by, created_date, updated_by, updated_date
		) VALUES ($1, $2, $3, $4, $5, $6)
		`

	log.Printf("recipeID: %v", recipeID)
	_, err := tx.Exec(ctx, query, procedureStep, recipeID, 1, time.Now(), 1, time.Now())
	if err != nil {
		log.Printf("Error inserting procedure step: %v", err)
		return err
	}

	log.Printf("Successfully inserted procedure step: %v", procedureStep)
	return nil
}

// GetProcedureByRecipeId retrieves all procedure steps for a recipe from the database.
// It requires a context and the ID of the recipe to retrieve.
// Returns a slice of procedure steps and an error if the retrieval fails.
func (pr * ProcedureRepository) GetProcedureByRecipeId(ctx context.Context, recipeID int) ([]string, error) {
	log.Println("Inside of GetProcedureByRecipeID")
	log.Printf("Retrieving procedure steps for recipe with ID: %v", recipeID)

	connection , err := pr.ConnectionPool.Acquire(ctx)

	if err != nil {
		log.Printf("Error getting a connection from the connection pool: %v\n", err)
		return nil, err
	}

	defer connection.Release()

	query := `SELECT step FROM procedure_steps WHERE recipe_id = $1`

	result, err := connection.Query(ctx, query, recipeID)

	if err != nil {
		log.Printf("Something went wrong with the following query: %v\n", query)
		return nil, err
	}

	var procedureSteps []string

	for result.Next() {
		
		var procedureStep string

		err := result.Scan(&procedureStep)

		if err != nil {
			log.Printf("Error scanning procedure step: %v\n", err)
			return nil, err
		}
		
		procedureSteps = append(procedureSteps, procedureStep)
	}

	if result.Err() != nil {
		log.Printf("Error retrieving procedure steps: %v\n", result.Err())
		return nil, err
	}

	return procedureSteps, nil
}

// TODO - implement
func(pr *ProcedureRepository) Get(ctx context.Context, recipeID int) ([]string, error) {
	log.Printf("Retrieving procedure steps for recipe with ID: %v", recipeID)
	
	return nil, nil
}


