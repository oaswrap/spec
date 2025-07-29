package core

import "github.com/swaggest/openapi-go"

type Reflector interface {
	AddOperation(oc openapi.OperationContext) error
	NewOperationContext(method, path string) (openapi.OperationContext, error)
	Spec() Spec
}

type Spec interface {
	MarshalYAML() ([]byte, error)
	MarshalJSON() ([]byte, error)
}
