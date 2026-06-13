package http

type Context interface {
	Param(name string) string
	Query(name string) string
	Bind(dest any) error
	JSON(code int, data any)
}

type Router interface {
	Group(prefix string)
	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
}

type HandlerFunc func(Context)
