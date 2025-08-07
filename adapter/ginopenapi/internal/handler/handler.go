package handler

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/ginopenapi/internal/constant"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/pkg/util"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5cdn"
)

type OpenAPIHandler struct {
	cfg       *openapi.Config
	generator spec.Generator
	once      sync.Once
	err       error
	schema    []byte
}

func NewOpenAPIHandler(cfg *openapi.Config, generator spec.Generator) *OpenAPIHandler {
	return &OpenAPIHandler{
		cfg:       cfg,
		generator: generator,
	}
}

func (h *OpenAPIHandler) Docs(c *gin.Context) {
	ui := v5cdn.NewHandlerWithConfig(h.swguiConfig())
	gin.WrapH(ui)(c)
}

func (h *OpenAPIHandler) swguiConfig() swgui.Config {
	cfg := h.cfg
	openapiPath := util.JoinPath(cfg.DocsPath, constant.OpenAPIFileName)
	if cfg.BaseURL != "" {
		openapiPath = util.JoinURL(cfg.BaseURL, openapiPath)
	}

	return swgui.Config{
		Title:              cfg.Title,
		SwaggerJSON:        openapiPath,
		BasePath:           cfg.DocsPath,
		ShowTopBar:         cfg.SwaggerConfig.ShowTopBar,
		HideCurl:           cfg.SwaggerConfig.HideCurl,
		JsonEditor:         cfg.SwaggerConfig.JsonEditor,
		PreAuthorizeApiKey: cfg.SwaggerConfig.PreAuthorizeApiKey,
		SettingsUI:         cfg.SwaggerConfig.SettingsUI,
		Proxy:              cfg.SwaggerConfig.Proxy,
	}
}

func (h *OpenAPIHandler) OpenAPIYaml(c *gin.Context) {
	h.once.Do(func() {
		schema, err := h.generator.GenerateSchema("yaml")
		if err != nil {
			h.err = err
			return
		}
		h.schema = schema
	})
	if h.err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate OpenAPI schema",
		})
		return
	}

	// Set the response headers and content type.
	// This ensures that the response is correctly formatted as YAML.
	// The headers also prevent caching of the OpenAPI schema.
	c.Header("Content-Type", "application/yaml")
	c.Header("Cache-Control", "no-cache, no-store, must-revalidate")
	c.Header("Pragma", "no-cache")
	c.Header("Expires", "0")

	c.Writer.WriteHeader(http.StatusOK)
	_, err := c.Writer.Write(h.schema)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "failed to write OpenAPI schema",
		})
		return
	}
}
