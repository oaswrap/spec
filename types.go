package spec

import (
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/openapi-go"
)

// Reflector is an interface for generating OpenAPI specifications.
type Reflector interface {
	// Add registers a new operation with the specified method and path.
	// It returns an error if the operation cannot be added.
	Add(method, path string, opts ...option.OperationOption)
	// Spec returns the generated OpenAPI specification.
	// It can be marshaled to JSON or YAML.
	Spec() Spec

	// Validate checks the generated OpenAPI specification for errors.
	Validate() error
}

// Spec is an interface for OpenAPI specifications.
type Spec interface {
	// MarshalYAML marshals the specification to YAML format.
	MarshalYAML() ([]byte, error)
	// MarshalJSON marshals the specification to JSON format.
	MarshalJSON() ([]byte, error)
}

// OperationContext is an interface for managing operation contexts in OpenAPI specifications.
type OperationContext interface {
	With(opts ...option.OperationOption) OperationContext
	Set(opt option.OperationOption)
	build() openapi.OperationContext
}
