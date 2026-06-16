package main

import (
	"context"
	"fmt"
	"reflect"

	"github.com/AyKrimino/go-structrest/pkg/adapters/db"
)

type DummyStore struct {
	users map[int]User
}

func NewDummyStore() *DummyStore {
	return &DummyStore{
		users: map[int]User{
			1: {
				ID:    1,
				Name:  "John Doe",
				Email: "jL2eF@example.com",
			},
			2: {
				ID:    2,
				Name:  "Biggie Smalls",
				Email: "biggie@smalls.com",
			},
			3: {
				ID:    3,
				Name:  "Tupac Amaru Shakur",
				Email: "tupac@shakur.com",
			},
		},
	}
}

func (s *DummyStore) Create(ctx context.Context, model any) error {
	id := len(s.users) + 1

	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Pointer {
		return fmt.Errorf("model must be a pointer")
	}

	elem := v.Elem()

	idField := elem.FieldByName("ID")
	if idField.IsValid() && idField.CanSet() {
		idField.SetInt(int64(id))
	}

	user, ok := elem.Interface().(User)
	if !ok {
		return fmt.Errorf("expected User model")
	}

	s.users[id] = user

	fmt.Printf("[DummyStore] Create: %+v\n", user)
	return nil
}

func (s *DummyStore) FindByID(ctx context.Context, model any, id any) error {
	fmt.Printf("[DummyStore] FindByID: %v\n", id)

	userID, ok := id.(int)
	if !ok {
		return db.ErrNotFound
	}

	user, exists := s.users[userID]
	if !exists {
		return db.ErrNotFound
	}

	reflect.ValueOf(model).Elem().Set(
		reflect.ValueOf(user),
	)

	return nil
}

func (s *DummyStore) FindAll(ctx context.Context, model any) error {
	fmt.Println("[DummyStore] FindAll")
	return nil
}

func (s *DummyStore) Update(ctx context.Context, model any) error {
	user, ok := reflect.ValueOf(model).Elem().Interface().(User)
	if !ok {
		return fmt.Errorf("expected User model")
	}

	if _, exists := s.users[user.ID]; !exists {
		return db.ErrNotFound
	}

	s.users[user.ID] = user

	fmt.Printf("[DummyStore] Update: %+v\n", user)
	return nil
}

func (s *DummyStore) Delete(ctx context.Context, model any) error {
	user, ok := reflect.ValueOf(model).Elem().Interface().(User)
	if !ok {
		return fmt.Errorf("expected User model")
	}

	if _, exists := s.users[user.ID]; !exists {
		return db.ErrNotFound
	}

	delete(s.users, user.ID)

	fmt.Printf("[DummyStore] Delete: %+v\n", user)
	return nil
}
