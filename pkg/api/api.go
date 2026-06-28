package api

import (
	"github.com/AyKrimino/go-structrest/internal/core"
	"github.com/AyKrimino/go-structrest/pkg/adapters/db"
	"github.com/AyKrimino/go-structrest/pkg/adapters/http"
)

// API is the main entry point for registering models and generating RESTful endpoints.
type API struct {
	router http.Router
	store  db.Store
}

// NewAPI initializes a new API instance with the specified HTTP router and database store.
func NewAPI(router http.Router, store db.Store) *API {
	return &API{
		router: router,
		store:  store,
	}
}

// Register generates standard CRUD endpoints (GET, POST, PUT, DELETE) for the provided model under the specified prefix.
func (a *API) Register(prefix string, model any) error {
	blueprint, err := core.BuildBlueprint(model)
	if err != nil {
		return err
	}

	handler := core.NewResourceHandler(blueprint, a.store)

	routerGroup := a.router.Group(prefix)
	{
		routerGroup.POST("", handler.HandleCreate)
		routerGroup.GET("/:id", handler.HandleGet)
		routerGroup.PUT("/:id", handler.HandleUpdate)
		routerGroup.DELETE("/:id", handler.HandleDelete)
		routerGroup.GET("", handler.HandleList)
	}

	return nil
}
