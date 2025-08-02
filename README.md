# oaswrap/spec

[![CI](https://github.com/oaswrap/spec/actions/workflows/ci.yml/badge.svg)](https://github.com/oaswrap/spec/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/oaswrap/spec/branch/main/graph/badge.svg)](https://codecov.io/gh/oaswrap/spec)
[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec.svg)](https://pkg.go.dev/github.com/oaswrap/spec)
[![Go Report Card](https://goreportcard.com/badge/github.com/oaswrap/spec)](https://goreportcard.com/report/github.com/oaswrap/spec)
[![License](https://img.shields.io/github/license/oaswrap/spec)](LICENSE)

A lightweight, framework-agnostic OpenAPI 3.x specification builder for Go that gives you complete control over your API documentation without vendor lock-in.

## Why oaswrap/spec?

- **🎯 Framework Agnostic** — Works with any Go web framework or as a standalone tool
- **⚡ Zero Dependencies** — Powered by [`swaggest/openapi-go`](https://github.com/swaggest/openapi-go) with minimal overhead
- **🔧 Programmatic Control** — Build specs in pure Go code with full type safety
- **🚀 Adapter Ecosystem** — Seamless integration with popular frameworks via dedicated adapters
- **📝 CI/CD Ready** — Generate specs at build time for documentation pipelines

## Quick Start

### Installation

```bash
go get github.com/oaswrap/spec
```

### Basic Usage (Standalone)

Perfect for generating OpenAPI specs in CI/CD, build scripts, or documentation tools:

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

## Framework Integration

For seamless HTTP server integration, use one of our framework adapters. Each adapter has its own repository with complete examples and documentation.

### Available Adapters

| Framework | Adapter Package |
|-----------|-----------------|
| **Gin** | [oaswrap/ginopenapi](https://github.com/oaswrap/ginopenapi) |
| **Echo** | [oaswrap/echoopenapi](https://github.com/oaswrap/echoopenapi) |
| **Fiber** | [oaswrap/fiberopenapi](https://github.com/oaswrap/fiberopenapi) |

Each adapter provides:
- Automatic spec generation from your routes
- Built-in Swagger UI documentation
- JSON/YAML spec endpoints
- Inline OpenAPI options with route definitions

Visit the individual adapter repositories for framework-specific examples and detailed integration guides.

## Configuration Options

The `option` package provides comprehensive OpenAPI configuration:

### Basic Information
```go
option.WithOpenAPIVersion("3.0.3") // Specify OpenAPI version (default is "3.0.3")
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
option.Tags(
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
option.OperationID("getUserByID") // Specify unique operation ID
option.Summary("Short description") // Brief summary of the operation
option.Description("Detailed description") // Full description of the operation
option.Tags("User Management", "Authentication") // Group operations by tags
option.Request(new(RequestModel)) // Define request body model
option.Response(200, new(ResponseModel), // Define response model
	option.ContentDescription("Successful response"), // Add description for response
	option.ContentType("application/json"), // Specify content type
	option.ContentDefault(true), // Mark as default response
)
option.Security("bearerAuth") // Apply security scheme to the route
option.Deprecated() // Mark route as deprecated
option.Hidden()     // Hide route from OpenAPI spec
```

Parameters (path, query, headers) are defined using struct tags in your request models:

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
// Apply tags, security, and other settings to all routes in the group
adminGroup := r.Group("/admin",
	option.GroupTags("Administration"),
	option.GroupSecurity("bearerAuth"),
	option.GroupDeprecated(), // Mark all routes as deprecated
)

// Hide internal routes from documentation
internalGroup := r.Group("/internal",
	option.GroupHidden(), // Exclude from OpenAPI spec
)
```

## Use Cases

### ✅ Use `spec` standalone when you:
- Generate OpenAPI files at **build time**
- Integrate with **CI/CD pipelines**
- Build **custom documentation tools**
- Need **static spec generation**
- Want **framework independence**

### ✅ Use framework adapters when you:
- Want **automatic spec generation** from routes
- Need **built-in documentation UI**
- Prefer **inline OpenAPI configuration**
- Want **live spec endpoints**

## Advanced Features

### Rich Schema Documentation
```go
// Use struct tags to generate detailed OpenAPI schemas
type CreateUserRequest struct {
	Name     string   `json:"name" required:"true" minLength:"2" maxLength:"50"`
	Email    string   `json:"email" required:"true" format:"email"`
	Age      int      `json:"age" minimum:"18" maximum:"120"`
	Tags     []string `json:"tags" maxItems:"10"`
}
```

For comprehensive struct tag documentation and advanced schema features, see the [swaggest/openapi-go features guide](https://github.com/swaggest/openapi-go?tab=readme-ov-file#features) and [swaggest/jsonschema-go field tags](https://github.com/swaggest/jsonschema-go?tab=readme-ov-file#field-tags).

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

Check out the [`examples/`](examples/) directory for complete working examples:

- **[Basic](examples/basic/)** — Standalone spec generation

For framework-specific examples, visit the individual adapter repositories:
- **Gin examples** — See [oaswrap/ginopenapi](https://github.com/oaswrap/ginopenapi)
- **Echo examples** — See [oaswrap/echoopenapi](https://github.com/oaswrap/echoopenapi)  
- **Fiber examples** — See [oaswrap/fiberopenapi](https://github.com/oaswrap/fiberopenapi)

## API Reference

For complete API documentation, visit [pkg.go.dev/github.com/oaswrap/spec](https://pkg.go.dev/github.com/oaswrap/spec).

Key packages:
- [`spec`](https://pkg.go.dev/github.com/oaswrap/spec) — Core router and spec builder
- [`option`](https://pkg.go.dev/github.com/oaswrap/spec/option) — All configuration options

## FAQ

**Q: Can I use this with my existing API?**  
A: Yes! You can either use the standalone version to document existing APIs, or gradually migrate to framework adapters.

**Q: How does this compare to swag/swaggo?**  
A: While swag uses code comments, oaswrap uses pure Go code for type safety and better IDE support. Both approaches have their merits.

**Q: Can I customize the generated documentation UI?**  
A: Framework adapters provide built-in UIs, but you can also serve the spec with any OpenAPI-compatible documentation tool.

**Q: Is this production ready?**  
A: The library is in active development. While the core functionality is solid, consider it beta software. We recommend thorough testing before production use.

## Roadmap

- [ ] Chi adapter
- [ ] HTTP adapter
- [ ] Stoplight support
- [ ] Redoc UI support

## Contributing

We welcome contributions! Here's how you can help:

1. **🐛 Report bugs** — Open an issue with reproduction steps
2. **💡 Suggest features** — Share your ideas for improvements
3. **📝 Improve docs** — Help make our documentation clearer
4. **🔧 Submit PRs** — Fix bugs or add features

Please check our issues and discussions before starting work on new features.

## License

[MIT License](LICENSE) — Created with ❤️ by [Ahmad Faiz](https://github.com/afkdevs).