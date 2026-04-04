# Framework Adapters

This directory contains framework-specific adapters for `oaswrap/spec` that provide seamless integration with popular Go web frameworks. Each adapter automatically generates OpenAPI 3.x specifications from your existing routes and handlers.

## Available Adapters

### Web Frameworks

| Framework | Adapter | Go Module | Description |
|-----------|---------|-----------|-------------|
| [Chi](https://github.com/go-chi/chi) | [`chiopenapi`](./chiopenapi) | `github.com/oaswrap/spec/adapter/chiopenapi` | Lightweight router with middleware support |
| [Echo v4](https://github.com/labstack/echo) | [`echoopenapi`](./echoopenapi) | `github.com/oaswrap/spec/adapter/echoopenapi` | High performance, extensible, minimalist framework |
| [Echo v5](https://github.com/labstack/echo) | [`echov5openapi`](./echov5openapi) | `github.com/oaswrap/spec/adapter/echov5openapi` | Echo v5 with updated Context API |
| [Fiber v2](https://github.com/gofiber/fiber) | [`fiberopenapi`](./fiberopenapi) | `github.com/oaswrap/spec/adapter/fiberopenapi` | Express-inspired framework built on Fasthttp |
| [Fiber v3](https://github.com/gofiber/fiber) | [`fiberv3openapi`](./fiberv3openapi) | `github.com/oaswrap/spec/adapter/fiberv3openapi` | Fiber v3 with updated Ctx interface and binding API |
| [Gin](https://github.com/gin-gonic/gin) | [`ginopenapi`](./ginopenapi) | `github.com/oaswrap/spec/adapter/ginopenapi` | Fast HTTP web framework with zero allocation |
| [net/http](https://pkg.go.dev/net/http) | [`httpopenapi`](./httpopenapi) | `github.com/oaswrap/spec/adapter/httpopenapi` | Standard library HTTP package |
| [HttpRouter](https://github.com/julienschmidt/httprouter) | [`httprouteropenapi`](./httprouteropenapi) | `github.com/oaswrap/spec/adapter/httprouteropenapi` | High performance HTTP request router |
| [Gorilla Mux](https://github.com/gorilla/mux) | [`muxopenapi`](./muxopenapi) | `github.com/oaswrap/spec/adapter/muxopenapi` | Powerful HTTP router and URL matcher |