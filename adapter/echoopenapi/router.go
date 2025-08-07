package echoopenapi

import (
	"io/fs"

	"github.com/labstack/echo/v4"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/echoopenapi/internal/constant"
	"github.com/oaswrap/spec/adapter/echoopenapi/internal/handler"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/parser"
)

type router struct {
	echoGroup  *echo.Group
	specRouter spec.Router
	gen        spec.Generator
}

// NewGenerator creates a new OpenAPI generator with the provided Echo instance and options.
//
// It initializes the OpenAPI configuration and sets up the necessary routes for serving
func NewGenerator(e *echo.Echo, opts ...option.OpenAPIOption) Generator {
	return NewRouter(e, opts...)
}

// NewRouter creates a new OpenAPI router with the provided Echo instance and options.
//
// It initializes the OpenAPI configuration and sets up the necessary routes for serving
func NewRouter(e *echo.Echo, opts ...option.OpenAPIOption) Generator {
	defaultOpts := []option.OpenAPIOption{
		option.WithTitle(constant.DefaultTitle),
		option.WithDescription(constant.DefaultDescription),
		option.WithVersion(constant.DefaultVersion),
		option.WithDocsPath(constant.DefaultDocsPath),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
		option.WithPathParser(parser.NewColonParamParser()),
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

	handler := handler.NewHandler(gen)

	rr.echoGroup.GET(handler.DocsFilePath(), handler.DocsFile)
	rr.echoGroup.GET(handler.DocsPath(), handler.Docs)

	return rr
}

func (r *router) Add(method, path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	echoRoute := r.echoGroup.Add(method, path, handler, m...)
	route := &route{echoRoute: echoRoute}

	if method == echo.CONNECT {
		// CONNECT method is not supported by OpenAPI, so we skip it
		return route
	}
	route.specRoute = r.specRouter.Add(method, path)

	return route
}

func (r *router) GET(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(echo.GET, path, handler, m...)
}

func (r *router) POST(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(echo.POST, path, handler, m...)
}

func (r *router) PUT(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(echo.PUT, path, handler, m...)
}

func (r *router) DELETE(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(echo.DELETE, path, handler, m...)
}

func (r *router) PATCH(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(echo.PATCH, path, handler, m...)
}

func (r *router) HEAD(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(echo.HEAD, path, handler, m...)
}

func (r *router) OPTIONS(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(echo.OPTIONS, path, handler, m...)
}

func (r *router) TRACE(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(echo.TRACE, path, handler, m...)
}

func (r *router) CONNECT(path string, handler echo.HandlerFunc, m ...echo.MiddlewareFunc) Route {
	return r.Add(echo.CONNECT, path, handler, m...)
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

func (r *router) File(path, file string) {
	r.echoGroup.File(path, file)
}

func (r *router) FileFS(path, file string, fs fs.FS, m ...echo.MiddlewareFunc) {
	r.echoGroup.FileFS(path, file, fs, m...)
}

func (r *router) Static(prefix, root string) {
	r.echoGroup.Static(prefix, root)
}

func (r *router) StaticFS(prefix string, fs fs.FS) {
	r.echoGroup.StaticFS(prefix, fs)
}

func (r *router) With(opts ...option.GroupOption) Router {
	r.specRouter.Use(opts...)
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
