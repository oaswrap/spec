package httprouteropenapi

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/oaswrap/spec"
	specui "github.com/oaswrap/spec-ui"
	"github.com/oaswrap/spec/adapter/httprouteropenapi/internal/constant"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/mapper"
	"github.com/oaswrap/spec/pkg/parser"
	"github.com/oaswrap/spec/pkg/util"
)

// NewRouter creates a new router with the given HTTP router and options.
func NewRouter(httpRouter *httprouter.Router, opts ...option.OpenAPIOption) Generator {
	return NewGenerator(httpRouter, opts...)
}

func NewGenerator(httpRouter *httprouter.Router, opts ...option.OpenAPIOption) Generator {
	defaultOpts := []option.OpenAPIOption{
		option.WithTitle(constant.DefaultTitle),
		option.WithDescription(constant.DefaultDescription),
		option.WithVersion(constant.DefaultVersion),
		option.WithStoplightElements(),
		option.WithCacheAge(0),
		option.WithPathParser(parser.NewColonParamParser()),
	}
	opts = append(defaultOpts, opts...)
	gen := spec.NewRouter(opts...)

	r := &router{
		router:     httpRouter,
		specRouter: gen,
		gen:        gen,
	}

	cfg := gen.Config()
	if cfg.DisableDocs {
		return r
	}

	handler := specui.NewHandler(mapper.SpecUIOpts(gen)...)

	httpRouter.Handler(http.MethodGet, cfg.DocsPath, handler.Docs())
	httpRouter.Handler(http.MethodGet, cfg.SpecPath, handler.Spec())

	return r
}

type router struct {
	router      *httprouter.Router
	prefix      string
	middlewares []func(http.Handler) http.Handler

	specRouter spec.Router
	gen        spec.Generator
}

func (sr *router) wrapHandler(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h(w, r, ps)
		})
		handler := http.Handler(handlerFunc)
		for i := len(sr.middlewares) - 1; i >= 0; i-- {
			m := sr.middlewares[i]
			handler = m(handler)
		}
		handler.ServeHTTP(w, r)
	}
}

func (r *router) pathOf(path string) string {
	if r.prefix == "" {
		return path
	}
	return util.JoinPath(r.prefix, path)
}

func (r *router) Handle(method, path string, handle httprouter.Handle) Route {
	if len(r.middlewares) > 0 {
		handle = r.wrapHandler(handle)
	}
	path = r.pathOf(path)
	r.router.Handle(method, path, handle)
	rr := &route{}
	if method != http.MethodConnect {
		rr.specRoute = r.specRouter.Add(method, path)
	}

	return rr
}

func (r *router) Handler(method, path string, handler http.Handler) Route {
	if len(r.middlewares) > 0 {
		for i := len(r.middlewares) - 1; i >= 0; i-- {
			handler = r.middlewares[i](handler)
		}
	}
	r.router.Handler(method, r.pathOf(path), handler)
	rr := &route{}
	if method != http.MethodConnect {
		rr.specRoute = r.specRouter.Add(method, path)
	}

	return rr
}

func (r *router) HandlerFunc(method, path string, handlerFunc http.HandlerFunc) Route {
	return r.Handler(method, path, http.HandlerFunc(handlerFunc))
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *router) GET(path string, handle httprouter.Handle) Route {
	return r.Handle(http.MethodGet, path, handle)
}

func (r *router) POST(path string, handle httprouter.Handle) Route {
	return r.Handle(http.MethodPost, path, handle)
}

func (r *router) PUT(path string, handle httprouter.Handle) Route {
	return r.Handle(http.MethodPut, path, handle)
}

func (r *router) DELETE(path string, handle httprouter.Handle) Route {
	return r.Handle(http.MethodDelete, path, handle)
}

func (r *router) PATCH(path string, handle httprouter.Handle) Route {
	return r.Handle(http.MethodPatch, path, handle)
}

func (r *router) HEAD(path string, handle httprouter.Handle) Route {
	return r.Handle(http.MethodHead, path, handle)
}

func (r *router) OPTIONS(path string, handle httprouter.Handle) Route {
	return r.Handle(http.MethodOptions, path, handle)
}

func (r *router) Lookup(method, path string) (httprouter.Handle, httprouter.Params, bool) {
	return r.router.Lookup(method, path)
}

func (r *router) ServeFiles(path string, root http.FileSystem) {
	r.router.ServeFiles(path, root)
}

func (r *router) Group(prefix string, middlewares ...func(http.Handler) http.Handler) Router {
	group := &router{
		router:      r.router,
		middlewares: append(r.middlewares, middlewares...),
		specRouter:  r.specRouter.Group(prefix),
		prefix:      r.pathOf(prefix),
	}
	return group
}

func (r *router) With(opts ...option.GroupOption) Router {
	r.specRouter.With(opts...)
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
