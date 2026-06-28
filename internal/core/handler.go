package core

import (
	"errors"
	"fmt"
	"log/slog"
	goHTTP "net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/AyKrimino/go-structrest/pkg/adapters/db"
	"github.com/AyKrimino/go-structrest/pkg/adapters/http"
)

// ResourceHandler manages the HTTP request lifecycle for a specific model,
// coordinating between the HTTP context, business logic hooks, and the database store.
type ResourceHandler struct {
	bluePrint *BlueprintModel
	store     db.Store
}

// NewResourceHandler creates a new handler instance for the given model blueprint and database store.
func NewResourceHandler(bluePrint *BlueprintModel, store db.Store) *ResourceHandler {
	return &ResourceHandler{
		bluePrint: bluePrint,
		store:     store,
	}
}

// HandleCreate processes an HTTP POST request to create a new resource.
func (h *ResourceHandler) HandleCreate(ctx http.Context) {
	var err error

	freshInstance := h.bluePrint.NewInstance()

	err = ctx.Bind(freshInstance)
	if err != nil {
		slog.Error("failed to bind request", "error", err)
		ctx.JSON(goHTTP.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err = RunHook(ctx.RequestContext(), freshInstance, BeforeCreate)
	if err != nil {
		slog.Error("failed to run before create hook", "error", err)
		ctx.JSON(goHTTP.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = h.store.Create(ctx.RequestContext(), freshInstance)
	if err != nil {
		slog.Error("failed to create resource", "error", err)
		ctx.JSON(goHTTP.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = RunHook(ctx.RequestContext(), freshInstance, AfterCreate)
	if err != nil {
		slog.Warn("failed to run after create hook", "warning", err)
	}

	ctx.JSON(goHTTP.StatusCreated, freshInstance)
}

// HandleGet processes an HTTP GET request to retrieve a single resource by its ID.
func (h *ResourceHandler) HandleGet(ctx http.Context) {
	idStr := ctx.Param("id")

	id, err := h.parsePrimaryKey(idStr)
	if err != nil {
		slog.Error("failed to parse primary key", "error", err)
		ctx.JSON(goHTTP.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	freshInstance := h.bluePrint.NewInstance()

	err = h.store.FindByID(ctx.RequestContext(), freshInstance, id)
	if err != nil {
		slog.Error("failed to find resource", "error", err)
		if errors.Is(err, db.ErrResourceNotFound) {
			ctx.JSON(goHTTP.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		ctx.JSON(goHTTP.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(goHTTP.StatusOK, freshInstance)
}

// HandleUpdate processes an HTTP PUT request to update an existing resource.
func (h *ResourceHandler) HandleUpdate(ctx http.Context) {
	idStr := ctx.Param("id")

	id, err := h.parsePrimaryKey(idStr)
	if err != nil {
		slog.Error("failed to parse primary key", "error", err)
		ctx.JSON(goHTTP.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	freshInstance := h.bluePrint.NewInstance()

	err = h.store.FindByID(ctx.RequestContext(), freshInstance, id)
	if err != nil {
		slog.Error("failed to find resource", "error", err)
		if errors.Is(err, db.ErrResourceNotFound) {
			ctx.JSON(goHTTP.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		ctx.JSON(goHTTP.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = ctx.Bind(freshInstance)
	if err != nil {
		slog.Error("failed to bind request", "error", err)
		ctx.JSON(goHTTP.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	err = RunHook(ctx.RequestContext(), freshInstance, BeforeUpdate)
	if err != nil {
		slog.Error("failed to run before update hook", "error", err)
		ctx.JSON(goHTTP.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = h.store.Update(ctx.RequestContext(), freshInstance)
	if err != nil {
		slog.Error("failed to update resource", "error", err)
		ctx.JSON(goHTTP.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = RunHook(ctx.RequestContext(), freshInstance, AfterUpdate)
	if err != nil {
		slog.Warn("failed to run after update hook", "warning", err)
	}

	ctx.JSON(goHTTP.StatusOK, freshInstance)
}

// HandleDelete processes an HTTP DELETE request to remove a resource.
func (h *ResourceHandler) HandleDelete(ctx http.Context) {
	idStr := ctx.Param("id")

	id, err := h.parsePrimaryKey(idStr)
	if err != nil {
		slog.Error("failed to parse primary key", "error", err)
		ctx.JSON(goHTTP.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	freshInstance := h.bluePrint.NewInstance()

	err = h.store.FindByID(ctx.RequestContext(), freshInstance, id)
	if err != nil {
		slog.Error("failed to find resource", "error", err)
		if errors.Is(err, db.ErrResourceNotFound) {
			ctx.JSON(goHTTP.StatusNotFound, map[string]string{"error": err.Error()})
			return
		}
		ctx.JSON(goHTTP.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = RunHook(ctx.RequestContext(), freshInstance, BeforeDelete)
	if err != nil {
		slog.Error("failed to run before delete hook", "error", err)
		ctx.JSON(goHTTP.StatusConflict, map[string]string{"error": err.Error()})
		return
	}

	err = h.store.Delete(ctx.RequestContext(), freshInstance)
	if err != nil {
		slog.Error("failed to delete resource", "error", err)
		ctx.JSON(goHTTP.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = RunHook(ctx.RequestContext(), freshInstance, AfterDelete)
	if err != nil {
		slog.Warn("failed to run after delete hook", "warning", err)
	}

	ctx.JSON(goHTTP.StatusNoContent, nil)
}

// HandleList processes an HTTP GET request to retrieve a paginated list of resources.
func (h *ResourceHandler) HandleList(ctx http.Context) {
	queryOpts, err := h.buildQueryOptions(ctx)
	if err != nil {
		ctx.JSON(goHTTP.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	freshInstance := h.bluePrint.NewInstance()

	t := reflect.TypeOf(freshInstance)
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	result := reflect.New(reflect.SliceOf(t))

	err = h.store.FindAll(ctx.RequestContext(), result.Interface(), queryOpts)
	if err != nil {
		slog.Error("failed to find resources", "error", err)
		ctx.JSON(goHTTP.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	ctx.JSON(goHTTP.StatusOK, result.Interface())
}

// parsePrimaryKey converts a string ID from the URL into the correct Go type based on the model's primary key definition.
func (h *ResourceHandler) parsePrimaryKey(idStr string) (any, error) {
	var t reflect.Kind

	for _, field := range h.bluePrint.Fields {
		if field.PrimaryKey {
			t = field.Kind
			break
		}
	}

	switch t {
	case reflect.String:
		return idStr, nil

	case reflect.Int:
		v, err := strconv.ParseInt(idStr, 10, 0)
		if err != nil {
			slog.Error("failed to cast id to int", "error", err)
			return nil, err
		}
		return int(v), nil

	case reflect.Int8:
		v, err := strconv.ParseInt(idStr, 10, 8)
		if err != nil {
			slog.Error("failed to cast id to int8", "error", err)
			return nil, err
		}
		return int8(v), nil

	case reflect.Int16:
		v, err := strconv.ParseInt(idStr, 10, 16)
		if err != nil {
			slog.Error("failed to cast id to int16", "error", err)
			return nil, err
		}
		return int16(v), nil

	case reflect.Int32:
		v, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			slog.Error("failed to cast id to int32", "error", err)
			return nil, err
		}
		return int32(v), nil

	case reflect.Int64:
		v, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			slog.Error("failed to cast id to int64", "error", err)
			return nil, err
		}
		return v, nil

	case reflect.Uint:
		v, err := strconv.ParseUint(idStr, 10, 0)
		if err != nil {
			slog.Error("failed to cast id to uint", "error", err)
			return nil, err
		}
		return uint(v), nil

	case reflect.Uint8:
		v, err := strconv.ParseUint(idStr, 10, 8)
		if err != nil {
			slog.Error("failed to cast id to uint8", "error", err)
			return nil, err
		}
		return uint8(v), nil

	case reflect.Uint16:
		v, err := strconv.ParseUint(idStr, 10, 16)
		if err != nil {
			slog.Error("failed to cast id to uint16", "error", err)
			return nil, err
		}
		return uint16(v), nil

	case reflect.Uint32:
		v, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			slog.Error("failed to cast id to uint32", "error", err)
			return nil, err
		}
		return uint32(v), nil

	case reflect.Uint64:
		v, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			slog.Error("failed to cast id to uint64", "error", err)
			return nil, err
		}
		return v, nil

	default:
		slog.Error("unsupported primary key type", "type", t)
		return nil, fmt.Errorf("unsupported primary key type: %s", t)
	}
}

// isValidSortField checks if the provided field name exists in the model's blueprint.
func (h *ResourceHandler) isValidSortField(field string) bool {
	for _, f := range h.bluePrint.Fields {
		if strings.EqualFold(f.Name, field) {
			return true
		}
	}

	return false
}

// parseIntQuery safely parses an integer from a query string, returning a default value if empty.
func parseIntQuery(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}

	return strconv.Atoi(value)
}

// buildQueryOptions extracts and validates pagination, sorting, and search parameters from the HTTP request.
func (h *ResourceHandler) buildQueryOptions(ctx http.Context) (db.QueryOptions, error) {
	var opts db.QueryOptions

	page, err := parseIntQuery(ctx.Query("page"), 1)
	if err != nil {
		return opts, err
	}

	limit, err := parseIntQuery(ctx.Query("limit"), 10)
	if err != nil {
		return opts, err
	}

	opts.Offset = (page - 1) * limit
	opts.Limit = limit

	opts.Order = ctx.Query("order")
	if opts.Order != "asc" && opts.Order != "desc" {
		opts.Order = "asc"
	}

	opts.Search = ctx.Query("search")

	opts.SortBy = ctx.Query("sort_by")
	if h.isValidSortField(opts.SortBy) {
		opts.SortBy = h.store.GetColumnName(h.bluePrint.NewInstance(), opts.SortBy)
	} else {
		opts.SortBy = ""
	}

	var searchableCols []string
	for _, goField := range h.bluePrint.GetSearchableFields() {
		searchableCols = append(searchableCols, h.store.GetColumnName(h.bluePrint.NewInstance(), goField))
	}
	opts.SearchableFields = searchableCols

	return opts, nil
}
