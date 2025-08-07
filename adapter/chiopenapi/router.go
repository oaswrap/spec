package chiopenapi

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/chiopenapi/internal/constant"
	"github.com/oaswrap/spec/module/specui"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

type router struct {
	chiRouter  chi.Router
	specRouter spec.Router
	gen        spec.Generator
}

var _ Router = (*router)(nil)

// NewRouter creates a new OpenAPI router with the specified Chi router and options.
//
// It initializes the OpenAPI generator and sets up the necessary handlers for OpenAPI documentation.
func NewRouter(r chi.Router, opts ...option.OpenAPIOption) Generator {
	return NewGenerator(r, opts...)
}

// NewGenerator creates a new OpenAPI generator with the specified Chi router and options.
//
// It initializes the OpenAPI configuration and sets up the necessary handlers for OpenAPI documentation.
func NewGenerator(r chi.Router, opts ...option.OpenAPIOption) Generator {
	defaultOpts := []option.OpenAPIOption{
		option.WithTitle(constant.DefaultTitle),
		option.WithDescription(constant.DefaultDescription),
		option.WithVersion(constant.DefaultVersion),
		option.WithDocsPath(constant.DefaultDocsPath),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
	}
	opts = append(defaultOpts, opts...)
	gen := spec.NewRouter(opts...)
	cfg := gen.Config()

	rr := &router{
		chiRouter:  r,
		specRouter: gen,
		gen:        gen,
	}

	if cfg.DisableDocs {
		return rr
	}

	handler := specui.NewHandler(gen)
	r.Get(handler.DocsPath(), handler.DocsFunc())
	r.Get(handler.DocsFilePath(), handler.DocsFileFunc())

	return rr
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chiRouter.ServeHTTP(w, req)
}

func (r *router) Use(middlewares ...func(http.Handler) http.Handler) {
	r.chiRouter.Use(middlewares...)
}

func (r *router) With(middlewares ...func(http.Handler) http.Handler) Router {
	cr := r.chiRouter.With(middlewares...)

	return &router{
		chiRouter:  cr,
		specRouter: r.specRouter,
		gen:        r.gen,
	}
}

func (r *router) Group(fn func(r Router), opts ...option.GroupOption) Router {
	var group *router
	r.chiRouter.Group(func(chiRouter chi.Router) {
		group = &router{
			chiRouter:  chiRouter,
			specRouter: r.specRouter.Group("/", opts...),
			gen:        r.gen,
		}
		fn(group)
	})
	return group
}

func (r *router) Route(pattern string, fn func(r Router), opts ...option.GroupOption) Router {
	var subRouter *router
	r.chiRouter.Route(pattern, func(chiRouter chi.Router) {
		subRouter = &router{
			chiRouter:  chiRouter,
			specRouter: r.specRouter.Group(pattern, opts...),
			gen:        r.gen,
		}
		fn(subRouter)
	})
	return subRouter
}

func (r *router) Handle(pattern string, h http.Handler) {
	r.chiRouter.Handle(pattern, h)
}

func (r *router) HandleFunc(pattern string, h http.HandlerFunc) {
	r.chiRouter.HandleFunc(pattern, h)
}

func (r *router) Mount(pattern string, h http.Handler) {
	r.chiRouter.Mount(pattern, h)
}

func (r *router) Method(method, pattern string, h http.Handler) Route {
	r.chiRouter.Method(method, pattern, h)
	if method == http.MethodConnect {
		// CONNECT method is not supported by OpenAPI, so we skip it
		return &route{}
	}
	sr := r.specRouter.Add(method, pattern)

	return &route{specRoute: sr}
}

func (r *router) MethodFunc(method, pattern string, h http.HandlerFunc) Route {
	r.chiRouter.MethodFunc(method, pattern, h)
	if method == http.MethodConnect {
		// CONNECT method is not supported by OpenAPI, so we skip it
		return &route{}
	}
	sr := r.specRouter.Add(method, pattern)

	return &route{specRoute: sr}
}

func (r *router) Connect(pattern string, h http.HandlerFunc) Route {
	return r.MethodFunc(http.MethodConnect, pattern, h)
}

func (r *router) Delete(pattern string, h http.HandlerFunc) Route {
	return r.MethodFunc(http.MethodDelete, pattern, h)
}

func (r *router) Get(pattern string, h http.HandlerFunc) Route {
	return r.MethodFunc(http.MethodGet, pattern, h)
}

func (r *router) Head(pattern string, h http.HandlerFunc) Route {
	return r.MethodFunc(http.MethodHead, pattern, h)
}

func (r *router) Options(pattern string, h http.HandlerFunc) Route {
	return r.MethodFunc(http.MethodOptions, pattern, h)
}

func (r *router) Patch(pattern string, h http.HandlerFunc) Route {
	return r.MethodFunc(http.MethodPatch, pattern, h)
}

func (r *router) Post(pattern string, h http.HandlerFunc) Route {
	return r.MethodFunc(http.MethodPost, pattern, h)
}

func (r *router) Put(pattern string, h http.HandlerFunc) Route {
	return r.MethodFunc(http.MethodPut, pattern, h)
}

func (r *router) Trace(pattern string, h http.HandlerFunc) Route {
	return r.MethodFunc(http.MethodTrace, pattern, h)
}

func (r *router) NotFound(h http.HandlerFunc) {
	r.chiRouter.NotFound(h)
}

func (r *router) MethodNotAllowed(h http.HandlerFunc) {
	r.chiRouter.MethodNotAllowed(h)
}

func (r *router) UseOptions(opts ...option.GroupOption) Router {
	r.specRouter.Use(opts...)
	return r
}

func (r *router) Validate() error {
	return r.gen.Validate()
}

func (r *router) GenerateSchema(formats ...string) ([]byte, error) {
	return r.gen.GenerateSchema(formats...)
}

func (r *router) MarshalYAML() ([]byte, error) {
	return r.gen.MarshalYAML()
}

func (r *router) MarshalJSON() ([]byte, error) {
	return r.gen.MarshalJSON()
}

func (r *router) WriteSchemaTo(filename string) error {
	return r.gen.WriteSchemaTo(filename)
}
