package main

import (
	"context"

	"github.com/uptrace/bun"
)

func SeedUsers(ctx context.Context, db *bun.DB) error {
	users := []User{
		{Name: "John Doe", Email: "john.doe@example.com"},
		{Name: "Jane Smith", Email: "jane.smith@example.com"},
		{Name: "Alice Johnson", Email: "alice.johnson@example.com"},
		{Name: "Bob Williams", Email: "bob.williams@example.com"},
		{Name: "Charlie Brown", Email: "charlie.brown@example.com"},
		{Name: "David Miller", Email: "david.miller@example.com"},
		{Name: "Emma Davis", Email: "emma.davis@example.com"},
		{Name: "Frank Wilson", Email: "frank.wilson@example.com"},
		{Name: "Grace Taylor", Email: "grace.taylor@example.com"},
		{Name: "Henry Anderson", Email: "henry.anderson@example.com"},
		{Name: "Isabella Thomas", Email: "isabella.thomas@example.com"},
		{Name: "Jack Moore", Email: "jack.moore@example.com"},
		{Name: "Katherine Martin", Email: "katherine.martin@example.com"},
		{Name: "Liam Jackson", Email: "liam.jackson@example.com"},
		{Name: "Mia White", Email: "mia.white@example.com"},
		{Name: "Noah Harris", Email: "noah.harris@example.com"},
		{Name: "Olivia Clark", Email: "olivia.clark@example.com"},
		{Name: "Peter Lewis", Email: "peter.lewis@example.com"},
		{Name: "Sophia Walker", Email: "sophia.walker@example.com"},
		{Name: "William Hall", Email: "william.hall@example.com"},
	}

	_, err := db.NewInsert().
		Model(&users).
		Exec(ctx)

	return err
}
