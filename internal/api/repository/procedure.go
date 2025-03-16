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
	db *pgxpool.Pool // Database connection pool
}

// NewProcedureRepository creates a new instance of ProcedureRepository.
// It requires a database connection pool to perform database operations.
func NewProcedureRepository(db *pgxpool.Pool) *ProcedureRepository {
	return &ProcedureRepository{db: db}
}

// Insert adds a new procedure step to the database within a transaction.
// It requires a context, the procedure step text, the associated recipe ID, and an active transaction.
// Returns an error if the insertion fails.
func (r *ProcedureRepository) Insert(ctx context.Context, procedureStep string, recipeID int, tx pgx.Tx) error {
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
