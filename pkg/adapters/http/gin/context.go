package gin_adapter

import (
	"context"

	"github.com/gin-gonic/gin"
)

// GinContext adapts Gin's *gin.Context to satisfy the http.Context interface.
type GinContext struct {
	ginContext *gin.Context
}

// RequestContext returns the standard Go context from the underlying Gin request.
func (ctx *GinContext) RequestContext() context.Context {
	return ctx.ginContext.Request.Context()
}

// Param extracts a path parameter from the Gin context.
func (ctx *GinContext) Param(name string) string {
	return ctx.ginContext.Param(name)
}

// Query extracts a query parameter from the Gin context.
func (ctx *GinContext) Query(name string) string {
	return ctx.ginContext.Query(name)
}

// Bind parses the JSON request body into the provided destination using Gin's ShouldBindJSON.
func (ctx *GinContext) Bind(dest any) error {
	return ctx.ginContext.ShouldBindJSON(dest)
}

// JSON sends a JSON response using Gin's native JSON method.
func (ctx *GinContext) JSON(code int, data any) {
	ctx.ginContext.JSON(code, data)
}
