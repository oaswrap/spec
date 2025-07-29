// Package fiberopenapi provides an OpenAPI 3.1 documentation generator
// for Fiber applications.
//
// It wraps fiber.Router with a chainable API for defining routes, groups,
// and OpenAPI metadata alongside your handlers.
//
// # Getting Started
//
// To use fiberopenapi, create a new wrapped router:
//
//	r := fiberopenapi.NewRouter(app,
//		fiberopenapi.WithTitle("My API"),
//		fiberopenapi.WithVersion("1.0.0"),
//		fiberopenapi.WithDescription("This is a sample API"),
//		fiberopenapi.WithDocsPath("/docs"),
//		fiberopenapi.WithServer("http://localhost:3000", "Local server"),
//		fiberopenapi.WithSecurity("bearerAuth", fiberopenapi.SecurityHTTPBearer()),
//	)
//
// You can then define routes just like native Fiber:
//
//	r.Get("/ping", handler).
//		With(option.Summary("Ping"), option.Description("Ping the server"))
//
// Or use nested groups:
//
//	api := r.Group("/api")
//	v1 := api.Group("/v1")
//
//	v1.Route("/auth", func(r fiberopenapi.Router) {
//		r.Post("/login", loginHandler).With(
//			option.Summary("Login"),
//			option.Request(new(LoginRequest)),
//			option.Response(200, new(Response[Token])),
//		)
//	}, option.WithRouteTags("Authentication"))
//
// When ready, validate and export the OpenAPI spec:
//
//	if err := r.Validate(); err != nil {
//		log.Fatal(err)
//	}
//
//	r.WriteSchemaTo("openapi.yaml")
//	r.WriteSchemaTo("openapi.json")
//
// # Options
//
// The package supports ConfigOptions for router-level config,
// as well as route-level options via the `option` package.
//
// See:
//   - WithTitle
//   - WithVersion
//   - WithDescription
//   - WithDocsPath
//   - WithServer
//   - WithSecurity
//
// For detailed OpenAPI options for operations, request bodies,
// responses, and security schemes, refer to the `option` subpackage.
//
// # Example
//
// See `examples/fiber-example` for a complete runnable setup.
//
// For full documentation, visit: https://pkg.go.dev/github.com/faizlabs/openapi-wrapper/fiberopenapi
package fiberopenapi
