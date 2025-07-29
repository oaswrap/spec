package core

import (
	"strings"

	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/faizlabs/openapi-wrapper/internal/mapper"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi31"
)

func newReflector31(cfg *openapiwrapper.Config) Reflector {
	reflector := openapi31.NewReflector()
	spec := reflector.Spec
	spec.Info.Title = cfg.Title
	spec.Info.Version = cfg.Version
	spec.Info.Description = cfg.Description
	spec.Servers = mapper.Servers31(cfg.Servers)

	if len(cfg.SecuritySchemes) > 0 {
		spec.Components = &openapi31.Components{}
		securitySchemes := make(map[string]openapi31.SecuritySchemeOrReference)
		for name, scheme := range cfg.SecuritySchemes {
			openapiScheme := mapper.SecurityScheme31(scheme)
			if openapiScheme == nil {
				continue // Skip invalid security schemes
			}
			securitySchemes[name] = openapi31.SecuritySchemeOrReference{
				SecurityScheme: openapiScheme,
			}
		}
		spec.Components.SecuritySchemes = securitySchemes
	}

	// Custom options for JSON schema generation
	reflector.DefaultOptions = append(reflector.DefaultOptions,
		jsonschema.InterceptProp(func(params jsonschema.InterceptPropParams) error {
			if !params.Processed {
				return nil
			}
			if v, ok := params.Field.Tag.Lookup("validate"); ok {
				if strings.Contains(v, "required") {
					params.ParentSchema.Required = append(params.ParentSchema.Required, params.Name)
				}
			}
			return nil
		}),
	)

	return &reflector31{reflector: reflector}
}

type reflector31 struct {
	reflector *openapi31.Reflector
}

func (r *reflector31) AddOperation(oc openapi.OperationContext) error {
	return r.reflector.AddOperation(oc)
}

func (r *reflector31) NewOperationContext(method, path string) (openapi.OperationContext, error) {
	return r.reflector.NewOperationContext(method, path)
}

func (r *reflector31) Spec() Spec {
	return r.reflector.Spec
}
