package main

import (
	"log"

	ginAdapter "github.com/AyKrimino/go-structrest/pkg/adapters/http/gin"
	structrest "github.com/AyKrimino/go-structrest/pkg/api"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()

	root := ginAdapter.NewGinRouter(&engine.RouterGroup)

	store := NewDummyStore()

	api := structrest.NewAPI(root, store)
	err := api.Register("/users", User{})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("server listening on :8080")

	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
