package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapters/echoopenapi/internal/handler"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestOpenAPIHandler_OpenAPIYaml_Success(t *testing.T) {
	// Setup
	spec := spec.NewGenerator(
		option.WithTitle("Test API"),
		option.WithDescription("This is a test API"),
		option.WithVersion("1.0.0"),
	)
	cfg := spec.Config()
	handler := handler.NewOpenAPIHandler(cfg, spec)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.OpenAPIYaml(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/x-yaml", rec.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache, no-store, must-revalidate", rec.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", rec.Header().Get("Pragma"))
	assert.Equal(t, "0", rec.Header().Get("Expires"))
}

func TestOpenAPIHandler_OpenAPIYaml_GeneratorError(t *testing.T) {
	// Setup
	spec := spec.NewGenerator(
		option.WithTitle("Test API"),
		option.WithDescription("This is a test API"),
		option.WithVersion("1.0.0"),
		option.WithOpenAPIVersion("2.0.0"),
	)
	cfg := spec.Config()
	handler := handler.NewOpenAPIHandler(cfg, spec)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.OpenAPIYaml(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.Contains(t, rec.Body.String(), "Internal Server Error")
	assert.Contains(t, rec.Body.String(), "Failed to generate OpenAPI schema")
}

func TestOpenAPIHandler_OpenAPIYaml_OnceCall(t *testing.T) {
	// Setup
	spec := spec.NewGenerator(
		option.WithTitle("Test API"),
		option.WithDescription("This is a test API"),
		option.WithVersion("1.0.0"),
	)
	cfg := spec.Config()
	handler := handler.NewOpenAPIHandler(cfg, spec)

	e := echo.New()

	// Execute multiple times
	for range 3 {
		req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handler.OpenAPIYaml(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestOpenAPIHandler_Docs(t *testing.T) {
	// Setup
	spec := spec.NewGenerator(
		option.WithTitle("Test API"),
		option.WithDescription("This is a test API"),
		option.WithVersion("1.0.0"),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
		option.WithDocsPath("/docs"),
	)
	cfg := spec.Config()
	handler := handler.NewOpenAPIHandler(cfg, spec)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.Docs(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
	assert.NotEmpty(t, rec.Body.String())
	assert.Contains(t, rec.Body.String(), "Swagger UI")
	assert.Contains(t, rec.Body.String(), "/docs/openapi.yaml")
}

func TestOpenAPIHandler_Docs_BaseURL(t *testing.T) {
	// Setup
	spec := spec.NewGenerator(
		option.WithTitle("Test API"),
		option.WithDescription("This is a test API"),
		option.WithVersion("1.0.0"),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
		option.WithBaseURL("http://localhost:3000"),
		option.WithDocsPath("/docs"),
	)
	cfg := spec.Config()
	handler := handler.NewOpenAPIHandler(cfg, spec)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Execute
	err := handler.Docs(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "text/html")
	assert.NotEmpty(t, rec.Body.String())
	assert.Contains(t, rec.Body.String(), "Swagger UI")
	assert.Contains(t, rec.Body.String(), "http://localhost:3000/docs/openapi.yaml")
}
