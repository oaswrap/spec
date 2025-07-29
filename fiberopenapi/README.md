# fiberopenapi

**`fiberopenapi`** is the **Fiber adapter** for [OpenAPI Wrapper](https://github.com/faizlabs/openapi-wrapper).  
It wraps a `fiber.Router` to help you generate an OpenAPI 3.1 specification alongside your Fiber routes, using a simple, chainable API.

---

## ✨ Features

- 📚 **Automatic OpenAPI generation** — no separate files, write docs inline.
- ⚡ **Wraps Fiber’s router** — works just like `fiber.Router` but adds `.With(...)` for OpenAPI metadata.
- 🏷️ **Supports groups** — define nested routes and tags.
- 🔐 **Supports security schemes** — API keys, bearer tokens.
- 📝 **Exports OpenAPI 3.1 spec** — as YAML or JSON.

---

## 📦 Install

```bash
go get github.com/faizlabs/openapi-wrapper/fiberopenapi
```

---

## 🚀 Example

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

	r := fiberopenapi.NewRouter(app,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
		option.WithDescription("Sample API"),
		option.WithDocsPath("/docs"),
		option.WithServer("http://localhost:3000", "Local server"),
	)

	r.Get("/ping", func(c *fiber.Ctx) error {
		return c.SendString("pong")
	}).With(
		option.Summary("Ping"),
		option.Description("Health check endpoint"),
	)

	// Validate & write spec (optional)
	if err := r.Validate(); err != nil {
		log.Fatalf("OpenAPI validation failed: %v", err)
	}
	r.WriteSchemaTo("openapi.yaml")

	app.Listen(":3000")
}
```

---

## 📚 Documentation

📖 Full API docs for `fiberopenapi` are on [pkg.go.dev](https://pkg.go.dev/github.com/faizlabs/openapi-wrapper/fiberopenapi).

- See `ConfigOption` helpers like `WithTitle`, `WithServer`, `WithSecurity`.
- Use `Route.With(...)` to attach summaries, descriptions, responses.
- Combine with the `option` package for reusable helpers.

---

## 📂 Related

- 🔗 [Core OpenAPI generator](https://pkg.go.dev/github.com/faizlabs/openapi-wrapper/core)  
- 🔗 [Reusable OpenAPI options](https://pkg.go.dev/github.com/faizlabs/openapi-wrapper/option)  
- 🔗 [Main repo README](https://github.com/faizlabs/openapi-wrapper)

---

## ✅ License

MIT — use freely, PRs welcome!

---

**Build better documented Fiber APIs — with `fiberopenapi`.** 🚀