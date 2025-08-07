package httpopenapi

import (
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

type route struct {
	specRoute spec.Route
}

var _ Route = (*route)(nil)

func (r *route) With(opts ...option.OperationOption) Route {
	if r.specRoute == nil {
		return r
	}
	r.specRoute.With(opts...)
	return r
}
