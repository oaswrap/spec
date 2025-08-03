# ginopenapi

[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec/adapters/ginopenapi.svg)](https://pkg.go.dev/github.com/oaswrap/spec/adapters/ginopenapi)

A lightweight adapter for the [Gin](https://github.com/gin-gonic/gin) web framework that automatically generates OpenAPI 3.x specifications from your routes using [`oaswrap/spec`](https://github.com/oaswrap/spec).

## Why ginopenapi?

- **⚡ Seamless Integration** — Works with your existing Gin routes and handlers
- **📝 Automatic Documentation** — Generate OpenAPI specs from route definitions and struct tags
- **🎯 Type Safety** — Full Go type safety for OpenAPI configuration
- **🔧 Built-in UI** — Swagger UI served automatically at `/docs`
- **🚀 Zero Overhead** — Minimal performance impact on your API

## Installation

```bash
go get github.com/oaswrap/spec/adapters/ginopenapi
```

## Quick Start

```go
package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec/adapters/ginopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	e := gin.Default()

	// Create a new OpenAPI router
	r := ginopenapi.NewRouter(e,
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

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

	e.Run(":3000")
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

func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}
	// Simulate login logic
	c.JSON(200, LoginResponse{Token: "example-token"})
}

func GetUserHandler(c *gin.Context) {
	var req GetUserRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(400, map[string]string{"error": "Invalid request"})
		return
	}
	// Simulate fetching user by ID
	user := User{ID: req.ID, Name: "John Doe"}
	c.JSON(200, user)
}
```

## Configuration Options

For all available configuration options, see the main [`oaswrap/spec`](https://github.com/oaswrap/spec#configuration-options) documentation.

## API Reference

- **Core**: [pkg.go.dev/github.com/oaswrap/spec](https://pkg.go.dev/github.com/oaswrap/spec)
- **Adapter**: [pkg.go.dev/github.com/oaswrap/spec/adapters/ginopenapi](https://pkg.go.dev/github.com/oaswrap/spec/adapters/ginopenapi)
- **Options**: [pkg.go.dev/github.com/oaswrap/spec/option](https://pkg.go.dev/github.com/oaswrap/spec/option)

## Contributing

We welcome contributions! Please open issues and PRs at the main [oaswrap/spec](https://github.com/oaswrap/spec) repository.

## License

[MIT License](LICENSE) — Created with ❤️ by [Ahmad Faiz](https://github.com/afkdevs).