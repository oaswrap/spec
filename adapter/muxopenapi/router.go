package muxopenapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/muxopenapi/internal/constant"
	"github.com/oaswrap/spec/option"
)

type router struct {
	muxRouter  *mux.Router
	specRouter spec.Router
	gen        spec.Generator
}

var _ Generator = (*router)(nil)

// NewRouter creates a new OpenAPI router with the provided mux.Router instance and options.
//
// It initializes the OpenAPI configuration and sets up the necessary routes for serving
func NewRouter(mux *mux.Router, opts ...option.OpenAPIOption) Generator {
	return NewGenerator(mux, opts...)
}

// NewGenerator creates a new OpenAPI router with the provided mux.Router instance and options.
//
// It initializes the OpenAPI configuration and sets up the necessary routes for serving
func NewGenerator(mux *mux.Router, opts ...option.OpenAPIOption) Generator {
	defaultOpts := []option.OpenAPIOption{
		option.WithTitle(constant.DefaultTitle),
		option.WithDescription(constant.DefaultDescription),
		option.WithVersion(constant.DefaultVersion),
	}
	opts = append(defaultOpts, opts...)
	gen := spec.NewRouter(opts...)
	rr := &router{
		muxRouter:  mux,
		specRouter: gen,
		gen:        gen,
	}
	cfg := gen.Config()
	if cfg.DisableDocs {
		return rr
	}

	mux.HandleFunc(cfg.DocsPath, gen.DocsHandlerFunc()).Methods(http.MethodGet)
	mux.HandleFunc(cfg.SpecPath, gen.SpecHandlerFunc()).Methods(http.MethodGet)

	return rr
}

func (r *router) Get(name string) Route {
	muxRoute := r.muxRouter.Get(name)

	return &route{
		muxRoute:   muxRoute,
		specRouter: r.specRouter,
	}
}

func (r *router) GetRoute(name string) Route {
	muxRoute := r.muxRouter.GetRoute(name)

	return &route{
		muxRoute:   muxRoute,
		specRouter: r.specRouter,
	}
}

func (r *router) HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) Route {
	return r.NewRoute().Path(path).HandlerFunc(handler)
}

func (r *router) Handle(path string, handler http.Handler) Route {
	return r.NewRoute().Path(path).Handler(handler)
}

func (r *router) Headers(pairs ...string) Route {
	return r.NewRoute().Headers(pairs...)
}

func (r *router) Host(tpl string) Route {
	return r.NewRoute().Host(tpl)
}

func (r *router) Methods(methods ...string) Route {
	return r.NewRoute().Methods(methods...)
}

func (r *router) Name(name string) Route {
	return r.NewRoute().Name(name)
}

func (r *router) NewRoute() Route {
	return &route{
		muxRoute:   r.muxRouter.NewRoute(),
		specRoute:  r.specRouter.NewRoute(),
		specRouter: r.specRouter,
	}
}

func (r *router) Path(tpl string) Route {
	return r.NewRoute().Path(tpl)
}

func (r *router) PathPrefix(tpl string) Route {
	return r.NewRoute().PathPrefix(tpl)
}

func (r *router) Queries(queries ...string) Route {
	return r.NewRoute().Queries(queries...)
}

func (r *router) Schemes(schemes ...string) Route {
	return r.NewRoute().Schemes(schemes...)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.muxRouter.ServeHTTP(w, req)
}

func (r *router) SkipClean(value bool) Router {
	r.muxRouter.SkipClean(value)
	return r
}

func (r *router) StrictSlash(value bool) Router {
	r.muxRouter.StrictSlash(value)
	return r
}

func (r *router) Use(middlewares ...mux.MiddlewareFunc) Router {
	r.muxRouter.Use(middlewares...)
	return r
}

func (r *router) UseEncodedPath() Router {
	r.muxRouter.UseEncodedPath()
	return r
}

func (r *router) With(opts ...option.GroupOption) Router {
	r.specRouter.Use(opts...)
	return r
}

func (r *router) GenerateSchema(formats ...string) ([]byte, error) {
	return r.gen.GenerateSchema(formats...)
}

func (r *router) MarshalJSON() ([]byte, error) {
	return r.gen.MarshalJSON()
}

func (r *router) MarshalYAML() ([]byte, error) {
	return r.gen.MarshalYAML()
}

func (r *router) Validate() error {
	return r.gen.Validate()
}

func (r *router) WriteSchemaTo(path string) error {
	return r.gen.WriteSchemaTo(path)
}
