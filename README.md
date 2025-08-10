# oaswrap/spec

[![CI](https://github.com/oaswrap/spec/actions/workflows/ci.yml/badge.svg)](https://github.com/oaswrap/spec/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/oaswrap/spec/graph/badge.svg?token=RIEIM9BAIW)](https://codecov.io/gh/oaswrap/spec)
[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec.svg)](https://pkg.go.dev/github.com/oaswrap/spec)
[![Go Report Card](https://goreportcard.com/badge/github.com/oaswrap/spec)](https://goreportcard.com/report/github.com/oaswrap/spec)
[![Go Version](https://img.shields.io/github/go-mod/go-version/oaswrap/spec)](https://github.com/oaswrap/spec/blob/main/go.mod)
[![License](https://img.shields.io/github/license/oaswrap/spec)](LICENSE)

A lightweight, framework-agnostic OpenAPI 3.x specification builder for Go that gives you complete control over your API documentation without vendor lock-in.

## Why oaswrap/spec?

- **🎯 Framework Agnostic** — Works with any Go web framework or as a standalone tool
- **⚡ Zero Dependencies** — Powered by [`swaggest/openapi-go`](https://github.com/swaggest/openapi-go) with minimal overhead
- **🔧 Programmatic Control** — Build specs in pure Go code with full type safety
- **📚 Built-in Documentation** — Multiple UI options: Swagger UI, Redoc, and Stoplight Elements
- **🚀 Adapter Ecosystem** — Seamless integration with popular frameworks via dedicated adapters
- **📝 CI/CD Ready** — Generate specs at build time for documentation pipelines

## Installation

```bash
go get github.com/oaswrap/spec
```

## Quick Start

### Static Spec Generation

For CI/CD pipelines and build-time spec generation:

```go
package main

import (
	"log"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

func main() {
	// Create a new OpenAPI router
	r := spec.NewRouter(
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithServer("https://api.example.com"),
	)

	// Add routes
	v1 := r.Group("/api/v1")
	
	v1.Post("/login",
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)

	v1.Get("/users/{id}",
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
```

📖 **[View the generated spec](https://rest.wiki/?https://raw.githubusercontent.com/oaswrap/spec/main/examples/basic/openapi.yaml)** on Rest.Wiki

### Basic Usage with Built-in Documentation

Perfect for adding OpenAPI documentation to any existing Go HTTP server:

```go
package main

import (
	"log"
	"net/http"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/option"
)

func main() {
	// Create a new OpenAPI router
	r := spec.NewRouter(
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithServer("https://api.example.com"),
	)

	// Add spec routes
	v1 := r.Group("/api/v1")

	v1.Post("/login",
		option.Summary("User login"),
		option.Request(new(LoginRequest)),
		option.Response(200, new(LoginResponse)),
	)

	v1.Get("/users/{id}",
		option.Summary("Get user by ID"),
		option.Request(new(GetUserRequest)),
		option.Response(200, new(User)),
	)

	// Get configuration for doc paths
	cfg := r.Config()

	// Setup HTTP handlers
	mux := http.NewServeMux()
	
	// Add built-in documentation handlers
	mux.HandleFunc("GET "+cfg.DocsPath, r.DocsHandlerFunc())
	mux.HandleFunc("GET "+cfg.SpecPath, r.SpecHandlerFunc())

	// Add your actual API routes here
	// mux.HandleFunc("POST /api/v1/login", loginHandler)
	// mux.HandleFunc("GET /api/v1/users/{id}", getUserHandler)

	log.Printf("🚀 OpenAPI docs available at: http://localhost:3000%s", cfg.DocsPath)
	log.Printf("📄 OpenAPI spec available at: http://localhost:3000%s", cfg.SpecPath)

	http.ListenAndServe(":3000", mux)
}
```

### Framework Integration

For seamless HTTP server integration, use one of our framework adapters:

| Framework | Package |
|-----------|---------|
| **Chi** | [`chiopenapi`](/adapter/chiopenapi) |
| **Echo** | [`echoopenapi`](/adapter/echoopenapi) |
| **Gin** | [`ginopenapi`](/adapter/ginopenapi) |
| **Fiber** | [`fiberopenapi`](/adapter/fiberopenapi) |
| **HTTP** | [`httpopenapi`](/adapter/httpopenapi) |
| **Mux** | [`muxopenapi`](/adapter/muxopenapi) |

Each adapter provides:
- ✅ Automatic spec generation from your routes
- 📚 Built-in documentation UI at `/docs`
- 📄 YAML spec endpoints at `/docs/openapi.yaml`
- 🔧 Inline OpenAPI options with route definitions

Visit the individual adapter repositories for framework-specific examples and detailed integration guides.

## Built-in Documentation Handlers

The `spec` package includes built-in handlers for serving OpenAPI documentation with multiple UI options, powered by [`oaswrap/spec-ui`](https://github.com/oaswrap/spec-ui):

### Supported Documentation UIs

- **Stoplight Elements** — Modern, component-based documentation (default)
- **Swagger UI** — Interactive API documentation
- **Redoc** — Clean, responsive documentation

### Handler Functions
```go
// DocsHandlerFunc returns a handler for interactive API documentation
// Supports multiple UI types: Stoplight Elements (default), Swagger UI, Redoc
mux.HandleFunc("GET /docs", r.DocsHandlerFunc())

// SpecHandlerFunc returns a handler for the OpenAPI spec (YAML format)
mux.HandleFunc("GET /docs/openapi.yaml", r.SpecHandlerFunc())
```

### UI Configuration
Choose your preferred documentation UI:

```go
r := spec.NewRouter(
	option.WithTitle("My API"),
	option.WithStoplightElements(),			// Available UI options: WithStoplightElements(), WithSwaggerUI(), WithRedoc(),
	option.WithDocsPath("/docs"),			// Custom docs path (default: "/docs")
	option.WithSpecPath("/openapi.yaml"),	// Custom spec path (default: "/docs/openapi.yaml")
)

cfg := r.Config()
// cfg.DocsPath == "/documentation"
// cfg.SpecPath == "/api/openapi.yaml"
```

## When to Use What?

### ✅ Use `spec` for static generation when you:
- Generate OpenAPI files at **build time**
- Integrate with **CI/CD pipelines**
- Build **custom documentation tools**
- Need **static spec generation**

### ✅ Use `spec` with built-in handlers when you:
- Want **lightweight documentation** for existing APIs
- Need **framework independence**
- Prefer **manual route documentation**
- Want to **add docs to legacy systems**
- Need **custom HTTP handler integration**

### ✅ Use framework adapters when you:
- Want **automatic spec generation** from routes
- Need **zero-configuration setup**
- Prefer **inline OpenAPI configuration**
- Want **route registration + documentation** in one step

## Configuration Options

The `option` package provides comprehensive OpenAPI configuration:

### Basic Information
```go
option.WithOpenAPIVersion("3.0.3") // Default: "3.0.3"
option.WithTitle("My API")
option.WithDescription("API description")
option.WithVersion("1.2.3")
option.WithContact(openapi.Contact{
	Name:  "Support Team",
	URL:   "https://support.example.com",
	Email: "support@example.com",
})
option.WithLicense(openapi.License{
	Name: "MIT License",
	URL:  "https://opensource.org/licenses/MIT",
})
option.WithExternalDocs("https://docs.example.com", "API Documentation")
option.WithTags(
	openapi.Tag{
		Name:        "User Management",
		Description: "Operations related to user management",
	},
	openapi.Tag{
		Name:        "Authentication", 
		Description: "Authentication related operations",
	},
)
```

### Servers
```go
option.WithServer("https://api.example.com")
option.WithServer("https://api-example.com/{version}",
	option.ServerDescription("Production Server"),
	option.ServerVariables(map[string]openapi.ServerVariable{
		"version": {
			Default:     "v1",
			Enum:        []string{"v1", "v2"},
			Description: "API version",
		},
	}),
)
```

### Security Schemes
```go
// Bearer token
option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer"))

// API Key
option.WithSecurity("apiKey", option.SecurityAPIKey("X-API-Key", "header"))

// OAuth2
option.WithSecurity("oauth2", option.SecurityOAuth2(
	openapi.OAuthFlows{
		Implicit: &openapi.OAuthFlowsImplicit{
			AuthorizationURL: "https://auth.example.com/authorize",
			Scopes: map[string]string{
				"read":  "Read access",
				"write": "Write access",
			},
		},
	},
))
```

### Route Documentation
```go
option.OperationID("getUserByID")					// Unique operation ID
option.Summary("Short description")					// Brief summary
option.Description("Detailed description")			// Full description
option.Tags("User Management", "Authentication")	// Group by tags
option.Request(new(RequestModel))					// Request body model
option.Response(200, new(ResponseModel),			// Response model
	option.ContentDescription("Successful response"),
	option.ContentType("application/json"),
	option.ContentDefault(true),
)
option.Security("bearerAuth")					// Apply security scheme
option.Deprecated()								// Mark as deprecated
option.Hidden()									// Hide from spec
```

### Parameter Definition
Define parameters using struct tags in your request models:

```go
type GetUserRequest struct {
	ID     string `path:"id" required:"true" description:"User identifier"`
	Limit  int    `query:"limit" description:"Maximum number of results"`
	APIKey string `header:"X-API-Key" description:"API authentication key"`
}
```

### Group-Level Configuration
Apply settings to all routes within a group:

```go
// Apply to all routes in the group
adminGroup := r.Group("/admin",
	option.GroupTags("Administration"),
	option.GroupSecurity("bearerAuth"),
	option.GroupDeprecated(),
)

// Hide internal routes from documentation
internalGroup := r.Group("/internal",
	option.GroupHidden(),
)
```

## Advanced Features

### Rich Schema Documentation
```go
type CreateUserRequest struct {
	Name     string   `json:"name" required:"true" minLength:"2" maxLength:"50"`
	Email    string   `json:"email" required:"true" format:"email"`
	Age      int      `json:"age" minimum:"18" maximum:"120"`
	Tags     []string `json:"tags" maxItems:"10"`
}
```

For comprehensive struct tag documentation, see [swaggest/openapi-go](https://github.com/swaggest/openapi-go?tab=readme-ov-file#features) and [swaggest/jsonschema-go](https://github.com/swaggest/jsonschema-go?tab=readme-ov-file#field-tags).

### Generic Response Types
```go
type APIResponse[T any] struct {
	Success   bool   `json:"success"`
	Data      T      `json:"data,omitempty"`
	Error     string `json:"error,omitempty"`
	Timestamp string `json:"timestamp"`
}

// Usage
option.Response(200, new(APIResponse[User]))
option.Response(200, new(APIResponse[[]Product]))
```

## Examples

Explore complete working examples in the [`examples/`](examples/) directory:

- **[Basic](examples/basic/)** — Standalone spec generation
- **[Basic-HTTP](examples/basic-http/)** — Built-in HTTP server with OpenAPI docs
- **[Petstore](examples/petstore/)** — Full Petstore API with routes and models

## API Reference

Complete documentation at [pkg.go.dev/github.com/oaswrap/spec](https://pkg.go.dev/github.com/oaswrap/spec).

Key packages:
- [`spec`](https://pkg.go.dev/github.com/oaswrap/spec) — Core router and spec builder
- [`option`](https://pkg.go.dev/github.com/oaswrap/spec/option) — Configuration options

## FAQ

**Q: Can I use this with my existing API?**  
A: Absolutely! Use the standalone version to document existing APIs, or gradually migrate to framework adapters.

**Q: How does this compare to swag/swaggo?**  
A: While swag uses code comments, oaswrap uses pure Go code for type safety and better IDE support. Both have their merits - swag is annotation-based while oaswrap is code-first.

**Q: How does this compare to Huma?**  
A: Both are excellent choices with different philosophies:
- **Huma** is a complete HTTP framework with built-in OpenAPI generation, validation, and middleware
- **oaswrap/spec** is a lightweight, framework-agnostic documentation builder that works with your existing setup
- Use **Huma** if you're building a new API and want an all-in-one solution with automatic validation
- Use **oaswrap** if you have existing code, prefer framework flexibility, or need standalone spec generation

**Q: Is this production ready?**  
A: The library is in active development. While core functionality is solid, consider it beta software. Thorough testing is recommended before production use.

**Q: How do I handle authentication in the generated docs?**  
A: Define security schemes using `option.WithSecurity()` and apply them to routes with `option.Security()`. The generated docs will include authentication UI.

**Q: Can I customize the documentation UI?**  
A: Yes! The built-in handlers support three different UI options: Stoplight Elements (default), Swagger UI, and Redoc. Use `option.WithSwaggerUI()`, `option.WithRedoc()`, or `option.WithStoplightElements()` to choose your preferred interface. All UIs are powered by [`oaswrap/spec-ui`](https://github.com/oaswrap/spec-ui) for consistent, modern documentation experiences.

## Contributing

We welcome contributions! Here's how you can help:

1. **🐛 Report bugs** — Open an issue with reproduction steps
2. **💡 Suggest features** — Share your ideas for improvements
3. **📝 Improve docs** — Help make our documentation clearer
4. **🔧 Submit PRs** — Fix bugs or add features

Please check existing issues and discussions before starting work on new features.

## License

[MIT License](LICENSE) — Created with ❤️ by [Ahmad Faiz](https://github.com/afkdevs).