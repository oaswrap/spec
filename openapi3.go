package spec

import (
	"github.com/oaswrap/spec/internal/mapper"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
)

func newReflector3(cfg *openapi.Config, jsonSchemaOpts []func(*jsonschema.ReflectContext)) reflector {
	reflector := openapi3.NewReflector()
	spec := reflector.Spec
	spec.Info.Title = cfg.Title
	spec.Info.Version = cfg.Version
	spec.Info.Description = cfg.Description
	spec.Info.Contact = mapper.OAS3Contact(cfg.Contact)
	spec.Info.License = mapper.OAS3License(cfg.License)

	spec.ExternalDocs = mapper.OAS3ExternalDocs(cfg.ExternalDocs)
	spec.Servers = mapper.OAS3Servers(cfg.Servers)
	spec.Tags = mapper.OAS3Tags(cfg.Tags)

	if len(cfg.SecuritySchemes) > 0 {
		spec.Components = &openapi3.Components{}
		securitySchemes := &openapi3.ComponentsSecuritySchemes{
			MapOfSecuritySchemeOrRefValues: make(map[string]openapi3.SecuritySchemeOrRef),
		}
		for name, scheme := range cfg.SecuritySchemes {
			openapiScheme := mapper.OAS3SecurityScheme(scheme)
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
	reflector.DefaultOptions = append(reflector.DefaultOptions, jsonSchemaOpts...)

	for _, opt := range cfg.TypeMappings {
		reflector.AddTypeMapping(opt.Src, opt.Dst)
	}

	errors := &SpecError{}

	return &reflector3{reflector: reflector, logger: cfg.Logger, errors: errors}
}

type reflector3 struct {
	logger    openapi.Logger
	errors    *SpecError
	reflector *openapi3.Reflector
}

func (r *reflector3) Spec() spec {
	return r.reflector.Spec
}

func (r *reflector3) Add(method, path string, opts ...option.OperationOption) {
	op, err := r.newOperationContext(method, path)
	if err != nil {
		r.errors.add(err)
		return
	}

	op.With(opts...)

	if err := r.addOperation(op); err != nil {
		r.errors.add(err)
		return
	}
}

func (r *reflector3) Validate() error {
	if r.errors.HasErrors() {
		return r.errors
	}
	return nil
}

func (r *reflector3) addOperation(oc operationContext) error {
	if oc == nil {
		return nil
	}
	openapiOC := oc.build()
	if openapiOC == nil {
		return nil
	}
	return r.reflector.AddOperation(openapiOC)
}

func (r *reflector3) newOperationContext(method, path string) (operationContext, error) {
	op, err := r.reflector.NewOperationContext(method, path)
	if err != nil {
		return nil, err
	}
	return &operationContextImpl{
		op:     op,
		logger: r.logger,
		cfg:    &option.OperationConfig{},
	}, nil
}
