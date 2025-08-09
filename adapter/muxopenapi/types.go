package muxopenapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oaswrap/spec/option"
)

type Generator interface {
	Router

	GenerateSchema(formats ...string) ([]byte, error)
	MarshalJSON() ([]byte, error)
	MarshalYAML() ([]byte, error)
	Validate() error
	WriteSchemaTo(path string) error
}

type Router interface {
	http.Handler

	Get(name string) Route
	GetRoute(name string) Route
	Handle(path string, handler http.Handler) Route
	HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) Route
	Headers(pairs ...string) Route
	Host(tpl string) Route
	Methods(methods ...string) Route
	Name(name string) Route
	NewRoute() Route
	Path(tpl string) Route
	PathPrefix(tpl string) Route
	Queries(queries ...string) Route
	Schemes(schemes ...string) Route
	SkipClean(value bool) Router
	StrictSlash(value bool) Router
	Use(middlewares ...mux.MiddlewareFunc) Router
	UseEncodedPath() Router

	With(opts ...option.GroupOption) Router
}

type Route interface {
	GetError() error
	GetHandler() http.Handler
	GetHostTemplate() (string, error)
	GetMethods() ([]string, error)
	GetName() string
	Handler(handler http.Handler) Route
	HandlerFunc(handler func(http.ResponseWriter, *http.Request)) Route
	Headers(pairs ...string) Route
	Host(tpl string) Route
	Methods(methods ...string) Route
	Name(name string) Route
	Path(tpl string) Route
	PathPrefix(tpl string) Route
	Queries(queries ...string) Route
	Schemes(schemes ...string) Route
	SkipClean() bool
	Subrouter() Router

	With(opts ...option.OperationOption) Route
}
