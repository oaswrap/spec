# echoopenapi

[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec/adapters/echoopenapi.svg)](https://pkg.go.dev/github.com/oaswrap/spec/adapters/echoopenapi)

A lightweight adapter for the [Echo](https://github.com/labstack/echo) web framework that automatically generates OpenAPI 3.x specifications from your routes using [`oaswrap/spec`](https://github.com/oaswrap/spec).

## Why echoopenapi?

- **⚡ Seamless Integration** — Works with your existing Echo routes and handlers
- **📝 Automatic Documentation** — Generate OpenAPI specs from route definitions and struct tags
- **🎯 Type Safety** — Full Go type safety for OpenAPI configuration
- **🔧 Built-in UI** — Swagger UI served automatically at `/docs`
- **🚀 Zero Overhead** — Minimal performance impact on your API

## Installation

```bash
go get github.com/oaswrap/spec/adapters/echoopenapi
```

## Quick Start

```go
package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/oaswrap/spec/adapters/echoopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	e := echo.New()

	// Create a new OpenAPI router
	r := echoopenapi.NewRouter(e,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
	)
	// Add routes
	v1 := r.Group("/api/v1")

	v1.POST("/login", LoginHandler).With(
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)

	v1.GET("/users/{id}", GetUserHandler).With(
		option.Summary("Get user by ID"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
	)

	// Generate OpenAPI spec
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}
	log.Println("✅ OpenAPI spec generated at openapi.yaml")

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

	e.Start(":3000")
}

type LoginRequest struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type GetUserRequest struct {
	ID string `path:"id" required:"true"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
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

## Advanced Features

### Route Groups with Common Settings
```go
// Apply settings to all routes in a group
adminAPI := api.Group("/admin").With(
    option.GroupTags("Administration"),
	option.GroupSecurity("bearerAuth"),
)

adminAPI.GET("/users", getUsersHandler).With(
	option.Summary("List all users"),
	option.Response(200, new([]User)),
)
```

## Configuration Options

For all available configuration options, see the main [`oaswrap/spec`](https://github.com/oaswrap/spec#configuration-options) documentation.

## API Reference

- **Core**: [pkg.go.dev/github.com/oaswrap/spec](https://pkg.go.dev/github.com/oaswrap/spec)
- **Adapter**: [pkg.go.dev/github.com/oaswrap/spec/adapters/echoopenapi](https://pkg.go.dev/github.com/oaswrap/spec/adapters/echoopenapi)
- **Options**: [pkg.go.dev/github.com/oaswrap/spec/option](https://pkg.go.dev/github.com/oaswrap/spec/option)

## Contributing

We welcome contributions! Please open issues and PRs at the main [oaswrap/spec](https://github.com/oaswrap/spec) repository.

## License

[MIT License](LICENSE) — Created with ❤️ by [Ahmad Faiz](https://github.com/afkdevs).