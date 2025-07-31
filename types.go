package spec

import (
	specopenapi "github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/openapi-go"
)

// Generator is an interface that defines methods for generating OpenAPI specifications.
type Generator interface {
	Router

	// Config returns the OpenAPI configuration used by the Generator.
	Config() *specopenapi.Config
	// GenerateSchema generates the OpenAPI schema in the specified format.
	// By default, it generates YAML. If "json" is specified, it generates JSON
	GenerateSchema(formats ...string) ([]byte, error)
	// MarshalYAML marshals the OpenAPI specification to YAML format.
	MarshalYAML() ([]byte, error)
	// MarshalJSON marshals the OpenAPI specification to JSON format.
	MarshalJSON() ([]byte, error)
	// Validate validates the OpenAPI specification.
	Validate() error
	// WriteSchemaTo writes the OpenAPI schema to the specified file path.
	// The file format is determined by the file extension: ".json" for JSON and ".yaml" for YAML.
	WriteSchemaTo(path string) error
}

// Router is an interface that defines methods for registering routes and operations
// in an OpenAPI specification. It allows for defining HTTP methods, paths, and operation options.
type Router interface {
	// Get registers a new GET operation with the specified path and options.
	Get(path string, opts ...option.OperationOption) Route
	// Post registers a new POST operation with the specified path and options.
	Post(path string, opts ...option.OperationOption) Route
	// Put registers a new PUT operation with the specified path and options.
	Put(path string, opts ...option.OperationOption) Route
	// Delete registers a new DELETE operation with the specified path and options.
	Delete(path string, opts ...option.OperationOption) Route
	// Patch registers a new PATCH operation with the specified path and options.
	Patch(path string, opts ...option.OperationOption) Route
	// Options registers a new OPTIONS operation with the specified path and options.
	Options(path string, opts ...option.OperationOption) Route
	// Head registers a new HEAD operation with the specified path and options.
	Head(path string, opts ...option.OperationOption) Route
	// Trace registers a new TRACE operation with the specified path and options.
	Trace(path string, opts ...option.OperationOption) Route
	// Add registers a new operation with the specified method and path.
	Add(method, path string, opts ...option.OperationOption) Route

	// Route registers a new route with the specified pattern and function.
	// The function receives a Router instance to define sub-routes.
	Route(pattern string, fn func(router Router), opts ...option.GroupOption) Router
	// Group creates a new sub-router with the specified prefix and options.
	Group(pattern string, opts ...option.GroupOption) Router
	// Use applies the provided options to the router.
	Use(opts ...option.GroupOption) Router
}

// Route defines a method for creating a new route with the specified options.
type Route interface {
	// With applies the provided operation options to the route.
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
