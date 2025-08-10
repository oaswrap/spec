package mapper

import (
	specui "github.com/oaswrap/spec-ui"
	"github.com/oaswrap/spec-ui/config"
	"github.com/oaswrap/spec/openapi"
)

func SpecUIOpts(cfg *openapi.Config) []specui.Option {
	opts := []specui.Option{
		specui.WithTitle(cfg.Title),
		specui.WithDocsPath(cfg.DocsPath),
		specui.WithSpecPath(cfg.SpecPath),
	}

	switch cfg.UIProvider {
	case openapi.UIProviderSwaggerUI:
		opts = append(opts, specui.WithSwaggerUI(SwaggerUI(cfg.SwaggerUIConfig)))
	case openapi.UIProviderStoplightElements:
		opts = append(opts, specui.WithStoplightElements(StoplightElements(cfg.StoplightElementsConfig)))
	case openapi.UIProviderRedoc:
		opts = append(opts, specui.WithRedoc(ReDoc(cfg.RedocConfig)))
	}

	return opts
}

func SwaggerUI(cfg *openapi.SwaggerUIConfig) config.SwaggerUI {
	if cfg == nil {
		return config.SwaggerUI{}
	}
	return config.SwaggerUI{
		ShowTopBar:         cfg.ShowTopBar,
		HideCurl:           cfg.HideCurl,
		JsonEditor:         cfg.JsonEditor,
		PreAuthorizeApiKey: cfg.PreAuthorizeApiKey,
		SettingsUI:         cfg.SettingsUI,
	}
}

func StoplightElements(cfg *openapi.StoplightElementsConfig) config.StoplightElements {
	if cfg == nil {
		return config.StoplightElements{}
	}
	return config.StoplightElements{
		HideExport:  cfg.HideExport,
		HideSchemas: cfg.HideSchemas,
		HideTryIt:   cfg.HideTryIt,
		Layout:      cfg.Layout,
		Logo:        cfg.Logo,
		Router:      cfg.Router,
	}
}

func ReDoc(cfg *openapi.RedocConfig) config.Redoc {
	if cfg == nil {
		return config.Redoc{}
	}
	return config.Redoc{
		DisableSearch:    cfg.DisableSearch,
		HideDownload:     cfg.HideDownload,
		HideSchemaTitles: cfg.HideSchemaTitles,
	}
}
