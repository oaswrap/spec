package handler_test

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/fiberopenapi/internal/handler"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestNewOpenAPIHandler(t *testing.T) {
	generator := spec.NewGenerator()
	cfg := generator.Config()

	handler := handler.NewOpenAPIHandler(cfg, generator)

	assert.NotNil(t, handler)
}

func TestOpenAPIHandler_OpenAPIYaml(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		generator := spec.NewGenerator()
		cfg := generator.Config()

		h := handler.NewOpenAPIHandler(cfg, generator)
		app := fiber.New()
		app.Get("/openapi.yaml", h.OpenAPIYaml)

		req := httptest.NewRequest("GET", "/openapi.yaml", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/x-yaml", resp.Header.Get(fiber.HeaderContentType))
	})
	t.Run("error", func(t *testing.T) {
		r := spec.NewGenerator()
		cfg := r.Config()

		r.Get("/user/{id}",
			option.Summary("Get User by ID"),
			option.Description("Retrieves a user by their unique ID"),
			option.Response(200, "User found"),
			option.Response(404, "User not found"),
		)
		h := handler.NewOpenAPIHandler(cfg, r)

		app := fiber.New()
		app.Get("/openapi.yaml", h.OpenAPIYaml)
		req := httptest.NewRequest("GET", "/openapi.yaml", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)
	})
}

func TestOpenAPIHandler_Docs(t *testing.T) {
	generator := spec.NewGenerator(option.WithSwaggerConfig(openapi.SwaggerConfig{}))
	cfg := generator.Config()

	h := handler.NewOpenAPIHandler(cfg, generator)
	app := fiber.New()
	app.Get("/docs", h.Docs)

	req := httptest.NewRequest("GET", "/docs", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/html", resp.Header.Get(fiber.HeaderContentType))

	generator = spec.NewGenerator(
		option.WithBaseURL("http://localhost:3000"),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
	)
	cfg = generator.Config()
	h = handler.NewOpenAPIHandler(cfg, generator)
	app = fiber.New()
	app.Get("/docs", h.Docs)
	req = httptest.NewRequest("GET", "/docs", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/html", resp.Header.Get(fiber.HeaderContentType))

	var bytes bytes.Buffer
	_, err = bytes.ReadFrom(resp.Body)
	assert.NoError(t, err)
	assert.Contains(t, bytes.String(), "Swagger UI")
	assert.Contains(t, bytes.String(), "http://localhost:3000/openapi.yaml")
}
