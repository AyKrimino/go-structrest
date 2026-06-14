package core

import (
	"context"
	"errors"
	"fmt"
	"testing"
)

type testUser struct{}

func (u *testUser) BeforeCreate(ctx context.Context) error {
	fmt.Println("Hashing password...")
	return nil
}

func (u *testUser) BeforeDelete(ctx context.Context) error {
	return errors.New("cannot delete user")
}

type testPost struct{}

func TestBeforeCreate(t *testing.T) {
	err := RunHook(context.Background(), &testUser{}, BeforeCreate)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}

func TestBeforeDelete(t *testing.T) {
	err := RunHook(context.Background(), &testUser{}, BeforeDelete)

	if err == nil {
		t.Fatal("expected error")
	}

	if err.Error() != "cannot delete user" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNoHooks(t *testing.T) {
	err := RunHook(context.Background(), &testPost{}, BeforeCreate)

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
}
