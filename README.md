# oaswrap/spec

[![CI](https://github.com/oaswrap/spec/actions/workflows/ci.yml/badge.svg)](https://github.com/oaswrap/spec/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/oaswrap/spec/branch/main/graph/badge.svg)](https://codecov.io/gh/oaswrap/spec)
[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec.svg)](https://pkg.go.dev/github.com/oaswrap/spec)
[![Go Report Card](https://goreportcard.com/badge/github.com/oaswrap/spec)](https://goreportcard.com/report/github.com/oaswrap/spec)
[![License](https://img.shields.io/github/license/oaswrap/spec)](LICENSE)

**`oaswrap/spec`** lets you build OpenAPI 3.x specs in pure Go — framework-agnostic and easy to integrate.

Describe your API operations, paths, and schemas once, then plug them into any router.  
No handlers, no routing — just pure OpenAPI generation.

Powered by [`swaggest/openapi-go`](https://github.com/swaggest/openapi-go) for robust schema generation.

## ✨ Features

- ✅ Programmatically build OpenAPI 3.x specs in pure Go.
- ✅ No runtime server — only schema generation logic.
- ✅ Designed for framework adapters (Gin, Echo, Fiber, etc.).
- ✅ Supports struct tags for request/response models.
- ✅ Export specs to JSON or YAML, and validate before publishing.

## 🔗 Related Projects

Need an integration? Check out these official adapters:
- [`oaswrap/ginopenapi`](https://github.com/oaswrap/ginopenapi) — Gin integration
- [`oaswrap/echoopenapi`](https://github.com/oaswrap/echoopenapi) — Echo integration
- [`oaswrap/fiberopenapi`](https://github.com/oaswrap/fiberopenapi) — Fiber integration

## 📦 Installation

```bash
go get github.com/oaswrap/spec
```

## 🚀 Usage Example

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

✨ **Live example:** View the generated spec on [Rest.Wiki](https://rest.wiki/?https://raw.githubusercontent.com/oaswrap/spec/main/examples/basic/openapi.yaml).

## 📚 Documentation

- All core configuration, router, server, and security options are defined in the [`option`](https://pkg.go.dev/github.com/oaswrap/spec/option) package.
- See the [full API reference on pkg.go.dev](https://pkg.go.dev/github.com/oaswrap/spec) for detailed usage, examples, and type definitions.
- This library uses [`swaggest/openapi-go`](https://github.com/swaggest/openapi-go) under the hood — see its docs for advanced struct tagging and schema reflection tips.


## 📄 License

This project is licensed under the [MIT License](LICENSE).

## 🤝 Contributing

PRs and issues are welcome! ❤️  
Made with care by [Ahmad Faiz](https://github.com/afkdevs)