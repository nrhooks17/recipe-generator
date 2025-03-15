package repository

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProcedureRepository struct {
	db *pgxpool.Pool
}

// Constructor for the ProcedureRepository
func NewProcedureRepository(db *pgxpool.Pool) *ProcedureRepository {
	return &ProcedureRepository{db: db}
}

// Insert method for submitting a step in a procedure. Uses a transaction to insert the step instead of a pool.
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
