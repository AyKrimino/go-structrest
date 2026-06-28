package main

import (
	"context"
	"database/sql"
	"errors"
	"log"

	bunStore "github.com/AyKrimino/go-structrest/pkg/adapters/db/bun"
	ginAdapter "github.com/AyKrimino/go-structrest/pkg/adapters/http/gin"
	structrest "github.com/AyKrimino/go-structrest/pkg/api"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID       int     `json:"id" bun:",pk,autoincrement" crud:"pk"`
	Username string  `json:"username" bun:",notnull" crud:"searchable"`
	Password string  `json:"password" bun:",notnull"`
	Role     string  `json:"role" bun:",notnull" crud:"searchable"`
	Balance  float64 `json:"balance" bun:",notnull"`
}

func main() {
	ctx := context.Background()

	// SQLite DB
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:")
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	db.NewCreateTable().Model((*Account)(nil)).Exec(ctx)

	engine := gin.Default()
	root := ginAdapter.NewGinRouter(&engine.RouterGroup)
	store := bunStore.NewBunStore(db)
	api := structrest.NewAPI(root, store)

	err = api.Register("/accounts", Account{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server listening on :8080")

	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

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

func (acc *Account) BeforeDelete(ctx context.Context) error {
	if acc.Role == "admin" {
		return errors.New("cannot delete admin account")
	}
	return nil
}
