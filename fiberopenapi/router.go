package fiberopenapi

import (
	"fmt"
	stdpath "path"
	"sync"

	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/faizlabs/openapi-wrapper/core"
	"github.com/faizlabs/openapi-wrapper/internal/errors"
	"github.com/faizlabs/openapi-wrapper/internal/util"
	"github.com/faizlabs/openapi-wrapper/option"
	"github.com/gofiber/fiber/v2"
	"github.com/swaggest/openapi-go"
)

// Router defines the interface for an OpenAPI router.
type Router interface {
	// Use applies middleware to the router.
	Use(args ...any) Router

	// Get registers a GET route.
	Get(path string, handler ...fiber.Handler) Route
	// Head registers a HEAD route.
	Head(path string, handler ...fiber.Handler) Route
	// Post registers a POST route.
	Post(path string, handler ...fiber.Handler) Route
	// Put registers a PUT route.
	Put(path string, handler ...fiber.Handler) Route
	// Patch registers a PATCH route.
	Patch(path string, handler ...fiber.Handler) Route
	// Delete registers a DELETE route.
	Delete(path string, handler ...fiber.Handler) Route
	// Connect registers a CONNECT route.
	Connect(path string, handler ...fiber.Handler) Route
	// Options registers an OPTIONS route.
	Options(path string, handler ...fiber.Handler) Route
	// Trace registers a TRACE route.
	Trace(path string, handler ...fiber.Handler) Route

	// Add registers a route with the specified method and path.
	Add(method, path string, handler ...fiber.Handler) Route
	// Static serves static files from the specified root directory.
	Static(prefix, root string, config ...fiber.Static) Router

	// Group creates a new sub-router with the specified prefix and handlers.
	// The prefix is prepended to all routes in the sub-router.
	Group(prefix string, handlers ...fiber.Handler) Router

	// Route creates a new sub-router with the specified prefix and applies options.
	Route(prefix string, fn func(router Router)) Router

	// With applies options to the router.
	// This allows you to configure tags, security, and visibility for the routes.
	With(opts ...option.RouteOption) Router

	// Validate checks for errors at OpenAPI router initialization.
	//
	// It returns an error if there are issues with the OpenAPI configuration.
	Validate() error

	// GenerateOpenAPISchema generates the OpenAPI schema in the specified format.
	// Supported formats are "json" and "yaml".
	// If no format is specified, "yaml" is used by default.
	GenerateOpenAPISchema(format ...string) ([]byte, error)

	WriteSchemaTo(filePath string) error
}

type router struct {
	prefix     string
	router     fiber.Router
	subRouters []*router
	routes     []*route
	core       *core.Generator

	errors *errors.MultiError
	logger openapiwrapper.Logger

	tags     []string
	security []option.RouteSecurityConfig
	hide     bool

	buildOnce  sync.Once
	schemaOnce sync.Once
	schema     []byte
}

func NewRouter(r fiber.Router, opts ...option.OpenAPIOption) Router {
	c := newConfig(opts...)

	rr := &router{
		prefix: "/",
		router: r,
		errors: &errors.MultiError{},
		logger: c.Logger,
	}
	// If OpenAPI is disabled, return the router without any OpenAPI functionality.
	// This allows the application to run without OpenAPI if desired.
	if c.DisableOpenAPI {
		return rr
	}

	generator, err := core.NewGenerator(c)
	if err != nil {
		rr.errors.Add(err)
		return rr
	}

	rr.core = generator
	r.Get(c.DocsPath, docsHandler(c))
	openapiPath := stdpath.Join(c.DocsPath, openapiFileName)
	r.Get(openapiPath, openAPIHandler(rr))

	return rr
}

func (r *router) Use(args ...any) Router {
	r.router.Use(args...)
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
	return r.Add(fiber.MethodConnect, path, handler...)
}

func (r *router) Options(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodOptions, path, handler...)
}

func (r *router) Trace(path string, handler ...fiber.Handler) Route {
	return r.Add(fiber.MethodTrace, path, handler...)
}

func (r *router) Add(method, path string, handler ...fiber.Handler) Route {
	rr := r.router.Add(method, path, handler...)

	route := &route{
		rr:       r,
		router:   rr,
		tags:     r.tags,
		security: r.security,
		hide:     r.hide,
	}
	if r.core != nil {
		cleanpath := r.cleanPath(path)
		cleanpath = util.ConvertPath(cleanpath)

		operation, err := r.core.NewOperationContext(method, cleanpath)
		if err != nil {
			r.errors.Add(err)
		}
		route.operation = operation
	}
	r.routes = append(r.routes, route)

	return route
}

