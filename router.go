package spec

import "github.com/oaswrap/spec/option"

// Router is an interface that defines methods for registering routes and operations
// in an OpenAPI specification. It allows for defining HTTP methods, paths, and operation options.
type Router interface {
	// Get registers a new GET operation with the specified path and options.
	Get(path string, opts ...option.OperationOption)
	// Post registers a new POST operation with the specified path and options.
	Post(path string, opts ...option.OperationOption)
	// Put registers a new PUT operation with the specified path and options.
	Put(path string, opts ...option.OperationOption)
	// Delete registers a new DELETE operation with the specified path and options.
	Delete(path string, opts ...option.OperationOption)
	// Patch registers a new PATCH operation with the specified path and options.
	Patch(path string, opts ...option.OperationOption)
	// Options registers a new OPTIONS operation with the specified path and options.
	Options(path string, opts ...option.OperationOption)
	// Head registers a new HEAD operation with the specified path and options.
	Head(path string, opts ...option.OperationOption)
	// Add registers a new operation with the specified method and path.
	Add(method, path string, opts ...option.OperationOption)
	// Trace registers a new TRACE operation with the specified path and options.
	Trace(path string, opts ...option.OperationOption)

	// Route registers a new route with the specified pattern and function.
	// The function receives a Router instance to define sub-routes.
	Route(pattern string, fn func(router Router), opts ...option.RouteOption) Router
	// Group creates a new sub-router with the specified prefix and options.
	Group(pattern string, opts ...option.RouteOption) Router
	// Use applies the provided options to the router.
	Use(opts ...option.RouteOption) Router
}
