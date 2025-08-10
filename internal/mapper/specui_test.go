package mapper_test

import (
	"testing"

	"github.com/oaswrap/spec-ui/config"
	"github.com/oaswrap/spec/internal/mapper"
	"github.com/oaswrap/spec/openapi"
	"github.com/stretchr/testify/assert"
)

func TestSpecUIOpts(t *testing.T) {
	tests := []struct {
		name      string
		cfg       *openapi.Config
		specuiCfg *config.SpecUI
	}{
		{
			name: "Swagger UI",
			cfg: &openapi.Config{
				UIProvider: openapi.UIProviderSwaggerUI,
				SwaggerUIConfig: &openapi.SwaggerUIConfig{
					ShowTopBar: true,
					HideCurl:   false,
					JsonEditor: true,
					PreAuthorizeApiKey: map[string]string{
						"api_key": "your_api_key",
					},
					SettingsUI: map[string]string{
						"theme":  "dark",
						"layout": "vertical",
					},
				},
			},
			specuiCfg: &config.SpecUI{
				Provider: config.ProviderSwaggerUI,
				SwaggerUI: config.SwaggerUI{
					ShowTopBar: true,
					HideCurl:   false,
					JsonEditor: true,
					PreAuthorizeApiKey: map[string]string{
						"api_key": "your_api_key",
					},
					SettingsUI: map[string]string{
						"theme":  "dark",
						"layout": "vertical",
					},
				},
			},
		},
		{
			name: "Stoplight Elements",
			cfg: &openapi.Config{
				UIProvider: openapi.UIProviderStoplightElements,
				StoplightElementsConfig: &openapi.StoplightElementsConfig{
					HideExport:  true,
					HideSchemas: false,
					HideTryIt:   true,
					Layout:      "sidebar",
					Logo:        "https://example.com/logo.png",
					Router:      "hash",
				},
			},
			specuiCfg: &config.SpecUI{
				Provider: config.ProviderStoplightElements,
				StoplightElements: config.StoplightElements{
					HideExport:  true,
					HideSchemas: false,
					HideTryIt:   true,
					Layout:      "sidebar",
					Logo:        "https://example.com/logo.png",
					Router:      "hash",
				},
			},
		},
		{
			name: "Redoc",
			cfg: &openapi.Config{
				UIProvider: openapi.UIProviderRedoc,
				RedocConfig: &openapi.RedocConfig{
					HideDownload:     true,
					DisableSearch:    true,
					HideSchemaTitles: true,
				},
			},
			specuiCfg: &config.SpecUI{
				Provider: config.ProviderRedoc,
				Redoc: config.Redoc{
					HideDownload:     true,
					DisableSearch:    true,
					HideSchemaTitles: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := tt.cfg
			specuiCfg := tt.specuiCfg
			opts := mapper.SpecUIOpts(cfg)
			for _, opt := range opts {
				opt(specuiCfg)
			}
			assert.Equal(t, specuiCfg, tt.specuiCfg)
		})
	}
}
