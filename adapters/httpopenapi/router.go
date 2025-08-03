package httpopenapi

import (
	"net/http"
	"slices"
	"strings"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapters/httpopenapi/internal/constant"
	"github.com/oaswrap/spec/adapters/httpopenapi/internal/handler"
	"github.com/oaswrap/spec/adapters/httpopenapi/internal/parser"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/util"
)

type router struct {
	mux        *http.ServeMux
	specRouter spec.Router
	gen        spec.Generator
}

var _ Router = (*router)(nil)

func NewRouter(mux *http.ServeMux, opts ...option.OpenAPIOption) Generator {
	return NewGenerator(mux, opts...)
}

func NewGenerator(mux *http.ServeMux, opts ...option.OpenAPIOption) Generator {
	defaultOpts := []option.OpenAPIOption{
		option.WithTitle(constant.DefaultTitle),
		option.WithDescription(constant.DefaultDescription),
		option.WithVersion(constant.DefaultVersion),
		option.WithDocsPath(constant.DefaultDocsPath),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
	}
	opts = append(defaultOpts, opts...)
	gen := spec.NewRouter(opts...)

	r := &router{
		mux:        mux,
		specRouter: gen,
		gen:        gen,
	}

	cfg := gen.Config()
	if cfg.DisableDocs {
		return r
	}

	handler := handler.NewOpenAPIHandler(cfg, gen)
	openapiPath := util.JoinPath(cfg.DocsPath, constant.OpenAPIFileName)
	mux.HandleFunc(http.MethodGet+" "+cfg.DocsPath, handler.Docs)
	mux.HandleFunc(http.MethodGet+" "+openapiPath, handler.OpenAPIYaml)

	return r
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *router) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) Route {
	if handler != nil {
		r.mux.HandleFunc(pattern, handler)
	}

	route := &route{}
	routePattern, err := parser.ParseRoutePattern(pattern)
	if err != nil || routePattern.Method == "" {
		return route
	}
	route.specRoute = r.specRouter.Add(routePattern.Method, routePattern.Path)

	return route
}

func (r *router) Handle(pattern string, handler http.Handler) Route {
	if handler != nil {
		r.mux.Handle(pattern, handler)
	}

	route := &route{}
	routePattern, err := parser.ParseRoutePattern(pattern)
	if err != nil || routePattern.Method == "" {
		return route
	}
	route.specRoute = r.specRouter.Add(routePattern.Method, routePattern.Path)

	return route
}

func (r *router) Group(prefix string, mux *http.ServeMux, middlewares ...func(http.Handler) http.Handler) Router {
	// Normalize prefix
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	prefix = strings.TrimSuffix(prefix, "/")

	if mux != nil && r.mux != nil {
		// Mount with proper pattern
		pattern := prefix + "/"
		if prefix == "/" {
			pattern = "/"
		}

		handler := http.Handler(mux)
		if len(middlewares) > 0 {
			// Apply middlewares in reverse order
			slices.Reverse(middlewares)

			for _, mw := range middlewares {
				handler = mw(handler)
			}
			handler = http.StripPrefix(prefix, handler)
		} else {
			handler = http.StripPrefix(prefix, mux)
		}

		r.mux.Handle(pattern, handler)
	}

	subRouter := r.specRouter.Group(prefix)

	return &router{
		mux:        mux,
		specRouter: subRouter,
		gen:        r.gen,
	}
}

func (r *router) With(opts ...option.GroupOption) Router {
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
