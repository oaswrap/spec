package irisopenapi

import (
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"

	"github.com/oaswrap/spec"
	specui "github.com/oaswrap/spec-ui"
	"github.com/oaswrap/spec/internal/mapper"
	"github.com/oaswrap/spec/internal/validate"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

type router struct {
	party      iris.Party
	specRouter spec.Router
	gen        spec.Generator
}

var _ Generator = (*router)(nil)

// NewRouter creates a new OpenAPI router with the provided Iris party and options.
func NewRouter(party iris.Party, opts ...option.OpenAPIOption) Generator {
	return NewGenerator(party, opts...)
}

// NewGenerator creates a new OpenAPI generator with the provided Iris party and options.
func NewGenerator(party iris.Party, opts ...option.OpenAPIOption) Generator {
	defaultOpts := []option.OpenAPIOption{
		option.WithTitle("Iris OpenAPI"),
		option.WithDescription("OpenAPI documentation for Iris applications"),
		option.WithVersion("1.0.0"),
		option.WithStoplightElements(),
		option.WithCacheAge(0),
		option.WithReflectorConfig(
			option.ParameterTagMapping(openapi.ParameterInPath, "param"),
		),
	}

	opts = append(defaultOpts, opts...)
	gen := spec.NewRouter(opts...)
	cfg := gen.Config()

	rr := &router{
		party:      party,
		specRouter: gen,
		gen:        gen,
	}
	if cfg.DisableDocs {
		return rr
	}

	handler := specui.NewHandler(mapper.SpecUIOpts(gen)...)

	party.Get(cfg.DocsPath, iris.FromStd(handler.Docs()))
	party.Get(cfg.SpecPath, iris.FromStd(handler.Spec()))

	if handler.AssetsEnabled() {
		party.Get(handler.AssetsPath()+"/{file:path}", iris.FromStd(handler.Assets()))
	}

	return rr
}

func (r *router) Handle(method, path string, handlers ...context.Handler) Route {
	irisRoute := r.party.Handle(method, path, handlers...)
	route := &route{irisRoute: irisRoute}

	if !validate.AllowsOperationMethod(r.gen.Config().OpenAPIVersion, method) {
		return route
	}
	route.specRoute = r.specRouter.Add(method, path)

	return route
}

func (r *router) Get(path string, handlers ...context.Handler) Route {
	return r.Handle(http.MethodGet, path, handlers...)
}

func (r *router) Post(path string, handlers ...context.Handler) Route {
	return r.Handle(http.MethodPost, path, handlers...)
}

func (r *router) Put(path string, handlers ...context.Handler) Route {
	return r.Handle(http.MethodPut, path, handlers...)
}

func (r *router) Delete(path string, handlers ...context.Handler) Route {
	return r.Handle(http.MethodDelete, path, handlers...)
}

func (r *router) Patch(path string, handlers ...context.Handler) Route {
	return r.Handle(http.MethodPatch, path, handlers...)
}

func (r *router) Head(path string, handlers ...context.Handler) Route {
	return r.Handle(http.MethodHead, path, handlers...)
}

func (r *router) Options(path string, handlers ...context.Handler) Route {
	return r.Handle(http.MethodOptions, path, handlers...)
}

func (r *router) Trace(path string, handlers ...context.Handler) Route {
	return r.Handle(http.MethodTrace, path, handlers...)
}

func (r *router) Connect(path string, handlers ...context.Handler) Route {
	return r.Handle(http.MethodConnect, path, handlers...)
}

func (r *router) Party(prefix string, handlers ...context.Handler) Router {
	subParty := r.party.Party(prefix, handlers...)
	subSpecRouter := r.specRouter.Group(prefix)

	return &router{
		party:      subParty,
		specRouter: subSpecRouter,
		gen:        r.gen,
	}
}

func (r *router) Use(handlers ...context.Handler) Router {
	r.party.Use(handlers...)
	return r
}

func (r *router) With(opts ...option.GroupOption) Router {
	r.specRouter.With(opts...)
	return r
}

func (r *router) Validate() error {
	return r.gen.Validate()
}

func (r *router) ValidateReport() error {
	return r.gen.ValidateReport()
}

func (r *router) GenerateSchema(format ...string) ([]byte, error) {
	return r.gen.GenerateSchema(format...)
}

func (r *router) MarshalYAML() ([]byte, error) {
	return r.gen.MarshalYAML()
}

func (r *router) MarshalJSON() ([]byte, error) {
	return r.gen.MarshalJSON()
}

func (r *router) WriteSchemaTo(filepath string) error {
	return r.gen.WriteSchemaTo(filepath)
}
