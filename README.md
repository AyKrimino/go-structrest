# go-structrest

[![Go Version](https://img.shields.io/badge/Go-1.26+-00ADD8?logo=go&logoColor=white)](https://go.dev)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Overview

**go-structrest** is a generic, reflection-based CRUD framework for Go.

It eliminates repetitive repository, service, and handler boilerplate for standard REST APIs. Define a struct, register it, and get a fully paginated, searchable REST API with lifecycle hooks.

## Installation

```bash
go get github.com/AyKrimino/go-structrest
```

## Quick Start

Define a struct with `crud` tags:

```go
package main

type User struct {
    ID    int    `json:"id" bun:",pk,autoincrement" crud:"pk,autoincrement"`
    Name  string `json:"name" bun:",notnull" crud:"searchable"`
    Email string `json:"email" bun:",notnull" crud:"searchable"`
}
```

Wire it up in `main.go`:

```go
package main

import (
    "context"
    "database/sql"
    "log"

    bunStore "github.com/AyKrimino/go-structrest/pkg/adapters/db/bun"
    ginAdapter "github.com/AyKrimino/go-structrest/pkg/adapters/http/gin"
    structrest "github.com/AyKrimino/go-structrest/pkg/api"
    "github.com/gin-gonic/gin"
    "github.com/uptrace/bun"
    "github.com/uptrace/bun/dialect/sqlitedialect"
    "github.com/uptrace/bun/driver/sqliteshim"
)

func main() {
    ctx := context.Background()

    sqldb, _ := sql.Open(sqliteshim.ShimName, "file::memory:")
    db := bun.NewDB(sqldb, sqlitedialect.New())
    db.NewCreateTable().Model((*User)(nil)).Exec(ctx)

    engine := gin.Default()
    root := ginAdapter.NewGinRouter(&engine.RouterGroup)
    store := bunStore.NewBunStore(db)
    api := structrest.NewAPI(root, store)

    api.Register("/users", User{})
    engine.Run(":8080")
}
```

Verify it works:

```bash
# Create a user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Alice", "email": "alice@example.com"}'

# List all users (paginated)
curl http://localhost:8080/users?page=1&limit=10
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
