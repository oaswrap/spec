package spec

import (
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/openapi-go"
)

type reflector interface {
	Add(method, path string, opts ...option.OperationOption)
	Spec() spec
	Validate() error
}

type spec interface {
	MarshalYAML() ([]byte, error)
	MarshalJSON() ([]byte, error)
}

// operationContext is an interface for managing operation contexts in OpenAPI specifications.
type operationContext interface {
	With(opts ...option.OperationOption) operationContext
	Set(opt option.OperationOption)
	build() openapi.OperationContext
}
