package main

import (
	"context"
	"fmt"
)

type DummyStore struct{}

func (s *DummyStore) Create(ctx context.Context, model any) error {
	fmt.Printf("[DummyStore] Create: %+v\n", model)
	return nil
}

func (s *DummyStore) FindByID(ctx context.Context, model any, id any) error {
	fmt.Printf("[DummyStore] FindByID: %v\n", id)
	return nil
}

func (s *DummyStore) FindAll(ctx context.Context, model any) error {
	fmt.Println("[DummyStore] FindAll")
	return nil
}

func (s *DummyStore) Update(ctx context.Context, model any) error {
	fmt.Printf("[DummyStore] Update: %+v\n", model)
	return nil
}

func (s *DummyStore) Delete(ctx context.Context, model any) error {
	fmt.Printf("[DummyStore] Delete: %+v\n", model)
	return nil
}
