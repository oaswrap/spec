package ginopenapi

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapters/ginopenapi/internal/constant"
	"github.com/oaswrap/spec/adapters/ginopenapi/internal/handler"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/parser"
	"github.com/oaswrap/spec/pkg/util"
)

// NewGenerator returns a new OpenAPI generator for Gin.
//
// It sets up the OpenAPI configuration and prepares the routes for serving docs.
func NewGenerator(ginRouter gin.IRouter, opts ...option.OpenAPIOption) Generator {
	return NewRouter(ginRouter, opts...)
}

// NewRouter returns a new Gin router with OpenAPI support.
//
// It configures the OpenAPI generator and attaches the routes for serving docs.
func NewRouter(ginRouter gin.IRouter, opts ...option.OpenAPIOption) Generator {
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
		ginRouter:  ginRouter,
		specRouter: gen,
		gen:        gen,
	}
	if cfg.DisableDocs {
		return rr
	}

	handler := handler.NewOpenAPIHandler(cfg, gen)

	openapiPath := util.JoinPath(cfg.DocsPath, constant.OpenAPIFileName)
	ginRouter.GET(cfg.DocsPath, handler.Docs)
	ginRouter.GET(openapiPath, handler.OpenAPIYaml)

	return rr
}

type router struct {
	ginRouter  gin.IRouter
	specRouter spec.Router
	gen        spec.Generator
}

var _ Generator = &router{}

// Handle registers a new route with the specified method and path, and returns a Route object.
func (r *router) Handle(method string, path string, handlers ...gin.HandlerFunc) Route {
	gr := r.ginRouter.Handle(method, path, handlers...)
	route := &route{ginRoute: gr}

	if strings.EqualFold(method, http.MethodConnect) {
		// CONNECT method is not supported by OpenAPI, so we skip it
		return route
	}
	route.specRoute = r.specRouter.Add(method, path)

	return route
}

// GET registers a new GET route with the specified path and handlers.
func (r *router) GET(path string, handlers ...gin.HandlerFunc) Route {
	return r.Handle(http.MethodGet, path, handlers...)
}

// POST registers a new POST route with the specified path and handlers.
func (r *router) POST(path string, handlers ...gin.HandlerFunc) Route {
	return r.Handle(http.MethodPost, path, handlers...)
}

// DELETE registers a new DELETE route with the specified path and handlers.
func (r *router) DELETE(path string, handlers ...gin.HandlerFunc) Route {
	return r.Handle(http.MethodDelete, path, handlers...)
}

// PATCH registers a new PATCH route with the specified path and handlers.
func (r *router) PATCH(path string, handlers ...gin.HandlerFunc) Route {
	return r.Handle(http.MethodPatch, path, handlers...)
}

// PUT registers a new PUT route with the specified path and handlers.
func (r *router) PUT(path string, handlers ...gin.HandlerFunc) Route {
	return r.Handle(http.MethodPut, path, handlers...)
}

// OPTIONS registers a new OPTIONS route with the specified path and handlers.
func (r *router) OPTIONS(path string, handlers ...gin.HandlerFunc) Route {
	return r.Handle(http.MethodOptions, path, handlers...)
}

// HEAD registers a new HEAD route with the specified path and handlers.
func (r *router) HEAD(path string, handlers ...gin.HandlerFunc) Route {
	return r.Handle(http.MethodHead, path, handlers...)
}

// Group creates a new route group with the specified prefix and handlers.
func (r *router) Group(prefix string, handlers ...gin.HandlerFunc) Router {
	ginGroup := r.ginRouter.Group(prefix, handlers...)
	specGroup := r.specRouter.Group(prefix)

	return &router{
		ginRouter:  ginGroup,
		specRouter: specGroup,
	}
}

// Use adds middleware to the router.
func (r *router) Use(middlewares ...gin.HandlerFunc) Router {
	r.ginRouter.Use(middlewares...)

	return r
}

// StaticFile serves a single file at the specified path.
func (r *router) StaticFile(path string, filepath string) Router {
	r.ginRouter.StaticFile(path, filepath)

	return r
}

// StaticFileFS serves a single file at the specified path using the provided file system.
func (r *router) StaticFileFS(path string, filepath string, fs http.FileSystem) Router {
	r.ginRouter.StaticFileFS(path, filepath, fs)

	return r
}

// Static serves static files from the specified root directory.
func (r *router) Static(path string, root string) Router {
	r.ginRouter.Static(path, root)

	return r
}

// StaticFS serves static files from the specified file system at the given path.
func (r *router) StaticFS(path string, fs http.FileSystem) Router {
	r.ginRouter.StaticFS(path, fs)

	return r
}

// With applies the specified options to the router.
// It allows for additional configuration of the OpenAPI router.
func (r *router) With(opts ...option.GroupOption) Router {
	r.specRouter.Use(opts...)

	return r
}

// Validate checks if the OpenAPI specification is valid.
func (r *router) Validate() error {
	return r.gen.Validate()
}

// GenerateSchema generates the OpenAPI schema in the specified format(s).
func (r *router) GenerateSchema(formats ...string) ([]byte, error) {
	return r.gen.GenerateSchema(formats...)
}

// MarshalYAML marshals the OpenAPI specification to YAML format.
func (r *router) MarshalYAML() ([]byte, error) {
	return r.gen.MarshalYAML()
}

// MarshalJSON marshals the OpenAPI specification to JSON format.
func (r *router) MarshalJSON() ([]byte, error) {
	return r.gen.MarshalJSON()
}

// WriteSchemaTo writes the OpenAPI schema to the specified file path.
func (r *router) WriteSchemaTo(filepath string) error {
	return r.gen.WriteSchemaTo(filepath)
}
