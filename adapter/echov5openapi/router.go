package echov5openapi

import (
	"io/fs"
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/oaswrap/spec"
	specui "github.com/oaswrap/spec-ui"
	"github.com/oaswrap/spec/adapter/echov5openapi/internal/constant"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/mapper"
	"github.com/oaswrap/spec/pkg/parser"
)

type router struct {
	echoGroup  *echo.Group
	specRouter spec.Router
	gen        spec.Generator
}

// NewRouter creates a new OpenAPI router with the provided Echo instance and options.
//
// It initializes the OpenAPI configuration and sets up the necessary routes for serving.
func NewRouter(e *echo.Echo, opts ...option.OpenAPIOption) Generator {
	return NewGenerator(e, opts...)
}

// NewGenerator creates a new OpenAPI generator with the provided Echo instance and options.
//
// It initializes the OpenAPI configuration and sets up the necessary routes for serving.
func NewGenerator(e *echo.Echo, opts ...option.OpenAPIOption) Generator {
	defaultOpts := []option.OpenAPIOption{
		option.WithTitle(constant.DefaultTitle),
		option.WithDescription(constant.DefaultDescription),
		option.WithVersion(constant.DefaultVersion),
		option.WithPathParser(parser.NewColonParamParser()),
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
		echoGroup:  e.Group(""),
		specRouter: gen,
		gen:        gen,
	}

	if cfg.DisableDocs {
		return rr
	}

	handler := specui.NewHandler(mapper.SpecUIOpts(gen)...)

	rr.echoGroup.GET(cfg.DocsPath, echo.WrapHandler(handler.Docs()))
	rr.echoGroup.GET(cfg.SpecPath, echo.WrapHandler(handler.Spec()))

	if handler.AssetsEnabled() {
		rr.echoGroup.GET(handler.AssetsPath()+"/*", echo.WrapHandler(handler.Assets()))
	}

	return rr
}

func (r *router) Add(method, path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	if handler == nil {
		handler = func(c *echo.Context) error {
			return c.NoContent(http.StatusNotImplemented)
		}
	}
	echoRoute := r.echoGroup.Add(method, path, handler, m...)
	route := &route{echoRoute: echoRoute}

	if method == http.MethodConnect {
		// CONNECT method is not supported by OpenAPI, so we skip it
		return route
	}
	route.specRoute = r.specRouter.Add(method, path)

	return route
}

func (r *router) GET(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(http.MethodGet, path, handler, m...)
}

func (r *router) POST(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(http.MethodPost, path, handler, m...)
}

func (r *router) PUT(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(http.MethodPut, path, handler, m...)
}

func (r *router) DELETE(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(http.MethodDelete, path, handler, m...)
}

func (r *router) PATCH(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(http.MethodPatch, path, handler, m...)
}

func (r *router) HEAD(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(http.MethodHead, path, handler, m...)
}

func (r *router) OPTIONS(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(http.MethodOptions, path, handler, m...)
}

func (r *router) TRACE(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(http.MethodTrace, path, handler, m...)
}

func (r *router) CONNECT(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(http.MethodConnect, path, handler, m...)
}

func (r *router) Group(prefix string, m ...echo.MiddlewareFunc) Router {
	group := r.echoGroup.Group(prefix, m...)
	specGroup := r.specRouter.Group(prefix)

	return &router{
		echoGroup:  group,
		specRouter: specGroup,
		gen:        r.gen,
	}
}

func (r *router) Use(m ...echo.MiddlewareFunc) Router {
	r.echoGroup.Use(m...)
	return r
}

func (r *router) File(path, file string, m ...echo.MiddlewareFunc) {
	r.echoGroup.File(path, file, m...)
}

func (r *router) FileFS(path, file string, fs fs.FS, m ...echo.MiddlewareFunc) {
	r.echoGroup.FileFS(path, file, fs, m...)
}

func (r *router) Static(prefix, root string, m ...echo.MiddlewareFunc) {
	r.echoGroup.Static(prefix, root, m...)
}

func (r *router) StaticFS(prefix string, fs fs.FS, m ...echo.MiddlewareFunc) {
	r.echoGroup.StaticFS(prefix, fs, m...)
}

func (r *router) With(opts ...option.GroupOption) Router {
	r.specRouter.With(opts...)
	return r
}

func (r *router) WriteSchemaTo(filepath string) error {
	return r.gen.WriteSchemaTo(filepath)
}

func (r *router) MarshalYAML() ([]byte, error) {
	return r.gen.MarshalYAML()
}

func (r *router) MarshalJSON() ([]byte, error) {
	return r.gen.MarshalJSON()
}

func (r *router) GenerateSchema(format ...string) ([]byte, error) {
	return r.gen.GenerateSchema(format...)
}

func (r *router) Validate() error {
	return r.gen.Validate()
}
