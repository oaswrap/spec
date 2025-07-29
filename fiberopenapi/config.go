package fiberopenapi

import (
	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/faizlabs/openapi-wrapper/option"
)

func newConfig(opts ...option.OpenAPIOption) *openapiwrapper.Config {
	cfg := &openapiwrapper.Config{
		OpenAPIVersion:  "3.1.0",
		Title:           "Fiber OpenAPI",
		Description:     nil,
		DisableOpenAPI:  false,
		DocsPath:        "/docs",
		SwaggerConfig:   &openapiwrapper.SwaggerConfig{},
		SecuritySchemes: make(map[string]*openapiwrapper.SecurityScheme),
		Logger:          &noopLogger{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

type noopLogger struct{}

func (noopLogger) Printf(format string, v ...any) {}
