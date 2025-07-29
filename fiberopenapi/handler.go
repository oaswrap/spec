package fiberopenapi

import (
	"path"

	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5cdn"
)

const openapiFileName = "openapi.yaml"

func swguiConfig(cfg *openapiwrapper.Config) swgui.Config {
	openapiPath := path.Join(cfg.DocsPath, openapiFileName)
	if cfg.BaseURL != "" {
		openapiPath = cfg.BaseURL + openapiPath
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

func docsHandler(cfg *openapiwrapper.Config) fiber.Handler {
	ui := v5cdn.NewHandlerWithConfig(swguiConfig(cfg))
	return adaptor.HTTPHandler(ui)
}

func openAPIHandler(r *router) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Validate the OpenAPI schema before serving it.
		// This ensures that the schema is valid and can be used by clients.
		if err := r.Validate(); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "OpenAPI schema validation failed")
		}
		// Generate the OpenAPI schema only once.
		// This avoids regenerating it on every request, which can be expensive.
		r.schemaOnce.Do(func() {
			schema, err := r.GenerateOpenAPISchema("yaml")
			if err != nil {
				r.errors.Add(err)
				return
			}
			r.schema = schema
		})
		// If there are errors during schema generation, return an error response.
		// Otherwise, serve the generated schema.
		if r.errors.HasErrors() {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate OpenAPI schema")
		}

		// Set the response headers and content type.
		// This ensures that the response is correctly formatted as YAML.
		// The headers also prevent caching of the OpenAPI schema.
		c.Set(fiber.HeaderContentType, "application/x-yaml")
		c.Set(fiber.HeaderCacheControl, "no-cache, no-store, must-revalidate")
		c.Set(fiber.HeaderPragma, "no-cache")
		c.Set(fiber.HeaderExpires, "0")
		return c.Send(r.schema)
	}
}
