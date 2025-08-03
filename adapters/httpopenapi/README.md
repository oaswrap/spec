# httpopenapi

[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec/adapters/httpopenapi.svg)](https://pkg.go.dev/github.com/oaswrap/spec/adapters/httpopenapi)

A lightweight adapter for the [HTTP](https://golang.org/pkg/net/http/) web framework that automatically generates OpenAPI 3.x specifications from your routes using [`oaswrap/spec`](https://github.com/oaswrap/spec).

## Why httpopenapi?

- **⚡ Seamless Integration** — Works with your existing HTTP routes and handlers
- **📝 Automatic Documentation** — Generate OpenAPI specs from route definitions and struct tags
- **🎯 Type Safety** — Full Go type safety for OpenAPI configuration
- **🔧 Built-in UI** — Swagger UI served automatically at `/docs`
- **🚀 Zero Overhead** — Minimal performance impact on your API

## Installation

```bash
go get github.com/oaswrap/spec/adapters/httpopenapi
```

## Quick Start

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/oaswrap/spec/adapters/httpopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	mainMux := http.NewServeMux()
	r := httpopenapi.NewGenerator(mainMux,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
	)

	r.Route("/api/v1", func(r httpopenapi.Router) {
		r.HandleFunc("POST /login", LoginHandler).With(
			option.Summary("User login"),
			option.Request(new(LoginRequest)),
			option.Response(200, new(LoginResponse)),
		)
		auth := r.Group("/", AuthMiddleware).With(
			option.GroupSecurity("bearerAuth"),
		)
		auth.HandleFunc("GET /users/{id}", GetUserHandler).With(
			option.Summary("Get user by ID"),
			option.Request(new(GetUserRequest)),
			option.Response(200, new(User)),
		)
	})

	// Generate OpenAPI spec
	if err := r.WriteSchemaTo("openapi.yaml"); err != nil {
		log.Fatal(err)
	}

	log.Printf("🚀 OpenAPI docs available at: %s", "http://localhost:3000/docs")

	// Start the server
	if err := http.ListenAndServe(":3000", mainMux); err != nil {
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
	ID string `path:"id" required:"true"`
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate authentication logic
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && authHeader == "Bearer example-token" {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	// Simulate login logic
	_ = json.NewEncoder(w).Encode(LoginResponse{Token: "example-token"})
}

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
	var req GetUserRequest
	id := r.PathValue("id")
	if id == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}
	req.ID = id
	// Simulate fetching user by ID
	user := User{ID: req.ID, Name: "John Doe"}
	_ = json.NewEncoder(w).Encode(user)
}
```

## Configuration Options

For all available configuration options, see the main [`oaswrap/spec`](https://github.com/oaswrap/spec#configuration-options) documentation.

## API Reference

- **Core**: [pkg.go.dev/github.com/oaswrap/spec](https://pkg.go.dev/github.com/oaswrap/spec)
- **Adapter**: [pkg.go.dev/github.com/oaswrap/spec/adapters/httpopenapi](https://pkg.go.dev/github.com/oaswrap/spec/adapters/httpopenapi)
- **Options**: [pkg.go.dev/github.com/oaswrap/spec/option](https://pkg.go.dev/github.com/oaswrap/spec/option)

## Contributing

We welcome contributions! Please open issues and PRs at the main [oaswrap/spec](https://github.com/oaswrap/spec) repository.

## License

[MIT License](LICENSE) — Created with ❤️ by [Ahmad Faiz](https://github.com/afkdevs).