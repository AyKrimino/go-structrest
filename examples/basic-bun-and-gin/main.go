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

	// SQLite DB
	sqldb, err := sql.Open(sqliteshim.ShimName, "file::memory:")
	if err != nil {
		panic(err)
	}

	db := bun.NewDB(sqldb, sqlitedialect.New())

	db.NewCreateTable().Model((*User)(nil)).Exec(ctx)

	err = SeedUsers(ctx, db)
	if err != nil {
		log.Fatal(err)
	}

	engine := gin.Default()
	root := ginAdapter.NewGinRouter(&engine.RouterGroup)
	store := bunStore.NewBunStore(db)
	api := structrest.NewAPI(root, store)

	err = api.Register("/users", User{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server listening on :8080")

	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
