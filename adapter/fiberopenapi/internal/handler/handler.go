package handler

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/fiberopenapi/internal/constant"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/pkg/util"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5cdn"
)

type OpenAPIHandler struct {
	cfg    *openapi.Config
	gen    spec.Generator
	once   sync.Once
	err    error
	schema []byte
}

func NewOpenAPIHandler(cfg *openapi.Config, gen spec.Generator) *OpenAPIHandler {
	return &OpenAPIHandler{
		cfg: cfg,
		gen: gen,
	}
}

func (h *OpenAPIHandler) OpenAPIYaml(c *fiber.Ctx) error {
	h.once.Do(func() {
		schema, err := h.gen.MarshalYAML()
		if err != nil {
			h.err = err
			return
		}
		h.schema = schema
	})
	if h.err != nil {
		return fiber.ErrInternalServerError
	}

	// Set the response headers and content type.
	// This ensures that the response is correctly formatted as YAML.
	// The headers also prevent caching of the OpenAPI schema.
	c.Set(fiber.HeaderContentType, "application/x-yaml")
	c.Set(fiber.HeaderCacheControl, "no-cache, no-store, must-revalidate")
	c.Set(fiber.HeaderPragma, "no-cache")
	c.Set(fiber.HeaderExpires, "0")

	return c.Send(h.schema)
}

func (h *OpenAPIHandler) Docs(c *fiber.Ctx) error {
	ui := v5cdn.NewHandlerWithConfig(h.swguiConfig())
	return adaptor.HTTPHandler(ui)(c)
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
