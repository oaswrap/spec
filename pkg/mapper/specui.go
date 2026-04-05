package mapper

import (
	"github.com/oaswrap/spec"
	specui "github.com/oaswrap/spec-ui"
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

	if cfg.UIOption != nil {
		opts = append(opts, cfg.UIOption)
	}

	return opts
}
