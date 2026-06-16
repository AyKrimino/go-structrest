package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/AyKrimino/go-structrest/pkg/adapters/db"
)

type DummyStore struct {
	users []User
}

func NewDummyStore() *DummyStore {
	return &DummyStore{
		users: []User{
			{
				ID:    1,
				Name:  "John Doe",
				Email: "jL2eF@example.com",
			},
		},
	}
}

func (s *DummyStore) Create(ctx context.Context, model any) error {
	fmt.Printf("[DummyStore] Create: %+v\n", model)
	return nil
}

func (s *DummyStore) FindByID(ctx context.Context, model any, id any) error {
	fmt.Printf("[DummyStore] FindByID: %v\n", id)

	userID, ok := id.(int)
	if !ok {
		return db.ErrNotFound
	}

	for _, user := range s.users {
		if user.ID == userID {
			reflect.ValueOf(model).Elem().Set(
				reflect.ValueOf(user),
			)
			return nil
		}
	}

	return db.ErrNotFound
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
