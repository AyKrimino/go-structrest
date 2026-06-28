
# go-structrest

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

**go-structrest** is a generic, reflection-based CRUD framework for Go.

It eliminates repetitive repository, service, and handler boilerplate for standard REST APIs. Define a struct, register it, and get a fully paginated, searchable REST API with lifecycle hooks.

> **⚠️ Development Status:** This project is under active development. The core features (CRUD operations, pagination, sorting, searching, and lifecycle hooks) are fully tested and stable for standard, flat data models. However, advanced ORM features like nested data structures (e.g., `has_many`, `belongs_to`) or complex many-to-many relationships are not yet fully supported or tested. The library works perfectly for the standard examples presented below.

## Installation

```bash
go get github.com/AyKrimino/go-structrest
```

## Quick Start

Define a struct with `crud` tags and your ORM's tags (e.g., `bun`):

```go
type Album struct {
	ID     int64   `json:"id" bun:"id,pk,autoincrement" crud:"pk"`
	Title  string  `json:"title" bun:",notnull" crud:"searchable"`
	Artist string  `json:"artist" bun:",notnull" crud:"searchable"`
	Price  float64 `json:"price" bun:",notnull"`
}
```

Wire it up in `main.go` (example using PostgreSQL and Bun):

```go
package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	bunStore "github.com/AyKrimino/go-structrest/pkg/adapters/db/bun"
	ginAdapter "github.com/AyKrimino/go-structrest/pkg/adapters/http/gin"
	structrest "github.com/AyKrimino/go-structrest/pkg/api"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
)

func main() {
	// 1. Initialize Database (PostgreSQL)
	// Note: Replace with your actual connection string
	sqldb, err := sql.Open("postgres", "postgres://user:password@localhost:5432/your_db?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db := bun.NewDB(sqldb, pgdialect.New())
	ctx := context.Background()

	// Create table automatically (for demo purposes)
	_, err = db.NewCreateTable().Model((*Album)(nil)).IfNotExists().Exec(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Initialize Structrest
	engine := gin.Default()
	root := ginAdapter.NewGinRouter(&engine.RouterGroup)
	store := bunStore.NewBunStore(db)
	api := structrest.NewAPI(root, store)

	// 3. Register the model
	err = api.Register("/albums", Album{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server is running on http://localhost:8080")
	engine.Run(":8080")
}
```

Verify it works:

```bash
# Create an album
curl -X POST http://localhost:8080/albums \
  -H "Content-Type: application/json" \
  -d '{"title": "Nevermind", "artist": "Nirvana", "price": 18.75}'

# List all albums (paginated and searchable)
curl "http://localhost:8080/albums?page=1&limit=10&search=nirvana"
```

## Lifecycle Hooks

Implement any of the six hook interfaces directly on your struct. Hooks run automatically during the corresponding CRUD operation.

| Interface | When it runs |
|-----------|-------------|
| `BeforeCreate(ctx) error` | After JSON binding, before insert |
| `AfterCreate(ctx) error` | After successful insert |
| `BeforeUpdate(ctx) error` | After finding + binding, before update |
| `AfterUpdate(ctx) error` | After successful update |
| `BeforeDelete(ctx) error` | After finding, before delete |
| `AfterDelete(ctx) error` | After successful delete |

Returning an error from a `Before` hook aborts the operation (HTTP 409). Errors from `After` hooks are logged but do not affect the response.

```go
func (acc *Account) BeforeCreate(ctx context.Context) error {
    hashedPassword, err := bcrypt.GenerateFromPassword(
        []byte(acc.Password),
        bcrypt.DefaultCost,
    )
    if err != nil {
        return err
    }
    acc.Password = string(hashedPassword)
    return nil
}
```

## Supported Adapters

### Database

| Adapter | Status |
|---------|--------|
| [Bun](https://bun.uptrace.dev/) (SQLite, Postgres, MySQL) | Supported |

### HTTP

| Adapter | Status |
|---------|--------|
| [Gin](https://github.com/gin-gonic/gin) | Supported |

Gorm and Chi support are planned.
