// Package repository provides data access objects for interacting with the database.
package repository

// Repository defines a generic interface for database operations.
// It provides a standard set of methods that all repositories should implement.
// The generic type parameter T represents the model type that the repository handles.
type Repository[T any] interface {
	// Insert adds a new item to the database.
	// Returns an error if the insertion fails.
	Insert() error
	
	// Get retrieves an item by its ID.
	// Returns the item and an error if the retrieval fails.
	Get(id int) (T, error)
	
	// GetAll retrieves all items.
	// Returns a slice of items and an error if the retrieval fails.
	GetAll() ([]T, error)
	
	// Update modifies an existing item in the database.
	// Returns an error if the update fails.
	Update(item T) error
	
	// Delete removes an item from the database.
	// Returns an error if the deletion fails.
	Delete(item T) error
}
