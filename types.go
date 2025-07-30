package spec

import "github.com/swaggest/openapi-go"

type Reflector interface {
	AddOperation(oc OperationContext) error
	NewOperationContext(method, path string) (OperationContext, error)
	Spec() Spec
}

type Spec interface {
	MarshalYAML() ([]byte, error)
	MarshalJSON() ([]byte, error)
}

type OperationContext interface {
	openapi.OperationContext

	unwrap() openapi.OperationContext
}
