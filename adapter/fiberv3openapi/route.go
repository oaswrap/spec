package fiberv3openapi

import (
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

// Route represents a single route in the OpenAPI specification.
type Route interface {
	// With applies the given options to the route.
	With(opts ...option.OperationOption) Route
}

type route struct {
	sr spec.Route
}

// With applies the given options to the route.
func (r *route) With(opts ...option.OperationOption) Route {
	if r.sr == nil {
		return r
	}
	r.sr.With(opts...)

	return r
}
