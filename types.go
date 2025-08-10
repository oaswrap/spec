package spec

import (
	"net/http"

	specopenapi "github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/openapi-go"
)

// Generator defines an interface for building and exporting OpenAPI specifications.
type Generator interface {
	Router

	// Config returns the OpenAPI configuration used by the Generator.
	Config() *specopenapi.Config

	// GenerateSchema generates the OpenAPI schema in the specified format.
	// By default, it generates YAML. Pass "json" to generate JSON instead.
	GenerateSchema(formats ...string) ([]byte, error)

	// MarshalYAML returns the OpenAPI specification marshaled as YAML.
	MarshalYAML() ([]byte, error)

	// MarshalJSON returns the OpenAPI specification marshaled as JSON.
	MarshalJSON() ([]byte, error)

	// Validate checks whether the OpenAPI specification is valid.
	Validate() error

	// WriteSchemaTo writes the OpenAPI schema to a file.
	// The format is inferred from the file extension: ".yaml" for YAML, ".json" for JSON.
	WriteSchemaTo(path string) error

	// DocsHandlerFunc returns a handler for serving the OpenAPI documentation.
	DocsHandlerFunc() http.HandlerFunc

	// SpecHandlerFunc returns a handler for serving the OpenAPI specification.
	SpecHandlerFunc() http.HandlerFunc
}

// Router defines methods for registering API routes and operations
// in an OpenAPI specification. It lets you describe HTTP methods, paths, and options.
type Router interface {
	// Get registers a GET operation for the given path and options.
	Get(path string, opts ...option.OperationOption) Route

	// Post registers a POST operation for the given path and options.
	Post(path string, opts ...option.OperationOption) Route

	// Put registers a PUT operation for the given path and options.
	Put(path string, opts ...option.OperationOption) Route

	// Delete registers a DELETE operation for the given path and options.
	Delete(path string, opts ...option.OperationOption) Route

	// Patch registers a PATCH operation for the given path and options.
	Patch(path string, opts ...option.OperationOption) Route

	// Options registers an OPTIONS operation for the given path and options.
	Options(path string, opts ...option.OperationOption) Route

	// Head registers a HEAD operation for the given path and options.
	Head(path string, opts ...option.OperationOption) Route

	// Trace registers a TRACE operation for the given path and options.
	Trace(path string, opts ...option.OperationOption) Route

	// Add registers an operation for the given HTTP method, path, and options.
	Add(method, path string, opts ...option.OperationOption) Route

	// NewRoute creates a new route with the given options.
	NewRoute(opts ...option.OperationOption) Route

	// Route registers a nested route under the given pattern.
	// The provided function receives a Router to define sub-routes.
	Route(pattern string, fn func(router Router), opts ...option.GroupOption) Router

	// Group creates a new sub-router with the given path prefix and group options.
	Group(pattern string, opts ...option.GroupOption) Router

	// Use applies one or more group options to the router.
	Use(opts ...option.GroupOption) Router
}

// Route represents a single API route in the OpenAPI specification.
type Route interface {
	// Method sets the HTTP method for the route.
	Method(method string) Route
	// Path sets the HTTP path for the route.
	Path(path string) Route
	// With applies additional operation options to the route.
	With(opts ...option.OperationOption) Route
}

type reflector interface {
	Add(method, path string, opts ...option.OperationOption)
	Spec() spec
	Validate() error
}

type spec interface {
	MarshalYAML() ([]byte, error)
	MarshalJSON() ([]byte, error)
}

type operationContext interface {
	With(opts ...option.OperationOption) operationContext
	build() openapi.OperationContext
}
