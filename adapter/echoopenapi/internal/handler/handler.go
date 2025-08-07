package handler

import (
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/echoopenapi/internal/constant"
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

func (h *OpenAPIHandler) OpenAPIYaml(c echo.Context) error {
	h.once.Do(func() {
		h.schema, h.err = h.generator.MarshalYAML()
	})
	if h.err != nil {
		return c.JSON(500, map[string]string{
			"title":  "Internal Server Error",
			"detail": "Failed to generate OpenAPI schema",
		})
	}

	res := c.Response()
	res.Header().Set("Content-Type", "application/x-yaml")
	res.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	res.Header().Set("Pragma", "no-cache")
	res.Header().Set("Expires", "0")

	return c.Blob(200, "application/x-yaml", h.schema)
}

func (h *OpenAPIHandler) Docs(c echo.Context) error {
	ui := v5cdn.NewHandlerWithConfig(h.swguiConfig())
	ui.ServeHTTP(c.Response(), c.Request())
	return nil
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
