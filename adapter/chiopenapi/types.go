package chiopenapi

import (
	"net/http"

	"github.com/oaswrap/spec/option"
)

// Generator is an interface that defines methods for generating OpenAPI schemas.
type Generator interface {
	Router

	// Validate checks if the OpenAPI schema is valid.
	Validate() error

	// GenerateSchema generates the OpenAPI schema in the specified formats.
	// Supported formats include "yaml", "json", etc.
	// If no formats are specified, it defaults to "yaml".
	GenerateSchema(formats ...string) ([]byte, error)

	// MarshalYAML generates the OpenAPI schema in YAML format.
	MarshalYAML() ([]byte, error)

	// MarshalJSON generates the OpenAPI schema in JSON format.
	MarshalJSON() ([]byte, error)

	// WriteSchemaTo writes the OpenAPI schema to a file in the specified format.
	WriteSchemaTo(filename string) error
}

// Router is an interface that defines methods for handling HTTP routes with OpenAPI support.
type Router interface {
	http.Handler
	// Use applies middleware to the router.
	Use(middlewares ...func(http.Handler) http.Handler)

	// With applies middleware to the router and returns a new Router instance.
	With(middlewares ...func(http.Handler) http.Handler) Router

	// Group creates a new sub-router with the specified options.
	Group(fn func(r Router), opts ...option.GroupOption) Router

	// Route creates a new sub-router for the specified pattern.
	Route(pattern string, fn func(r Router), opts ...option.GroupOption) Router

	// Mount mounts a sub-router at the specified pattern.
	Mount(pattern string, h http.Handler)

	// Handle registers a handler for the specified pattern.
	Handle(pattern string, h http.Handler)
	HandleFunc(pattern string, h http.HandlerFunc)

	// Method registers a handler for the specified HTTP method and pattern.
	Method(method, pattern string, h http.Handler) Route
	MethodFunc(method, pattern string, h http.HandlerFunc) Route

	// HTTP methods
	Connect(pattern string, h http.HandlerFunc) Route
	Delete(pattern string, h http.HandlerFunc) Route
	Get(pattern string, h http.HandlerFunc) Route
	Head(pattern string, h http.HandlerFunc) Route
	Options(pattern string, h http.HandlerFunc) Route
	Patch(pattern string, h http.HandlerFunc) Route
	Post(pattern string, h http.HandlerFunc) Route
	Put(pattern string, h http.HandlerFunc) Route
	Trace(pattern string, h http.HandlerFunc) Route

	// NotFound sets the handler for 404 Not Found responses.
	NotFound(h http.HandlerFunc)
	// MethodNotAllowed sets the handler for 405 Method Not Allowed responses.
	MethodNotAllowed(h http.HandlerFunc)

	// UseOptions applies OpenAPI group options to this router.
	UseOptions(opts ...option.GroupOption) Router
}

// Route represents a single Chi route with OpenAPI metadata.
type Route interface {
	// With applies OpenAPI operation options to this route.
	With(opts ...option.OperationOption) Route
}
