package fiberopenapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapters/fiberopenapi/internal/constant"
	"github.com/oaswrap/spec/adapters/fiberopenapi/internal/handler"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/parser"
	"github.com/oaswrap/spec/pkg/util"
)

// NewGenerator creates a new OpenAPI generator with the specified Fiber router and options.
//
// It initializes the OpenAPI router and sets up the necessary routes for OpenAPI documentation.
func NewGenerator(r fiber.Router, opts ...option.OpenAPIOption) Generator {
	return NewRouter(r, opts...)
}

// NewRouter creates a new OpenAPI router with the specified Fiber router and options.
//
// It initializes the OpenAPI generator and sets up the necessary routes for OpenAPI documentation.
func NewRouter(r fiber.Router, opts ...option.OpenAPIOption) Generator {
	defaultOpts := []option.OpenAPIOption{
		option.WithTitle(constant.DefaultTitle),
		option.WithDescription(constant.DefaultDescription),
		option.WithVersion(constant.DefaultVersion),
		option.WithDocsPath(constant.DefaultDocsPath),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
		option.WithPathParser(parser.NewColonParamParser()),
	}
	opts = append(defaultOpts, opts...)
	gen := spec.NewGenerator(opts...)
	cfg := gen.Config()

	rr := &router{
		fiberRouter: r,
		specRouter:  gen,
		gen:         gen,
	}

	// If docs are disabled, return the router without adding docs routes.
	if cfg.DisableDocs {
		return rr
	}

	handler := handler.NewOpenAPIHandler(cfg, gen)
	openapiPath := util.JoinPath(cfg.DocsPath, constant.OpenAPIFileName)

	r.Get(cfg.DocsPath, handler.Docs)
	r.Get(openapiPath, handler.OpenAPIYaml)

	return rr
}

type router struct {
	fiberRouter fiber.Router
	specRouter  spec.Router
	gen         spec.Generator
}

func (r *router) Use(args ...any) Router {
	r.fiberRouter.Use(args...)
	return r
}

func (r *router) Get(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodGet, path, handler...)
}

func (r *router) Head(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodHead, path, handler...)
}

func (r *router) Post(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodPost, path, handler...)
}

func (r *router) Put(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodPut, path, handler...)
}

func (r *router) Patch(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodPatch, path, handler...)
}

func (r *router) Delete(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodDelete, path, handler...)
}

func (r *router) Connect(path string, handler ...fiber.Handler) Route {
	fr := r.fiberRouter.Connect(path, handler...)

	return &route{fr: fr}
}

func (r *router) Options(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodOptions, path, handler...)
}

func (r *router) Trace(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodTrace, path, handler...)
}

func (r *router) Add(method, path string, handler ...fiber.Handler) Route {
	fr := r.fiberRouter.Add(method, path, handler...)
	sr := r.specRouter.Add(method, path)

	route := &route{
		fr: fr,
		sr: sr,
	}

	return route
}

func (r *router) Static(prefix, root string, config ...fiber.Static) Router {
	r.fiberRouter.Static(prefix, root, config...)
	return r
}

func (r *router) Group(prefix string, handlers ...fiber.Handler) Router {
	rr := r.fiberRouter.Group(prefix, handlers...)
	sr := r.specRouter.Group(prefix)

	return &router{
		fiberRouter: rr,
		specRouter:  sr,
	}
}

func (r *router) Route(prefix string, fn func(router Router)) Router {
	fr := r.fiberRouter.Group(prefix)
	sr := r.specRouter.Group(prefix)

	subRouter := &router{
		fiberRouter: fr,
		specRouter:  sr,
	}

	fn(subRouter)

	return subRouter
}

func (r *router) With(opts ...option.GroupOption) Router {
	r.specRouter.Use(opts...)
	return r
}

func (r *router) Validate() error {
	return r.gen.Validate()
}

func (r *router) GenerateOpenAPISchema(formats ...string) ([]byte, error) {
	return r.gen.GenerateSchema(formats...)
}

func (r *router) MarshalYAML() ([]byte, error) {
	return r.gen.MarshalYAML()
}

func (r *router) MarshalJSON() ([]byte, error) {
	return r.gen.MarshalJSON()
}

func (r *router) WriteSchemaTo(path string) error {
	return r.gen.WriteSchemaTo(path)
}
