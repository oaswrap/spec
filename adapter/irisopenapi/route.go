package irisopenapi

import (
	irisrouter "github.com/kataras/iris/v12/core/router"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

type route struct {
	irisRoute *irisrouter.Route
	specRoute spec.Route
}

var _ Route = (*route)(nil)

func (r *route) Method() string {
	if r.irisRoute == nil {
		return ""
	}
	return r.irisRoute.Method
}

func (r *route) Path() string {
	if r.irisRoute == nil {
		return ""
	}
	return r.irisRoute.Path
}

func (r *route) Name() string {
	if r.irisRoute == nil {
		return ""
	}
	return r.irisRoute.Name
}

func (r *route) With(opts ...option.OperationOption) Route {
	if r.specRoute == nil {
		return r
	}

	r.specRoute.With(opts...)
	return r
}
