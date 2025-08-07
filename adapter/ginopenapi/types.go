package ginopenapi

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/option"
)

// Generator defines an Gin-compatible OpenAPI generator.
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

// Router defines an OpenAPI-aware Gin router.
//
// It wraps Gin routes and supports OpenAPI metadata.
type Router interface {
	// Handle registers a new route with the given method, path, and handler.
	Handle(method string, path string, handlers ...gin.HandlerFunc) Route

	// GET registers a new GET route.
	GET(path string, handlers ...gin.HandlerFunc) Route

	// POST registers a new POST route.
	POST(path string, handlers ...gin.HandlerFunc) Route

	// DELETE registers a new DELETE route.
	DELETE(path string, handlers ...gin.HandlerFunc) Route

	// PATCH registers a new PATCH route.
	PATCH(path string, handlers ...gin.HandlerFunc) Route

	// PUT registers a new PUT route.
	PUT(path string, handlers ...gin.HandlerFunc) Route

	// OPTIONS registers a new OPTIONS route.
	OPTIONS(path string, handlers ...gin.HandlerFunc) Route

	// HEAD registers a new HEAD route.
	HEAD(path string, handlers ...gin.HandlerFunc) Route

	// Group creates a new sub-group with the given prefix and middleware.
	Group(prefix string, handlers ...gin.HandlerFunc) Router

	// Use adds global middleware.
	Use(middlewares ...gin.HandlerFunc) Router

	// StaticFile serves a single static file.
	StaticFile(path string, filepath string) Router

	// StaticFileFS serves a static file from the given filesystem.
	StaticFileFS(path string, filepath string, fs http.FileSystem) Router

	// Static serves static files from a directory under the given prefix.
	Static(path string, root string) Router

	// StaticFS serves static files from the given filesystem.
	StaticFS(path string, fs http.FileSystem) Router

	// With applies OpenAPI group options to this router.
	With(opts ...option.GroupOption) Router
}

// Route represents a single Echo route with OpenAPI metadata.
type Route interface {
	// With applies OpenAPI operation options to this route.
	With(opts ...option.OperationOption) Route
}
