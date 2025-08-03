# fiberopenapi

[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec/adapters/fiberopenapi.svg)](https://pkg.go.dev/github.com/oaswrap/spec/adapters/fiberopenapi)

A lightweight adapter for the [Fiber](https://github.com/gofiber/fiber) web framework that automatically generates OpenAPI 3.x specifications from your routes using [`oaswrap/spec`](https://github.com/oaswrap/spec).

## Why fiberopenapi?

- **⚡ Seamless Integration** — Works with your existing Fiber routes and handlers
- **📝 Automatic Documentation** — Generate OpenAPI specs from route definitions and struct tags
- **🎯 Type Safety** — Full Go type safety for OpenAPI configuration
- **🔧 Built-in UI** — Swagger UI served automatically at `/docs`
- **🚀 Zero Overhead** — Minimal performance impact on your API

## Installation

```bash
go get github.com/oaswrap/spec/adapters/fiberopenapi
```

## Quick Start

```go
package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/spec/adapters/fiberopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	app := fiber.New()

	// Create a new OpenAPI router
	r := fiberopenapi.NewRouter(app,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
	)
	// Add routes
	v1 := r.Group("/api/v1")

	v1.Post("/login", LoginHandler).With(
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)

	v1.Get("/users/{id}", GetUserHandler).With(
		option.Summary("Get user by ID"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
	)

	// Generate OpenAPI spec
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ OpenAPI spec generated at openapi.yaml")
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

func LoginHandler(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(map[string]string{"error": "Invalid request"})
	}
	// Simulate login logic
	return c.Status(200).JSON(LoginResponse{Token: "example-token"})
}

func GetUserHandler(c *fiber.Ctx) error {
	var req GetUserRequest
	if err := c.ParamsParser(&req); err != nil {
		return c.Status(400).JSON(map[string]string{"error": "Invalid request"})
	}
	// Simulate fetching user by ID
	user := User{ID: req.ID, Name: "John Doe"}
	return c.Status(200).JSON(user)
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
- **Adapter**: [pkg.go.dev/github.com/oaswrap/spec/adapters/fiberopenapi](https://pkg.go.dev/github.com/oaswrap/spec/adapters/fiberopenapi)
- **Options**: [pkg.go.dev/github.com/oaswrap/spec/option](https://pkg.go.dev/github.com/oaswrap/spec/option)

## Contributing

We welcome contributions! Please open issues and PRs at the main [oaswrap/spec](https://github.com/oaswrap/spec) repository.

## License

[MIT License](LICENSE) — Created with ❤️ by [Ahmad Faiz](https://github.com/afkdevs).