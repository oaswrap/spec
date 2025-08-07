package echoopenapi

import (
	"io/fs"

	"github.com/labstack/echo/v4"
	"github.com/oaswrap/spec/option"
)

// Generator defines an Echo-compatible OpenAPI generator.
//
// It combines routing and OpenAPI schema generation.
type Generator interface {
	Router

	// Validate checks if the OpenAPI specification is valid.
	Validate() error

	// GenerateSchema generates the OpenAPI schema.
	// Defaults to YAML. Pass "json" to generate JSON.
	GenerateSchema(format ...string) ([]byte, error)

	// MarshalYAML marshals the OpenAPI schema to YAML.
	MarshalYAML() ([]byte, error)

	// MarshalJSON marshals the OpenAPI schema to JSON.
	MarshalJSON() ([]byte, error)

	// WriteSchemaTo writes the schema to the given file.
	// The format is inferred from the file extension.
	WriteSchemaTo(filepath string) error
}

// Router defines an OpenAPI-aware Echo router.
//
// It wraps Echo routes and supports OpenAPI metadata.
type Router interface {
	// GET registers a new GET route.
	GET(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route

	// POST registers a new POST route.
	POST(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route

	// PUT registers a new PUT route.
	PUT(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route

	// DELETE registers a new DELETE route.
	DELETE(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route

	// PATCH registers a new PATCH route.
	PATCH(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route

	// HEAD registers a new HEAD route.
	HEAD(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route

	// OPTIONS registers a new OPTIONS route.
	OPTIONS(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route

	// TRACE registers a new TRACE route.
	TRACE(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route

	// CONNECT registers a new CONNECT route.
	CONNECT(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route

	// Add registers a new route with the given method, path, and handler.
	Add(method, path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route

	// Group creates a new sub-group with the given prefix and middleware.
	Group(prefix string, m ...echo.MiddlewareFunc) Router

	// Use adds global middleware.
	Use(m ...echo.MiddlewareFunc) Router

	// File serves a single static file.
	File(path, file string)

	// FileFS serves a static file from the given filesystem.
	FileFS(path, file string, fs fs.FS, m ...echo.MiddlewareFunc)

	// Static serves static files from a directory under the given prefix.
	Static(prefix, root string)

	// StaticFS serves static files from the given filesystem.
	StaticFS(prefix string, fs fs.FS)

	// With applies OpenAPI group options to this router.
	With(opts ...option.GroupOption) Router
}

// Route represents a single Echo route with OpenAPI metadata.
type Route interface {
	// Method returns the HTTP method (GET, POST, etc.).
	Method() string

	// Path returns the route path.
	Path() string

	// Name returns the route name.
	Name() string

	// With applies OpenAPI operation options to this route.
	With(opts ...option.OperationOption) Route
}
