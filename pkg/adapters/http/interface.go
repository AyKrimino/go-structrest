package http

import "context"

type Context interface {
	RequestContext() context.Context
	Param(name string) string
	Query(name string) string
	Bind(dest any) error
	JSON(code int, data any)
}

type Router interface {
	Group(prefix string) Router
	GET(path string, handler HandlerFunc)
	POST(path string, handler HandlerFunc)
	PUT(path string, handler HandlerFunc)
	DELETE(path string, handler HandlerFunc)
}

type HandlerFunc func(Context)
