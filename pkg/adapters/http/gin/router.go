package gin_adapter

import (
	"github.com/AyKrimino/go-structrest/pkg/adapters/http"
	"github.com/gin-gonic/gin"
)

// GinRouter adapts the Gin framework's router to work with structrest.
type GinRouter struct {
	group *gin.RouterGroup
}

// NewGinRouter creates a new GinRouter instance from an existing Gin RouterGroup.
func NewGinRouter(group *gin.RouterGroup) http.Router {
	return &GinRouter{
		group: group,
	}
}

// Group creates a new sub-router with the specified prefix.
func (r *GinRouter) Group(prefix string) http.Router {
	return &GinRouter{
		group: r.group.Group(prefix),
	}
}

// GET registers a handler for HTTP GET requests on the given path.
func (r *GinRouter) GET(path string, handler http.HandlerFunc) {
	r.group.GET(path, wrapHandler(handler))
}

// POST registers a handler for HTTP POST requests on the given path.
func (r *GinRouter) POST(path string, handler http.HandlerFunc) {
	r.group.POST(path, wrapHandler(handler))
}

// PUT registers a handler for HTTP PUT requests on the given path.
func (r *GinRouter) PUT(path string, handler http.HandlerFunc) {
	r.group.PUT(path, wrapHandler(handler))
}

// DELETE registers a handler for HTTP DELETE requests on the given path.
func (r *GinRouter) DELETE(path string, handler http.HandlerFunc) {
	r.group.DELETE(path, wrapHandler(handler))
}

// wrapHandler translates a structrest HandlerFunc into a native Gin HandlerFunc.
func wrapHandler(handler http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(&GinContext{
			ginContext: c,
		})
	}
}
