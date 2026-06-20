package db

import (
	"context"
	"errors"
)

type Store interface {
	Create(ctx context.Context, model any) error
	FindByID(ctx context.Context, model any, id any) error
	FindAll(ctx context.Context, model any, opts QueryOptions) error
	Update(ctx context.Context, model any) error
	Delete(ctx context.Context, model any) error
}

type QueryOptions struct {
	Limit  int
	Offset int
	SortBy string
	Order  string // 'asc' | 'desc'
	Search string
}

var ErrResourceNotFound = errors.New("resource not found")
