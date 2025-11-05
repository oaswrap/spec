# Framework Adapters

This directory contains framework-specific adapters for `oaswrap/spec` that provide seamless integration with popular Go web frameworks. Each adapter automatically generates OpenAPI 3.x specifications from your existing routes and handlers.

## Available Adapters

### Web Frameworks

| Framework | Adapter | Go Module | Description |
|-----------|---------|-----------|-------------|
| [Chi](https://github.com/go-chi/chi) | [`chiopenapi`](./chiopenapi) | `github.com/oaswrap/spec/adapter/chiopenapi` | Lightweight router with middleware support |
| [Echo](https://github.com/labstack/echo) | [`echoopenapi`](./echoopenapi) | `github.com/oaswrap/spec/adapter/echoopenapi` | High performance, extensible, minimalist framework |
| [Fiber](https://github.com/gofiber/fiber) | [`fiberopenapi`](./fiberopenapi) | `github.com/oaswrap/spec/adapter/fiberopenapi` | Express-inspired framework built on Fasthttp |
| [Gin](https://github.com/gin-gonic/gin) | [`ginopenapi`](./ginopenapi) | `github.com/oaswrap/spec/adapter/ginopenapi` | Fast HTTP web framework with zero allocation |
| [net/http](https://pkg.go.dev/net/http) | [`httpopenapi`](./httpopenapi) | `github.com/oaswrap/spec/adapter/httpopenapi` | Standard library HTTP package |
| [HttpRouter](https://github.com/julienschmidt/httprouter) | [`httprouteropenapi`](./httprouteropenapi) | `github.com/oaswrap/spec/adapter/httprouteropenapi` | High performance HTTP request router |
| [Gorilla Mux](https://github.com/gorilla/mux) | [`muxopenapi`](./muxopenapi) | `github.com/oaswrap/spec/adapter/muxopenapi` | Powerful HTTP router and URL matcher |