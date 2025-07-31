# oaswrap/spec

[![CI](https://github.com/oaswrap/spec/actions/workflows/ci.yml/badge.svg)](https://github.com/oaswrap/spec/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/oaswrap/spec/branch/main/graph/badge.svg)](https://codecov.io/gh/oaswrap/spec)
[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec.svg)](https://pkg.go.dev/github.com/oaswrap/spec)
[![Go Report Card](https://goreportcard.com/badge/github.com/oaswrap/spec)](https://goreportcard.com/report/github.com/oaswrap/spec)
[![License](https://img.shields.io/github/license/oaswrap/spec)](LICENSE)

**`oaswrap/spec`** is a lightweight, framework-agnostic OpenAPI 3.x specification builder for Go.  
It provides the core logic to describe your API operations, paths, parameters, and schemas — without locking you into any specific web framework.

This makes it easy to use with any router — directly or through adapters for frameworks like Fiber, Gin, Echo, and more.

Under the hood, `oaswrap/spec` uses [`swaggest/openapi-go`](https://github.com/swaggest/openapi-go) for robust OpenAPI schema generation.

> ✅ Looking for a ready-to-use Gin integration? Check out [`oaswrap/ginopenapi`](https://github.com/oaswrap/ginopenapi).

> ✅ Looking for a ready-to-use Fiber integration? Check out [`oaswrap/fiberopenapi`](https://github.com/oaswrap/fiberopenapi).

## ✨ Features

- ✅ Programmatically build OpenAPI 3.x specs in pure Go.
- ✅ Powered by [`swaggest/openapi-go`](https://github.com/swaggest/openapi-go).
- ✅ No runtime web server logic — focused purely on schema generation.
- ✅ Designed to be wrapped by adapters for popular frameworks.
- ✅ Supports struct tags for request/response models.
- ✅ Write specs to JSON or YAML, validate before serving or publishing.

## 📦 Installation

```bash
go get github.com/oaswrap/spec
```

## ⚡️ Quick Example

```go
package main

import (
	"log"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

func main() {
	// Create a new OpenAPI router with basic info and security scheme
	r := spec.NewRouter(
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithDescription("Example API"),
		option.WithServer("https://api.example.com"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
	)

	// Versioned API group
	v1 := r.Group("/api/v1")

	// Auth routes
	v1.Route("/auth", func(r spec.Router) {
		r.Post("/login",
			option.Summary("User Login"),
			option.Request(new(LoginRequest)),
			option.Response(200, new(Response[Token])),
		)

		r.Get("/me",
			option.Summary("Get Profile"),
			option.Security("bearerAuth"),
			option.Response(200, new(Response[User])),
		)
	}, option.GroupTags("Authentication"))

	// Generate the OpenAPI file
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}

	log.Println("✅ OpenAPI schema generated at openapi.yaml")
}

// Example request & response structs

type LoginRequest struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

type Token struct {
	AccessToken string `json:"access_token"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Response[T any] struct {
	Status int `json:"status"`
	Data   T   `json:"data"`
}
```
✨ See it live: You can view the generated OpenAPI documentation for this example online at [Rest.Wiki Viewer](https://rest.wiki/?https://raw.githubusercontent.com/oaswrap/spec/main/examples/basic/openapi.yaml).

## 📚 Documentation

For detailed usage instructions, see the [pkg.go.dev documentation](https://pkg.go.dev/github.com/oaswrap/spec).

## 📄 License

This project is licensed under the [MIT License](LICENSE).

**Made with ❤️ by [oaswrap](https://github.com/oaswrap)**