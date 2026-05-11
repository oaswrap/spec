package irisopenapi

import (
	"github.com/kataras/iris/v12/context"

	"github.com/oaswrap/spec/option"
)

// Generator defines an Iris-compatible OpenAPI generator.
type Generator interface {
	Router

	// Validate checks if the OpenAPI specification is valid.
	Validate() error

	// ValidateReport validates the schema and returns all findings.
	ValidateReport() error

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

// Router defines an OpenAPI-aware Iris router.
type Router interface {
	// Handle registers a new route with the given method, path, and handlers.
	Handle(method string, path string, handlers ...context.Handler) Route

	// Get registers a new GET route.
	Get(path string, handlers ...context.Handler) Route

	// Post registers a new POST route.
	Post(path string, handlers ...context.Handler) Route

	// Put registers a new PUT route.
	Put(path string, handlers ...context.Handler) Route

	// Delete registers a new DELETE route.
	Delete(path string, handlers ...context.Handler) Route

	// Patch registers a new PATCH route.
	Patch(path string, handlers ...context.Handler) Route

	// Head registers a new HEAD route.
	Head(path string, handlers ...context.Handler) Route

	// Options registers a new OPTIONS route.
	Options(path string, handlers ...context.Handler) Route

	// Trace registers a new TRACE route.
	Trace(path string, handlers ...context.Handler) Route

	// Connect registers a new CONNECT route.
	// Connect operations are emitted only for OpenAPI 3.2.0.
	Connect(path string, handlers ...context.Handler) Route

	// Party creates a new sub-party with the given prefix and middleware.
	Party(prefix string, handlers ...context.Handler) Router

	// Use adds middleware.
	Use(handlers ...context.Handler) Router

	// With applies OpenAPI group options to this router.
	With(opts ...option.GroupOption) Router
}

// Route represents a single Iris route with OpenAPI metadata.
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
