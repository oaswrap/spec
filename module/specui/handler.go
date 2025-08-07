package specui

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/pkg/util"
	"github.com/swaggest/swgui"
	"github.com/swaggest/swgui/v5cdn"
)

// Handler provides methods to serve OpenAPI documentation and files.
type Handler struct {
	cfg       *openapi.Config
	generator spec.Generator
}

// NewHandler creates a new Handler instance with the provided OpenAPI generator.
func NewHandler(generator spec.Generator) *Handler {
	return &Handler{
		cfg:       generator.Config(),
		generator: generator,
	}
}

// DocsPath returns the path for the OpenAPI documentation.
func (h *Handler) DocsPath() string {
	return h.cfg.DocsPath
}

// DocsFilePath returns the path for the OpenAPI schema file.
func (h *Handler) DocsFilePath() string {
	return util.JoinPath(h.DocsPath(), Filename)
}

// Docs returns an HTTP handler that serves the OpenAPI documentation UI.
func (h *Handler) Docs() http.Handler {
	ui := v5cdn.NewHandlerWithConfig(h.swguiConfig())
	return ui
}

// DocsFunc returns an HTTP handler function that serves the OpenAPI documentation UI.
func (h *Handler) DocsFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ui := h.Docs()
		ui.ServeHTTP(w, r)
	}
}

func (h *Handler) swguiConfig() swgui.Config {
	cfg := h.cfg
	openapiPath := h.DocsFilePath()
	if cfg.BaseURL != "" {
		openapiPath = util.JoinURL(cfg.BaseURL, openapiPath)
	}

	return swgui.Config{
		Title:              cfg.Title,
		SwaggerJSON:        openapiPath,
		BasePath:           h.DocsPath(),
		ShowTopBar:         cfg.SwaggerConfig.ShowTopBar,
		HideCurl:           cfg.SwaggerConfig.HideCurl,
		JsonEditor:         cfg.SwaggerConfig.JsonEditor,
		PreAuthorizeApiKey: cfg.SwaggerConfig.PreAuthorizeApiKey,
		SettingsUI:         cfg.SwaggerConfig.SettingsUI,
		Proxy:              cfg.SwaggerConfig.Proxy,
	}
}

// DocsFile returns an HTTP handler that serves the OpenAPI schema file.
func (h *Handler) DocsFile() http.Handler {
	var (
		schema []byte
		err    error
		once   sync.Once
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		once.Do(func() {
			schema, err = h.generator.MarshalYAML()
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to generate OpenAPI schema"})
			return
		}

		w.Header().Set("Content-Type", "application/x-yaml")
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		_, err := w.Write(schema)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "failed to write OpenAPI schema"})
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

// DocsFileFunc returns an HTTP handler function that serves the OpenAPI schema file.
func (h *Handler) DocsFileFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ui := h.DocsFile()
		ui.ServeHTTP(w, r)
	}
}
