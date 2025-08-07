package specui_test

import (
	"net/http/httptest"
	"testing"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/module/specui"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Docs(t *testing.T) {
	gen := spec.NewGenerator(
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
	)
	handler := specui.NewHandler(gen)
	ui := handler.DocsFunc()
	assert.NotNil(t, ui)

	req := httptest.NewRequest("GET", "/docs", nil)
	rec := httptest.NewRecorder()
	ui.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), "API Documentation")
}

func TestHandler_Docs_BaseURL(t *testing.T) {
	gen := spec.NewGenerator(
		option.WithBaseURL("http://example.com"),
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
	)
	handler := specui.NewHandler(gen)

	req := httptest.NewRequest("GET", "/docs", nil)
	rec := httptest.NewRecorder()
	ui := handler.DocsFunc()
	ui.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.Contains(t, rec.Body.String(), "http://example.com/openapi.yaml")
}

func TestHandler_DocsFile(t *testing.T) {
	gen := spec.NewGenerator(
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
	)
	handler := specui.NewHandler(gen)

	req := httptest.NewRequest("GET", "/openapi.yaml", nil)
	rec := httptest.NewRecorder()
	ui := handler.DocsFileFunc()
	ui.ServeHTTP(rec, req)

	assert.Equal(t, 200, rec.Code)
	assert.NotEmpty(t, rec.Body.String())
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/x-yaml")
	assert.Contains(t, rec.Body.String(), "openapi: 3.0.3")
}

func TestHandler_DocsFile_Error(t *testing.T) {
	gen := spec.NewGenerator(
		option.WithSwaggerConfig(openapi.SwaggerConfig{}),
		option.WithOpenAPIVersion("2.0.0"),
	)
	handler := specui.NewHandler(gen)

	req := httptest.NewRequest("GET", "/openapi.yaml", nil)
	rec := httptest.NewRecorder()
	ui := handler.DocsFileFunc()
	ui.ServeHTTP(rec, req)

	assert.Equal(t, 500, rec.Code)
	assert.Contains(t, rec.Body.String(), "failed to generate OpenAPI schema")
}
