package httpopenapi

import (
	"net/http"

	"github.com/oaswrap/spec/option"
)

// Generator is an interface for generating OpenAPI documentation for HTTP applications.
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

// Router is an interface for handling HTTP routes with OpenAPI support.
type Router interface {
	http.Handler

	// HandleFunc registers a handler function for the specified pattern.
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) Route

	// Handle registers a handler for the specified pattern.
	Handle(pattern string, handler http.Handler) Route

	// Group creates a new sub-router with the specified options.
	Group(prefix string, mux *http.ServeMux, middlewares ...func(http.Handler) http.Handler) Router

	// With applies group level options to the router.
	With(opts ...option.GroupOption) Router
}

// Route is an interface for defining a route with OpenAPI options.
type Route interface {
	// With applies operation options to the route.
	With(opts ...option.OperationOption) Route
}
