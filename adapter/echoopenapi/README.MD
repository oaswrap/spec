# echoopenapi

[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec/adapter/echoopenapi.svg)](https://pkg.go.dev/github.com/oaswrap/spec/adapter/echoopenapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/oaswrap/spec/adapter/echoopenapi)](https://goreportcard.com/report/github.com/oaswrap/spec/adapter/echoopenapi)

A lightweight adapter for the [Echo](https://github.com/labstack/echo) web framework that automatically generates OpenAPI 3.x specifications from your routes using [`oaswrap/spec`](https://github.com/oaswrap/spec).

> **Note:** This adapter is for Echo v4. For Echo v5, use [`echov5openapi`](../echov5openapi).

## Features

- **⚡ Seamless Integration** — Works with your existing Echo routes and handlers
- **📝 Automatic Documentation** — Generate OpenAPI specs from route definitions and struct tags
- **🎯 Type Safety** — Full Go type safety for OpenAPI configuration
- **🔧 Multiple UI Options** — Swagger UI, Stoplight Elements, ReDoc, Scalar or RapiDoc served automatically at `/docs`
- **📄 YAML Export** — OpenAPI spec available at `/docs/openapi.yaml`
- **🚀 Zero Overhead** — Minimal performance impact on your API

## Installation

```bash
go get github.com/oaswrap/spec/adapter/echoopenapi
```

## Quick Start

```go
package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/oaswrap/spec/adapter/echoopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	e := echo.New()

	// Create a new OpenAPI router
	r := echoopenapi.NewRouter(e,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
	)
	// Add routes
	v1 := r.Group("/api/v1")

	v1.POST("/login", LoginHandler).With(
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)

	auth := v1.Group("", AuthMiddleware).With(
		option.GroupSecurity("bearerAuth"),
	)
	auth.GET("/users/:id", GetUserHandler).With(
		option.Summary("Get user by ID"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
	)

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

	if err := e.Start(":3000"); err != nil {
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
	ID string `param:"id" required:"true"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Simulate authentication logic
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" && authHeader == "Bearer example-token" {
			return next(c)
		}
		return c.JSON(401, map[string]string{"error": "Unauthorized"})
	}
}

func LoginHandler(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}
	// Simulate login logic
	return c.JSON(200, LoginResponse{Token: "example-token"})
}

func GetUserHandler(c echo.Context) error {
	var req GetUserRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, map[string]string{"error": "Invalid request"})
	}
	// Simulate fetching user by ID
	user := User{ID: req.ID, Name: "John Doe"}
	return c.JSON(200, user)
}
```

## Documentation Features

### Built-in Endpoints
When you create a echoopenapi router, the following endpoints are automatically available:

- **`/docs`** — Interactive UI documentation
- **`/docs/openapi.yaml`** — Raw OpenAPI specification in YAML format

If you want to disable the built-in UI, you can do so by passing `option.WithDisableDocs()` when creating the router:

```go
r := echoopenapi.NewRouter(c,
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
r := echoopenapi.NewRouter(c,
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

For more struct tag options, see the [swaggest/openapi-go](https://github.com/swaggest/openapi-go?tab=readme-ov-file#features).

## Example

Check out the [examples directory](/adapter/echoopenapi/example) for more complete implementations and use cases.

## Best Practices

1. **Organize with Tags** — Group related operations using `option.Tags()`
2. **Document Everything** — Use `option.Summary()` and `option.Description()` for all routes
3. **Define Error Responses** — Include common error responses (400, 401, 404, 500)
4. **Use Validation Tags** — Leverage struct tags for request validation documentation
5. **Security First** — Define and apply appropriate security schemes
6. **Version Your API** — Use route groups for API versioning (`/api/v1`, `/api/v2`)

## API Reference

- **Spec**: [pkg.go.dev/github.com/oaswrap/spec](https://pkg.go.dev/github.com/oaswrap/spec)
- **Echo Adapter**: [pkg.go.dev/github.com/oaswrap/spec/adapter/echoopenapi](https://pkg.go.dev/github.com/oaswrap/spec/adapter/echoopenapi)
- **Options**: [pkg.go.dev/github.com/oaswrap/spec/option](https://pkg.go.dev/github.com/oaswrap/spec/option)
- **Spec UI**: [pkg.go.dev/github.com/oaswrap/spec-ui](https://pkg.go.dev/github.com/oaswrap/spec-ui)

## Contributing

We welcome contributions! Please open issues and PRs at the main [oaswrap/spec](https://github.com/oaswrap/spec) repository.

## License

[MIT](../../LICENSE)