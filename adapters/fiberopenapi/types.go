package fiberopenapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/spec/option"
)

// Generator defines the interface for generating OpenAPI schemas.
type Generator interface {
	Router

	// Validate checks for errors at OpenAPI router initialization.
	Validate() error

	// GenerateOpenAPISchema generates the OpenAPI schema in the specified format.
	GenerateOpenAPISchema(format ...string) ([]byte, error)
	// MarshalYAML marshals the OpenAPI schema to YAML format.
	MarshalYAML() ([]byte, error)
	// MarshalJSON marshals the OpenAPI schema to JSON format.
	MarshalJSON() ([]byte, error)

	// WriteSchemaTo writes the OpenAPI schema to a file.
	WriteSchemaTo(filePath string) error
}

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
	With(opts ...option.GroupOption) Router
}
