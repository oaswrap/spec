package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/ginopenapi/internal/handler"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)

	gen := spec.NewGenerator(
		option.WithDocsPath("/docs"),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
	)
	handler := handler.NewHandler(gen)
	assert.NotNil(t, handler)

	t.Run("Docs", func(t *testing.T) {
		e := gin.New()
		e.GET(handler.DocsPath(), handler.Docs)

		req := httptest.NewRequest(http.MethodGet, handler.DocsPath(), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, 200, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})
	t.Run("DocsFile", func(t *testing.T) {
		e := gin.New()
		e.GET(handler.DocsFilePath(), handler.DocsFile)

		req := httptest.NewRequest(http.MethodGet, handler.DocsFilePath(), nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, 200, rec.Code)
		assert.NotEmpty(t, rec.Body.String())
	})
}
