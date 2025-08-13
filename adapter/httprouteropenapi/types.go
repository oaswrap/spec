package httprouteropenapi

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/oaswrap/spec/option"
)

// Generator is an interface for generating OpenAPI specifications.
type Generator interface {
	Router

	// GenerateSchema generates the OpenAPI schema for the router.
	GenerateSchema(formats ...string) ([]byte, error)

	// MarshalJSON marshals the schema to JSON.
	MarshalJSON() ([]byte, error)

	// MarshalYAML marshals the schema to YAML.
	MarshalYAML() ([]byte, error)

	// Validate validates the schema.
	Validate() error

	// WriteSchemaTo writes the schema to a file.
	WriteSchemaTo(path string) error
}

// Router is an interface for handling HTTP requests.
type Router interface {
	http.Handler

	// Handle registers a new route with the given method, path, and handler.
	Handle(method, path string, handle httprouter.Handle) Route
	// Handler registers a new route with the given method, path, and handler.
	Handler(method, path string, handler http.Handler) Route
	// HandlerFunc registers a new route with the given method, path, and handler function.
	HandlerFunc(method, path string, handler http.HandlerFunc) Route

	// GET registers a new GET route with the given path and handler.
	GET(path string, handle httprouter.Handle) Route
	// POST registers a new POST route with the given path and handler.
	POST(path string, handle httprouter.Handle) Route
	// PUT registers a new PUT route with the given path and handler.
	PUT(path string, handle httprouter.Handle) Route
	// DELETE registers a new DELETE route with the given path and handler.
	DELETE(path string, handle httprouter.Handle) Route
	// PATCH registers a new PATCH route with the given path and handler.
	PATCH(path string, handle httprouter.Handle) Route
	// HEAD registers a new HEAD route with the given path and handler.
	HEAD(path string, handle httprouter.Handle) Route
	// OPTIONS registers a new OPTIONS route with the given path and handler.
	OPTIONS(path string, handle httprouter.Handle) Route

	// Group creates a new route group with the given prefix and middlewares.
	Group(prefix string, middlewares ...func(http.Handler) http.Handler) Router

	// Lookup retrieves the route for the given method and path.
	Lookup(method, path string) (httprouter.Handle, httprouter.Params, bool)
	// ServeFiles serves static files from the given root.
	ServeFiles(path string, root http.FileSystem)

	// With adds the given options to the group.
	With(opts ...option.GroupOption) Router
}

// Route is an interface for handling route-specific options.
type Route interface {
	// With adds the given options to the route.
	With(opts ...option.OperationOption) Route
}