func (r *router) Static(prefix, root string, config ...fiber.Static) Router {
	r.router.Static(prefix, root, config...)
	return r
}

func (r *router) Group(prefix string, handlers ...fiber.Handler) Router {
	rr := r.router.Group(prefix, handlers...)

	subRouter := &router{
		router:   rr,
		core:     r.core,
		prefix:   r.cleanPath(prefix),
		logger:   r.logger,
		errors:   r.errors,
		tags:     r.tags,
		security: r.security,
		hide:     r.hide,
	}
	r.subRouters = append(r.subRouters, subRouter)

	return subRouter
}

func (r *router) Route(prefix string, fn func(router Router)) Router {
	fr := r.router.Group(prefix) // or Route, same

	subRouter := &router{
		router:   fr,
		core:     r.core,
		prefix:   r.cleanPath(prefix),
		logger:   r.logger,
		errors:   r.errors,
		tags:     r.tags,
		security: r.security,
		hide:     r.hide,
	}

	fn(subRouter)

	r.subRouters = append(r.subRouters, subRouter)

	return subRouter
}

func (r *router) With(opts ...option.RouteOption) Router {
	cfg := &option.RouteConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	if len(cfg.Tags) > 0 {
		r.logger.Printf("Adding tags to router: %v for path: %s", cfg.Tags, r.prefix)
		r.tags = append(r.tags, cfg.Tags...)
	}
	if len(cfg.Security) > 0 {
		r.logger.Printf("Adding security to router: %v for path: %s", cfg.Security, r.prefix)
		r.security = append(r.security, cfg.Security...)
	}
	if cfg.Hide && !r.hide {
		r.logger.Printf("Hiding router for path: %s", r.prefix)
		r.hide = true
	}
	return r
}

func (r *router) Validate() error {
	r.addOperations()
	if r.errors.HasErrors() {
		return r.errors
	}
	return nil
}

func (r *router) GenerateOpenAPISchema(formats ...string) ([]byte, error) {
	if r.core == nil {
		return nil, fmt.Errorf("OpenAPI is disabled, cannot generate schema")
	}
	r.addOperations()
	if r.errors.HasErrors() {
		return nil, r.errors
	}
	return r.core.GenerateSchema(formats...)
}

func (r *router) WriteSchemaTo(path string) error {
	if r.core == nil {
		return fmt.Errorf("OpenAPI is disabled, cannot write schema")
	}
	r.addOperations()
	if r.errors.HasErrors() {
		return r.errors
	}
	return r.core.WriteSchemaTo(path)
}

// addOperations builds the OpenAPI operations for the router and its sub-routers.
// This method is called only once to avoid redundant operations.
func (r *router) addOperations() {
	// Ensure that operations are built only once.
	r.buildOnce.Do(func() {
		for _, operation := range r.build() {
			if operation == nil {
				continue
			}

			r.logger.Printf("Adding operation for path: %s with method: %s", operation.PathPattern(), operation.Method())
			if err := r.core.AddOperation(operation); err != nil {
				r.errors.Add(err)
			}
		}
	})
}

// build constructs the OpenAPI operations for the router and its sub-routers.
func (r *router) build() []openapi.OperationContext {
	var operations []openapi.OperationContext

	for _, route := range r.routes {
		if route.operation == nil || route.hide || r.hide {
			continue
		}

		path := route.operation.Method() + " " + route.operation.PathPattern()

		// Add tags from the route to the operation
		tags := append(r.tags, route.tags...)
		if len(tags) > 0 {
			r.logger.Printf("Adding tags: %v for path: %s", tags, path)
			route.operation.SetTags(tags...)
		}

		// Add security schemes from the route to the operation
		securities := append(r.security, route.security...)
		for _, sec := range securities {
			r.logger.Printf("Adding security scheme: %s with scopes: %v for path: %s", sec.Name, sec.Scopes, path)
			route.operation.AddSecurity(sec.Name, sec.Scopes...)
		}

		operations = append(operations, route.operation)
	}

	// Recursively build operations for sub-routers
	for _, subRouter := range r.subRouters {
		operations = append(operations, subRouter.build()...)
	}

	return operations
}

func (r *router) cleanPath(paths ...string) string {
	return stdpath.Clean(stdpath.Join(append([]string{r.prefix}, paths...)...))
}
