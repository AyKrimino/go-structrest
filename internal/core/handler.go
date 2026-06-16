package core

import (
	"errors"
	"log/slog"
	goHTTP "net/http"
	"reflect"
	"strconv"

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

func (h *ResourceHandler) HandleGet(ctx http.Context) {
	var (
		err error
		id  any
	)

	idStr := ctx.Param("id")

	var t reflect.Kind
	for _, field := range h.bluePrint.Fields {
		if field.PrimaryKey {
			t = field.Kind
		}
	}

	// TODO: handle all the other types
	switch t {
	case reflect.Int:
		id, err = strconv.Atoi(idStr)
		if err != nil {
			slog.Error("failed to cast id to int", "error", err)
			ctx.JSON(goHTTP.StatusBadRequest, err)
			return
		}
	case reflect.String:
		id = idStr
	default:
		slog.Error("unknown primary key type", "type", t)
		ctx.JSON(goHTTP.StatusBadRequest, err)
		return
	}

	freshInstance := h.bluePrint.NewInstance()

	err = h.store.FindByID(ctx.RequestContext(), freshInstance, id)
	if err != nil {
		slog.Error("failed to find resource", "error", err)
		if errors.Is(err, db.ErrNotFound) {
			ctx.JSON(goHTTP.StatusNotFound, err)
			return
		}
		ctx.JSON(goHTTP.StatusInternalServerError, err)
		return
	}

	ctx.JSON(goHTTP.StatusOK, freshInstance)
}
