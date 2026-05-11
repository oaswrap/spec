# irisopenapi

[![Go Reference](https://pkg.go.dev/badge/github.com/oaswrap/spec/adapter/irisopenapi.svg)](https://pkg.go.dev/github.com/oaswrap/spec/adapter/irisopenapi)
[![Go Report Card](https://goreportcard.com/badge/github.com/oaswrap/spec/adapter/irisopenapi)](https://goreportcard.com/report/github.com/oaswrap/spec/adapter/irisopenapi)

A lightweight adapter for the [Iris](https://github.com/kataras/iris) web framework that automatically generates OpenAPI 3.x specifications from your routes using [`oaswrap/spec`](https://github.com/oaswrap/spec).

## Installation

```bash
go get github.com/oaswrap/spec/adapter/irisopenapi
```

## Quick Start

```go
package main

import (
	"github.com/kataras/iris/v12"
	"github.com/oaswrap/spec/adapter/irisopenapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	app := iris.New()
	r := irisopenapi.NewRouter(app,
		option.WithTitle("My API"),
		option.WithVersion("1.0.0"),
	)

	r.Get("/ping", func(ctx iris.Context) {
		_ = ctx.JSON(map[string]string{"message": "pong"})
	}).With(
		option.OperationID("ping"),
		option.Response(200, new(struct {
			Message string `json:"message"`
		})),
	)

	_ = app.Listen(":8080")
}
```

## Built-in Endpoints

- `/docs` for interactive API docs UI.
- `/docs/openapi.yaml` for raw OpenAPI YAML.

Use `option.WithDisableDocs(true)` to disable docs endpoints.
