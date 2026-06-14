package core

import (
	"log/slog"
	goHTTP "net/http"
	"reflect"

	"github.com/AyKrimino/go-structrest/pkg/adapters/db"
	"github.com/AyKrimino/go-structrest/pkg/adapters/http"
)

type ResourceHandler struct {
	bluePrint *BlueprintModel
	store     db.Store
}

func NewResourceHandler(bluePrint *BlueprintModel, store db.Store) *ResourceHandler {
	return &ResourceHandler{
		bluePrint: bluePrint,
		store:     store,
	}
}

func (h *ResourceHandler) HandleCreate(ctx http.Context) {
	var err error

	freshInstance := h.bluePrint.NewInstance()

	err = ctx.Bind(freshInstance)
	if err != nil {
		slog.Error("failed to bind request", "error", err)
		ctx.JSON(goHTTP.StatusBadRequest, err)
		return
	}

	err = RunHook(ctx.RequestContext(), freshInstance, BeforeCreate)
	if err != nil {
		slog.Error("failed to run before create hook", "error", err)
		ctx.JSON(goHTTP.StatusConflict, err)
		return
	}

	err = h.store.Create(ctx.RequestContext(), freshInstance)
	if err != nil {
		slog.Error("failed to create resource", "error", err)
		ctx.JSON(goHTTP.StatusInternalServerError, err)
		return
	}

	err = RunHook(ctx.RequestContext(), freshInstance, AfterCreate)
	if err != nil {
		slog.Warn("failed to run after create hook", "warning", err)
	}

	ctx.JSON(goHTTP.StatusCreated, freshInstance)
}
