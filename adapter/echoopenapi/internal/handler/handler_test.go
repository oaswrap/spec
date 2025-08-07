package handler_test

import (
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/echoopenapi/internal/handler"
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
		e := echo.New()
		e.GET(handler.DocsPath(), handler.Docs)

		req := httptest.NewRequest(echo.GET, handler.DocsPath(), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, 200, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})
	t.Run("DocsFile", func(t *testing.T) {
		e := echo.New()
		e.GET(handler.DocsFilePath(), handler.DocsFile)

		req := httptest.NewRequest(echo.GET, handler.DocsFilePath(), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, 200, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})
}
