package mapper

import (
	"github.com/oaswrap/spec"
	specui "github.com/oaswrap/spec-ui"
	"github.com/oaswrap/spec-ui/config"
)

func SpecUIOpts(gen spec.Generator) []specui.Option {
	cfg := gen.Config()
	opts := []specui.Option{
		specui.WithTitle(cfg.Title),
		specui.WithDocsPath(cfg.DocsPath),
		specui.WithSpecPath(cfg.SpecPath),
		specui.WithSpecGenerator(gen),
	}
	if cfg.CacheAge != nil {
		opts = append(opts, specui.WithCacheAge(*cfg.CacheAge))
	}

	switch cfg.UIProvider {
	case config.ProviderSwaggerUI:
		opts = append(opts, specui.WithSwaggerUI(*cfg.SwaggerUIConfig))
	case config.ProviderStoplightElements:
		opts = append(opts, specui.WithStoplightElements(*cfg.StoplightElementsConfig))
	case config.ProviderReDoc:
		opts = append(opts, specui.WithReDoc(*cfg.ReDocConfig))
	case config.ProviderScalar:
		opts = append(opts, specui.WithScalar(*cfg.ScalarConfig))
	case config.ProviderRapiDoc:
		opts = append(opts, specui.WithRapiDoc(*cfg.RapiDocConfig))
	}

	return opts
}
