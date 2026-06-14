package gin_adapter

import (
	"context"

	"github.com/gin-gonic/gin"
)

type GinContext struct {
	ginContext *gin.Context
}

func (ctx *GinContext) RequestContext() context.Context {
	return ctx.ginContext.Request.Context()
}

func (ctx *GinContext) Param(name string) string {
	return ctx.ginContext.Param(name)
}

func (ctx *GinContext) Query(name string) string {
	return ctx.ginContext.Query(name)
}

func (ctx *GinContext) Bind(dest any) error {
	return ctx.ginContext.ShouldBindJSON(dest)
}

func (ctx *GinContext) JSON(code int, data any) {
	ctx.ginContext.JSON(code, data)
}
