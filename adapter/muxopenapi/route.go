package muxopenapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

type route struct {
	muxRoute   *mux.Route
	specRoute  spec.Route
	specRouter spec.Router

	pathPrefix string
}

var _ Route = (*route)(nil)

func (r *route) GetError() error {
	return r.muxRoute.GetError()
}

func (r *route) GetHandler() http.Handler {
	return r.muxRoute.GetHandler()
}

func (r *route) GetHostTemplate() (string, error) {
	return r.muxRoute.GetHostTemplate()
}

func (r *route) GetMethods() ([]string, error) {
	return r.muxRoute.GetMethods()
}

func (r *route) GetName() string {
	return r.muxRoute.GetName()
}

func (r *route) Handler(handler http.Handler) Route {
	r.muxRoute.Handler(handler)
	return r
}

func (r *route) HandlerFunc(handler func(http.ResponseWriter, *http.Request)) Route {
	r.muxRoute.HandlerFunc(handler)
	return r
}

func (r *route) Headers(pairs ...string) Route {
	r.muxRoute.Headers(pairs...)
	return r
}

func (r *route) Host(tpl string) Route {
	r.muxRoute.Host(tpl)
	return r
}

func (r *route) Methods(methods ...string) Route {
	r.muxRoute.Methods(methods...)
	if len(methods) > 0 {
		r.specRoute.With(option.Method(methods[0]))
	}
	return r
}

func (r *route) Name(name string) Route {
	r.muxRoute.Name(name)
	return r
}

func (r *route) Path(tpl string) Route {
	r.muxRoute.Path(tpl)
	r.specRoute.With(option.Path(tpl))
	return r
}

func (r *route) PathPrefix(tpl string) Route {
	r.muxRoute.PathPrefix(tpl)
	r.pathPrefix = tpl
	return r
}

func (r *route) Queries(queries ...string) Route {
	r.muxRoute.Queries(queries...)
	return r
}

func (r *route) Schemes(schemes ...string) Route {
	r.muxRoute.Schemes(schemes...)
	return r
}

func (r *route) SkipClean() bool {
	return r.muxRoute.SkipClean()
}

func (r *route) Subrouter() Router {
	return &router{
		muxRouter:  r.muxRoute.Subrouter(),
		specRouter: r.specRouter.Group(r.pathPrefix),
	}
}

func (r *route) With(opts ...option.OperationOption) Route {
	r.specRoute.With(opts...)
	return r
}
