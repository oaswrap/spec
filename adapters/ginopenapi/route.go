package ginopenapi

import (
	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)


type route struct {
	ginRoute  gin.IRoutes
	specRoute spec.Route
}

var _ Route = &route{}

// With applies the specified options to the route.
func (r *route) With(opts ...option.OperationOption) Route {
	r.specRoute.With(opts...)
	return r
}
