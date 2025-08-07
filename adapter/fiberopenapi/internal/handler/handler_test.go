package handler_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/fiberopenapi/internal/handler"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	gen := spec.NewGenerator(
		option.WithDocsPath("/docs"),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
	)
	handler := handler.NewHandler(gen)
	assert.NotNil(t, handler)

	t.Run("Docs", func(t *testing.T) {
		app := fiber.New()
		app.Get(handler.DocsPath(), handler.Docs)

		req := httptest.NewRequest(fiber.MethodGet, handler.DocsPath(), nil)
		rec, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, 200, rec.StatusCode)
	})
	t.Run("DocsFile", func(t *testing.T) {
		app := fiber.New()
		app.Get(handler.DocsFilePath(), handler.DocsFile)

		req := httptest.NewRequest(fiber.MethodGet, handler.DocsFilePath(), nil)
		rec, err := app.Test(req, -1)
		assert.NoError(t, err)
		assert.Equal(t, 200, rec.StatusCode)
	})
}
