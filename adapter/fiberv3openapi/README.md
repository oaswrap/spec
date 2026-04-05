# fiberv3openapi

[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec/adapter/fiberv3openapi.svg)](https://pkg.go.dev/github.com/oaswrap/spec/adapter/fiberv3openapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/oaswrap/spec/adapter/fiberv3openapi)](https://goreportcard.com/report/github.com/oaswrap/spec/adapter/fiberv3openapi)

A lightweight adapter for [Fiber v3](https://github.com/gofiber/fiber) that automatically generates OpenAPI 3.x specifications from your routes using [`oaswrap/spec`](https://github.com/oaswrap/spec).

> **Note:** This adapter is for Fiber v3. For Fiber v2, use [`fiberopenapi`](../fiberopenapi).

## Features

- **⚡ Seamless Integration** — Works with your existing Fiber v3 routes and handlers
- **📝 Automatic Documentation** — Generate OpenAPI specs from route definitions and struct tags
- **🎯 Type Safety** — Full Go type safety for OpenAPI configuration
- **🔧 Multiple UI Options** — Swagger UI, Stoplight Elements, ReDoc, Scalar or RapiDoc served automatically at `/docs`
- **📄 YAML Export** — OpenAPI spec available at `/docs/openapi.yaml`
- **🚀 Zero Overhead** — Minimal performance impact on your API

## Installation

```bash
go get github.com/oaswrap/spec/adapter/fiberv3openapi
```

## Quick Start

```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/oaswrap/spec/adapter/fiberv3openapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	app := fiber.New()

	// Create a new OpenAPI router
	r := fiberv3openapi.NewRouter(app,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
	)
	// Add routes
	v1 := r.Group("/api/v1")
	v1.Post("/login", LoginHandler).With(
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)
	auth := v1.Group("/", AuthMiddleware).With(
		option.GroupSecurity("bearerAuth"),
	)
	auth.Get("/users/:id", GetUserHandler).With(
		option.Summary("Get user by ID"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
	)

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

	if err := app.Listen(":3000"); err != nil {
		log.Fatal(err)
	}
}

type LoginRequest struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type GetUserRequest struct {
	ID string `uri:"id" required:"true"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func AuthMiddleware(c fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader != "" && authHeader == "Bearer example-token" {
		return c.Next()
	}
	return c.Status(401).JSON(map[string]string{"error": "Unauthorized"})
}

func LoginHandler(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return c.Status(400).JSON(map[string]string{"error": "Invalid request"})
	}
	return c.Status(200).JSON(LoginResponse{Token: "example-token"})
}

func GetUserHandler(c fiber.Ctx) error {
	var req GetUserRequest
	if err := c.Bind().URI(&req); err != nil {
		return c.Status(400).JSON(map[string]string{"error": "Invalid request"})
	}
	user := User{ID: req.ID, Name: "John Doe"}
	return c.Status(200).JSON(user)
}
```

## Fiber v3 Breaking Changes

If you're migrating from the `fiberopenapi` (v2) adapter, note these key differences:

| | Fiber v2 | Fiber v3 |
|---|---|---|
| **Handler signature** | `func(c *fiber.Ctx) error` | `func(c fiber.Ctx) error` |
| **Path param struct tag** | `` `params:"id"` `` | `` `uri:"id"` `` |
| **Body parsing** | `c.BodyParser(&req)` | `c.Bind().Body(&req)` |
| **Params parsing** | `c.ParamsParser(&req)` | `c.Bind().URI(&req)` |
| **Static files** | `r.Static(prefix, root)` | Use `static` middleware directly on Fiber app |

## Documentation Features

### Built-in Endpoints
When you create a fiberv3openapi router, the following endpoints are automatically available:

- **`/docs`** — Interactive UI documentation
- **`/docs/openapi.yaml`** — Raw OpenAPI specification in YAML format

To disable the built-in UI:

```go
r := fiberv3openapi.NewRouter(app,
	option.WithTitle("My API"),
	option.WithVersion("1.0.0"),
	option.WithDisableDocs(),
)
```

### Supported Documentation UIs
Choose from multiple UI options, powered by [`oaswrap/spec-ui`](https://github.com/oaswrap/spec-ui):

- **Stoplight Elements** — Modern, clean design (default)
- **Swagger UI** — Classic interface with try-it functionality
- **ReDoc** — Three-panel responsive layout
- **Scalar** — Beautiful and fast interface
- **RapiDoc** — Highly customizable

```go
r := fiberv3openapi.NewRouter(app,
	option.WithTitle("My API"),
	option.WithVersion("1.0.0"),
	option.WithScalar(), // Use Scalar as the documentation UI
)
```

### Rich Schema Documentation
Use struct tags to generate detailed OpenAPI schemas. **Note: These tags are used only for OpenAPI spec generation and documentation - they do not perform actual request validation.**

```go
type CreateProductRequest struct {
	Name        string   `json:"name" required:"true" minLength:"1" maxLength:"100"`
	Description string   `json:"description" maxLength:"500"`
	Price       float64  `json:"price" required:"true" minimum:"0" maximum:"999999.99"`
	Category    string   `json:"category" required:"true" enum:"electronics,books,clothing"`
	Tags        []string `json:"tags" maxItems:"10"`
	InStock     bool     `json:"in_stock" default:"true"`
}
```

For more struct tag options, see [swaggest/openapi-go](https://github.com/swaggest/openapi-go?tab=readme-ov-file#features).

## Example

Check out the [examples directory](/adapter/fiberv3openapi/example) for more complete implementations and use cases.

## Best Practices

1. **Organize with Tags** — Group related operations using `option.Tags()`
2. **Document Everything** — Use `option.Summary()` and `option.Description()` for all routes
3. **Define Error Responses** — Include common error responses (400, 401, 404, 500)
4. **Use Validation Tags** — Leverage struct tags for request validation documentation
5. **Security First** — Define and apply appropriate security schemes
6. **Version Your API** — Use route groups for API versioning (`/api/v1`, `/api/v2`)

## API Reference

- **Spec**: [pkg.go.dev/github.com/oaswrap/spec](https://pkg.go.dev/github.com/oaswrap/spec)
- **Fiber v3 Adapter**: [pkg.go.dev/github.com/oaswrap/spec/adapter/fiberv3openapi](https://pkg.go.dev/github.com/oaswrap/spec/adapter/fiberv3openapi)
- **Options**: [pkg.go.dev/github.com/oaswrap/spec/option](https://pkg.go.dev/github.com/oaswrap/spec/option)
- **Spec UI**: [pkg.go.dev/github.com/oaswrap/spec-ui](https://pkg.go.dev/github.com/oaswrap/spec-ui)

## Contributing

We welcome contributions! Please open issues and PRs at the main [oaswrap/spec](https://github.com/oaswrap/spec) repository.

## License

[MIT](../../LICENSE)
