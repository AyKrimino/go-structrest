package gin_adapter

import (
	"github.com/AyKrimino/go-structrest/pkg/adapters/http"
	"github.com/gin-gonic/gin"
)

type GinRouter struct {
	group *gin.RouterGroup
}

func NewGinRouter(group *gin.RouterGroup) http.Router {
	return &GinRouter{
		group: group,
	}
}

func (r *GinRouter) Group(prefix string) http.Router {
	return &GinRouter{
		group: r.group.Group(prefix),
	}
}

func (r *GinRouter) GET(path string, handler http.HandlerFunc) {
	r.group.GET(path, wrapHandler(handler))
}

func (r *GinRouter) POST(path string, handler http.HandlerFunc) {
	r.group.POST(path, wrapHandler(handler))
}

func (r *GinRouter) PUT(path string, handler http.HandlerFunc) {
	r.group.PUT(path, wrapHandler(handler))
}

func (r *GinRouter) DELETE(path string, handler http.HandlerFunc) {
	r.group.DELETE(path, wrapHandler(handler))
}

func wrapHandler(handler http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(&GinContext{
			ginContext: c,
		})
	}
}
