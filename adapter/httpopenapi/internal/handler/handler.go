package handler

import (
	"net/http"
	"sync"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/adapter/httpopenapi/internal/constant"
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

func (h *OpenAPIHandler) Docs(w http.ResponseWriter, r *http.Request) {
	ui := v5cdn.NewHandlerWithConfig(h.swguiConfig())
	ui.ServeHTTP(w, r)
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

func (h *OpenAPIHandler) OpenAPIYaml(w http.ResponseWriter, r *http.Request) {
	h.once.Do(func() {
		schema, err := h.generator.MarshalYAML()
		if err != nil {
			h.err = err
			http.Error(w, "Failed to generate OpenAPI schema", http.StatusInternalServerError)
			return
		}
		h.schema = schema
	})

	if h.err != nil {
		http.Error(w, h.err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/x-yaml")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(h.schema)
	if err != nil {
		http.Error(w, "Failed to write OpenAPI schema", http.StatusInternalServerError)
		return
	}
}
