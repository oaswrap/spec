package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestNewOpenAPIHandler(t *testing.T) {
	spec := spec.NewGenerator()
	cfg := spec.Config()

	handler := NewOpenAPIHandler(cfg, spec)

	assert.NotNil(t, handler)
	assert.Equal(t, cfg, handler.cfg)
	assert.Equal(t, spec, handler.generator)
}

func TestOpenAPIHandler_Docs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	spec := spec.NewGenerator(option.WithSwaggerConfig(openapi.SwaggerConfig{}))
	cfg := spec.Config()
	handler := NewOpenAPIHandler(cfg, spec)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/docs", nil)

	handler.Docs(c)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestOpenAPIHandler_Docs_BaseURL(t *testing.T) {
	gin.SetMode(gin.TestMode)

	spec := spec.NewGenerator(
		option.WithBaseURL("http://example.com"),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
	)
	cfg := spec.Config()
	handler := NewOpenAPIHandler(cfg, spec)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/docs", nil)

	handler.Docs(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "http://example.com/openapi.yaml")
}

func TestOpenAPIHandler_OpenAPIYaml_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	spec := spec.NewGenerator(option.WithSwaggerConfig(openapi.SwaggerConfig{}))
	cfg := spec.Config()

	handler := NewOpenAPIHandler(cfg, spec)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/openapi.yaml", nil)

	handler.OpenAPIYaml(c)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/yaml", w.Header().Get("Content-Type"))
	assert.Equal(t, "no-cache, no-store, must-revalidate", w.Header().Get("Cache-Control"))
	assert.Equal(t, "no-cache", w.Header().Get("Pragma"))
	assert.Equal(t, "0", w.Header().Get("Expires"))
}

func TestOpenAPIHandler_OpenAPIYaml_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	spec := spec.NewGenerator(option.WithOpenAPIVersion("2.0.0"), option.WithSwaggerConfig(openapi.SwaggerConfig{}))
	cfg := spec.Config()

	handler := NewOpenAPIHandler(cfg, spec)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/openapi.yaml", nil)

	handler.OpenAPIYaml(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "failed to generate OpenAPI schema")
}

func TestOpenAPIHandler_OpenAPIYaml_CacheOnce(t *testing.T) {
	gin.SetMode(gin.TestMode)
	spec := spec.NewGenerator(option.WithSwaggerConfig(openapi.SwaggerConfig{}))
	cfg := spec.Config()

	handler := NewOpenAPIHandler(cfg, spec)

	// First call
	w1 := httptest.NewRecorder()
	c1, _ := gin.CreateTestContext(w1)
	c1.Request = httptest.NewRequest("GET", "/openapi.yaml", nil)
	handler.OpenAPIYaml(c1)

	// Second call - should use cached result
	w2 := httptest.NewRecorder()
	c2, _ := gin.CreateTestContext(w2)
	c2.Request = httptest.NewRequest("GET", "/openapi.yaml", nil)
	handler.OpenAPIYaml(c2)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, http.StatusOK, w2.Code)
}
