package db

import (
	"context"
	"errors"
)

// Store defines the contract for database operations. Implementations must handle
// ORM-specific logic for querying and persisting models.
type Store interface {
	// Create inserts a new record into the database.
	Create(ctx context.Context, model any) error

	// FindByID retrieves a single record by its primary key and populates the provided model.
	FindByID(ctx context.Context, model any, id any) error

	// FindAll retrieves a list of records based on the provided query options.
	FindAll(ctx context.Context, model any, opts QueryOptions) error

	// Update modifies an existing record in the database.
	Update(ctx context.Context, model any) error

	// Delete removes a record from the database.
	Delete(ctx context.Context, model any) error

	// GetColumnName resolves the actual database column name for a given Go struct field name.
	GetColumnName(model any, goFieldName string) string
}

// QueryOptions contains parameters for filtering, sorting, and paginating database queries.
type QueryOptions struct {
	Limit  int
	Offset int
	SortBy string
	Order  string // 'asc' | 'desc'
	Search string

	// SearchableFields contains the database column names that should be included in search queries.
	SearchableFields []string
}

// ErrResourceNotFound is returned when a requested resource does not exist in the database.
var ErrResourceNotFound = errors.New("resource not found")
