package main

import (
	"log"

	"github.com/AyKrimino/go-structrest/internal/core"
	ginAdapter "github.com/AyKrimino/go-structrest/pkg/adapters/http/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	root := ginAdapter.NewGinRouter(&engine.RouterGroup)

	blueprint, err := core.BuildBlueprint(User{})
	if err != nil {
		log.Fatal(err)
	}

	handler := core.NewResourceHandler(blueprint, &DummyStore{})

	root.POST("/users", handler.HandleCreate)

	log.Println("server listening on :8080")

	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
