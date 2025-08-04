# httpopenapi

[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec/adapters/httpopenapi.svg)](https://pkg.go.dev/github.com/oaswrap/spec/adapters/httpopenapi)

A lightweight adapter for the [net/http](https://pkg.go.dev/net/http) package that automatically generates OpenAPI 3.x specifications from your routes using [`oaswrap/spec`](https://github.com/oaswrap/spec).

## Features

- **⚡ Seamless Integration** — Works with your existing net/http routes and handlers
- **📝 Automatic Documentation** — Generate OpenAPI specs from route definitions and struct tags
- **🎯 Type Safety** — Full Go type safety for OpenAPI configuration
- **🔧 Built-in UI** — Swagger UI served automatically at `/docs`
- **📄 YAML Export** — OpenAPI spec available at `/docs/openapi.yaml`
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

## Documentation Features

### Built-in Endpoints
When you create a httpopenapi router, the following endpoints are automatically available:

- **`/docs`** — Interactive Swagger UI documentation
- **`/docs/openapi.yaml`** — Raw OpenAPI specification in YAML format

If you want to disable the built-in UI, you can do so by passing `option.WithDisableDocs()` when creating the router:

```go
r := httpopenapi.NewRouter(c,
	option.WithTitle("My API"),
	option.WithVersion("1.0.0"),
	option.WithDisableDocs(),
)
```

### Rich Schema Documentation
Use struct tags to generate detailed OpenAPI schemas:

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

## Examples

Check out complete examples in the main repository:
- [Basic HTTP Example](https://github.com/oaswrap/spec/tree/main/examples/adapters/httpopenapi/basic)

## Best Practices

1. **Organize with Tags** — Group related operations using `option.Tags()`
2. **Document Everything** — Use `option.Summary()` and `option.Description()` for all routes
3. **Define Error Responses** — Include common error responses (400, 401, 404, 500)
4. **Use Validation Tags** — Leverage struct tags for request validation documentation
5. **Security First** — Define and apply appropriate security schemes
6. **Version Your API** — Use route groups for API versioning (`/api/v1`, `/api/v2`)

## API Reference

- **Spec**: [pkg.go.dev/github.com/oaswrap/spec](https://pkg.go.dev/github.com/oaswrap/spec)
- **HTTP Adapter**: [pkg.go.dev/github.com/oaswrap/spec/adapters/httpopenapi](https://pkg.go.dev/github.com/oaswrap/spec/adapters/httpopenapi)
- **Options**: [pkg.go.dev/github.com/oaswrap/spec/option](https://pkg.go.dev/github.com/oaswrap/spec/option)

## Contributing

We welcome contributions! Please open issues and PRs at the main [oaswrap/spec](https://github.com/oaswrap/spec) repository.

## License

[MIT License](LICENSE) — Created with ❤️ by [Ahmad Faiz](https://github.com/afkdevs).