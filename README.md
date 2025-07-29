# OpenAPI Wrapper

[![Go Reference](https://pkg.go.dev/badge/github.com/faizlabs/openapi-wrapper.svg)](https://pkg.go.dev/github.com/faizlabs/openapi-wrapper)
[![Go Report Card](https://goreportcard.com/badge/github.com/faizlabs/openapi-wrapper)](https://goreportcard.com/report/github.com/faizlabs/openapi-wrapper)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

**OpenAPI Wrapper** is a flexible library that helps you automatically generate **OpenAPI 3.x** specifications for your Go web APIs.  
It works by wrapping your router — starting with **Fiber** — and adding OpenAPI metadata in an idiomatic way.

## ✨ Features

- 📚 **Automatic OpenAPI generation** — every route can add spec details inline.
- 🔗 **Chainable API** — define routes just like native `fiber.Router`, add OpenAPI docs with `.With(...)`.
- ⚙️ **Fully configurable** — set spec version (**3.0** or **3.1**), titles, servers, security schemes.
- 🔒 **Security schemes** — define and apply API key or bearer authentication.
- ⚡️ **Route groups** — groups inherit prefixes and options naturally.
- 📂 **Framework-agnostic core** — only `fiberopenapi` depends on Fiber.

## 📦 Install

```bash
go get github.com/faizlabs/openapi-wrapper
```

## 📚 Documentation

Full API documentation is available on [pkg.go.dev](https://pkg.go.dev/github.com/faizlabs/openapi-wrapper):

- [openapi-wrapper](https://pkg.go.dev/github.com/faizlabs/openapi-wrapper)
- [fiberopenapi](https://pkg.go.dev/github.com/faizlabs/openapi-wrapper/fiberopenapi)
- [option](https://pkg.go.dev/github.com/faizlabs/openapi-wrapper/option)

## 🚀 Usage with Fiber

Below is a complete example using **Fiber** with **OpenAPI Wrapper**:

```go
package main

import (
	"log"

	"github.com/faizlabs/openapi-wrapper/fiberopenapi"
	"github.com/faizlabs/openapi-wrapper/option"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	// Initialize OpenAPI router with config.
	// Use WithOpenAPIVersion to choose 3.0 or 3.1
	r := fiberopenapi.NewRouter(app,
		option.WithOpenAPIVersion("3.1"), // Or "3.0"
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithDescription("This is a sample API"),
		option.WithServer("http://localhost:3000", "Local server"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer()),
	)

	// Group routes
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Example route with chained metadata
	v1.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	}).With(
		option.Summary("Ping"),
		option.Description("Check server health"),
	)

	// Example grouped route with inline Route
	v1.Route("/auth", func(r fiberopenapi.Router) {
		r.Post("/login", dummyHandler).With(
			option.Summary("User Login"),
			option.Description("Authenticate user and return token"),
			option.Request(new(LoginRequest)),
			option.Response(200, new(Response[Token])),
		)
	}).With(option.RouteTags("Authentication"))

	// Validate spec & write files (optional)
	if err := r.Validate(); err != nil {
		log.Fatal(err)
	}
	r.WriteSchemaTo("openapi.yaml")
	r.WriteSchemaTo("openapi.json")

	app.Listen(":3000")
}

func dummyHandler(c *fiber.Ctx) error {
	return c.SendString("dummy")
}

type LoginRequest struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"password123"`
}

type Token struct {
	AccessToken string `json:"access_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

type Response[T any] struct {
	Status int   `json:"status" example:"200"`
	Data   T     `json:"data"`
}
```

## 🏷️ Struct Tags

This package uses [openapi-go](https://github.com/swaggest/openapi-go) for schema generation. Supported tags include:

- `json`: JSON field name.
- `query`: Query parameter name.
- `path`: Path parameter name.
- `header`: Header parameter name.
- `formData`: Form data parameter name.
- `title`: Custom title for the field.
- `description`: Description for the field.

For more details, see the [openapi-go](https://github.com/swaggest/openapi-go?tab=readme-ov-file#features) and [jsonschema-go](https://github.com/swaggest/jsonschema-go#field-tags) documentation.

## ✅ Highlights

- No extra Swagger annotations — just idiomatic Go.
- `.Get()`, `.Post()` accept multiple handlers, just like native Fiber.
- `.With(...)` adds OpenAPI metadata to the last defined route.
- Route groups allow scoped options like tags or security.
- `WriteSchemaTo` writes JSON/YAML spec files.
- `Validate` checks your OpenAPI spec at runtime.
- Supports **OpenAPI 3.0** and **3.1** — configurable per project.

## 📄 License

This project is licensed under the MIT License — see the [LICENSE](LICENSE) file for details.