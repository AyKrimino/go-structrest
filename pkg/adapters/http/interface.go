package http

import "context"

// Context abstracts the underlying HTTP framework's context, allowing the core engine
// to interact with requests and responses without direct dependencies.
type Context interface {
	// RequestContext returns the standard Go context associated with the HTTP request.
	RequestContext() context.Context

	// Param extracts a path parameter from the URL (e.g., :id).
	Param(name string) string

	// Query extracts a query parameter from the URL (e.g., ?search=abc).
	Query(name string) string

	// Bind parses the request body into the provided destination struct.
	Bind(dest any) error

	// JSON sends a JSON response with the specified status code and data.
	JSON(code int, data any)
}

// Router abstracts the underlying HTTP router, allowing for dynamic route registration.
type Router interface {
	// Group creates a new sub-router with the specified prefix.
	Group(prefix string) Router

	// GET registers a handler for HTTP GET requests.
	GET(path string, handler HandlerFunc)

	// POST registers a handler for HTTP POST requests.
	POST(path string, handler HandlerFunc)

	// PUT registers a handler for HTTP PUT requests.
	PUT(path string, handler HandlerFunc)

	// DELETE registers a handler for HTTP DELETE requests.
	DELETE(path string, handler HandlerFunc)
}

// HandlerFunc defines the signature for handling HTTP requests within the framework.
type HandlerFunc func(Context)
