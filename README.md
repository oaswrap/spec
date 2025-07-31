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
	r := spec.NewGenerator(
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
	)

	r.Post("/login",
		option.Summary("User Login"),
		option.Description("Logs in a user and returns a token"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(TokenResponse)),
	)

	if err := r.Validate(); err != nil {
		log.Fatal(err)
	}

	_ = r.WriteSchemaTo("openapi.yaml")
}

type LoginRequest struct {
	Username string `json:"username" required:"true"`
	Password string `json:"password" required:"true"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}
```

## 📚 Documentation

For detailed usage instructions, see the [pkg.go.dev documentation](https://pkg.go.dev/github.com/oaswrap/spec).

## 📄 License

This project is licensed under the [MIT License](LICENSE).

**Made with ❤️ by [oaswrap](https://github.com/oaswrap)**