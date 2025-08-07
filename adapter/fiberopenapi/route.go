package fiberopenapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

// Route represents a single route in the OpenAPI specification.
type Route interface {
	// Name sets the name for the route.
	Name(name string) Route
	// With applies the given options to the route.
	With(opts ...option.OperationOption) Route
}

type route struct {
	fr fiber.Router
	sr spec.Route
}

// Name sets the name for the route.
func (r *route) Name(name string) Route {
	r.fr.Name(name)

	return r
}

// With applies the given options to the route.
func (r *route) With(opts ...option.OperationOption) Route {
	if r.sr == nil {
		return r
	}
	r.sr.With(opts...)

	return r
}
