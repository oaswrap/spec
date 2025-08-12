package mapper_test

import (
	"testing"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec-ui/config"
	"github.com/oaswrap/spec/option"
	"github.com/oaswrap/spec/pkg/mapper"
	"github.com/stretchr/testify/assert"
)

func TestSpecUIOpts(t *testing.T) {
	tests := []struct {
		name   string
		gen    spec.Generator
		assert func(t *testing.T, cfg *config.SpecUI)
	}{
		{
			name: "Swagger UI",
			gen: spec.NewGenerator(
				option.WithSwaggerUI(),
			),
			assert: func(t *testing.T, cfg *config.SpecUI) {
				assert.Equal(t, config.ProviderSwaggerUI, cfg.Provider)
				assert.NotNil(t, cfg.SwaggerUI)
			},
		},
		{
			name: "Stoplight Elements",
			gen: spec.NewGenerator(
				option.WithStoplightElements(),
			),
			assert: func(t *testing.T, cfg *config.SpecUI) {
				assert.Equal(t, config.ProviderStoplightElements, cfg.Provider)
				assert.NotNil(t, cfg.StoplightElements)
			},
		},
		{
			name: "ReDoc",
			gen: spec.NewGenerator(
				option.WithReDoc(),
			),
			assert: func(t *testing.T, cfg *config.SpecUI) {
				assert.Equal(t, config.ProviderReDoc, cfg.Provider)
				assert.NotNil(t, cfg.ReDoc)
			},
		},
		{
			name: "Scalar",
			gen: spec.NewGenerator(
				option.WithScalar(),
			),
			assert: func(t *testing.T, cfg *config.SpecUI) {
				assert.Equal(t, config.ProviderScalar, cfg.Provider)
				assert.NotNil(t, cfg.Scalar)
			},
		},
		{
			name: "RapiDoc",
			gen: spec.NewGenerator(
				option.WithRapiDoc(),
			),
			assert: func(t *testing.T, cfg *config.SpecUI) {
				assert.Equal(t, config.ProviderRapiDoc, cfg.Provider)
				assert.NotNil(t, cfg.RapiDoc)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := mapper.SpecUIOpts(tt.gen)
			cfg := &config.SpecUI{}
			for _, opt := range opts {
				opt(cfg)
			}
			tt.assert(t, cfg)
		})
	}
}
