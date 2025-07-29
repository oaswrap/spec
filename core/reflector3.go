package core

import (
	"strings"

	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/faizlabs/openapi-wrapper/internal/mapper"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
)

func newReflector3(cfg *openapiwrapper.Config) Reflector {
	reflector := openapi3.NewReflector()
	spec := reflector.Spec
	spec.Info.Title = cfg.Title
	spec.Info.Version = cfg.Version
	spec.Info.Description = cfg.Description
	spec.Servers = mapper.Servers3(cfg.Servers)

	if len(cfg.SecuritySchemes) > 0 {
		spec.Components = &openapi3.Components{}
		securitySchemes := &openapi3.ComponentsSecuritySchemes{
			MapOfSecuritySchemeOrRefValues: make(map[string]openapi3.SecuritySchemeOrRef),
		}
		for name, scheme := range cfg.SecuritySchemes {
			openapiScheme := mapper.SecurityScheme3(scheme)
			if openapiScheme == nil {
				continue // Skip invalid security schemes
			}
			securitySchemes.MapOfSecuritySchemeOrRefValues[name] = openapi3.SecuritySchemeOrRef{
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

	return &reflector3{reflector: reflector}
}

type reflector3 struct {
	reflector *openapi3.Reflector
}

func (r *reflector3) AddOperation(oc openapi.OperationContext) error {
	return r.reflector.AddOperation(oc)
}

func (r *reflector3) NewOperationContext(method, path string) (openapi.OperationContext, error) {
	return r.reflector.NewOperationContext(method, path)
}

func (r *reflector3) Spec() Spec {
	return r.reflector.Spec
}
