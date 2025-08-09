package muxopenapi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oaswrap/spec/option"
)

// Generator is an interface that defines methods for generating OpenAPI schemas.
type Generator interface {
	Router

	// GenerateSchema generates the OpenAPI schema for the router.
	GenerateSchema(formats ...string) ([]byte, error)

	// MarshalJSON marshals the schema to JSON.
	MarshalJSON() ([]byte, error)

	// MarshalYAML marshals the schema to YAML.
	MarshalYAML() ([]byte, error)

	// Validate validates the schema.
	Validate() error

	// WriteSchemaTo writes the schema to a file.
	WriteSchemaTo(path string) error
}

// Router is an interface that defines methods for handling HTTP routes with OpenAPI support.
type Router interface {
	http.Handler

	// Get returns a route registered with the given name.
	Get(name string) Route

	// GetRoute returns a route registered with the given name.
	GetRoute(name string) Route

	// Handle registers a new route with a matcher for the URL path. See Route.Path() and Route.Handler().
	Handle(path string, handler http.Handler) Route

	// HandleFunc registers a new route with a matcher for the URL path. See Route.Path() and Route.HandlerFunc().
	HandleFunc(path string, handler func(http.ResponseWriter, *http.Request)) Route

	// Headers registers a new route with a matcher for request header values. See Route.Headers().
	Headers(pairs ...string) Route

	// Host registers a new route with a matcher for the URL host. See Route.Host().
	Host(host string) Route

	// Methods registers a new route with a matcher for HTTP methods. See Route.Methods().
	Methods(methods ...string) Route

	// Name registers a new route with a name. See Route.Name().
	Name(name string) Route

	// NewRoute registers an empty route.
	NewRoute() Route

	// Path registers a new route with a matcher for the URL path. See Route.Path().
	Path(path string) Route

	// PathPrefix registers a new route with a matcher for the URL path prefix. See Route.PathPrefix().
	PathPrefix(prefix string) Route

	// Queries registers a new route with a matcher for URL query values. See Route.Queries().
	Queries(pairs ...string) Route

	// Schemes registers a new route with a matcher for URL schemes. See Route.Schemes().
	Schemes(schemes ...string) Route

	// SkipClean defines the path cleaning behaviour for new routes. The initial value is false.
	SkipClean(value bool) Router

	// StrictSlash defines the trailing slash behavior for new routes. The initial value is false.
	StrictSlash(value bool) Router

	// Use appends a MiddlewareFunc to the chain. Middleware can be used to intercept or otherwise modify requests and/or responses, and are executed in the order that they are applied to the Router.
	Use(middlewares ...mux.MiddlewareFunc) Router

	// UseEncodedPath tells the router to match the encoded original path to the routes.
	UseEncodedPath() Router

	// With applies OpenAPI group options to this router.
	With(opts ...option.GroupOption) Router
}

// Route defines the interface for a route that can handle HTTP requests.
type Route interface {
	// GetError returns an error resulted from building the route, if any.
	GetError() error

	// GetHandler returns the handler for the route, if any.
	GetHandler() http.Handler

	// GetName returns the name for the route, if any.
	GetName() string

	// Handler sets a handler for the route.
	Handler(handler http.Handler) Route

	// HandlerFunc sets a handler function for the route.
	HandlerFunc(handler func(http.ResponseWriter, *http.Request)) Route

	// Headers adds a matcher for request header values. It accepts a sequence of key/value pairs to be matched.
	Headers(pairs ...string) Route

	// Host adds a matcher for the URL host. It accepts a template with zero or more URL variables enclosed by {}.
	Host(host string) Route

	// Methods adds a matcher for HTTP methods. It accepts a sequence of one or more methods to be matched, e.g.: "GET", "POST", "PUT".
	Methods(methods ...string) Route

	// Name sets the name for the route, used to build URLs. It is an error to call Name more than once on a route.
	Name(name string) Route

	// Path adds a matcher for the URL path. It accepts a template with zero or more URL variables enclosed by {}. The template must start with a "/".
	Path(path string) Route

	// PathPrefix adds a matcher for the URL path prefix. This matches if the given template is a prefix of the full URL path.
	PathPrefix(prefix string) Route

	// Queries adds a matcher for URL query values. It accepts a sequence of key/value pairs.
	Queries(pairs ...string) Route

	// Schemes adds a matcher for URL schemes. It accepts a sequence of schemes to be matched, e.g.: "http", "https".
	Schemes(schemes ...string) Route

	// SkipClean reports whether path cleaning is enabled for this route via Router.SkipClean.
	SkipClean() bool

	// Subrouter creates a subrouter for the route.
	Subrouter() Router

	// With applies OpenAPI operation options to this route.
	With(opts ...option.OperationOption) Route
}
