package spec

import (
	"github.com/oaswrap/spec/internal/mapper"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi31"
)

func newReflector31(cfg *openapi.Config, jsonSchemaOpts []func(*jsonschema.ReflectContext)) reflector {
	reflector := openapi31.NewReflector()
	spec := reflector.Spec
	spec.Info.Title = cfg.Title
	spec.Info.Version = cfg.Version
	spec.Info.Description = cfg.Description
	spec.Info.Contact = mapper.OAS31Contact(cfg.Contact)
	spec.Info.License = mapper.OAS31License(cfg.License)

	spec.ExternalDocs = mapper.OAS31ExternalDocs(cfg.ExternalDocs)
	spec.Servers = mapper.OAS31Servers(cfg.Servers)
	spec.Tags = mapper.OAS31Tags(cfg.Tags)

	if len(cfg.SecuritySchemes) > 0 {
		spec.Components = &openapi31.Components{}
		securitySchemes := make(map[string]openapi31.SecuritySchemeOrReference)
		for name, scheme := range cfg.SecuritySchemes {
			openapiScheme := mapper.OAS31SecurityScheme(scheme)
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
	reflector.DefaultOptions = append(reflector.DefaultOptions, jsonSchemaOpts...)

	for _, opt := range cfg.TypeMappings {
		reflector.AddTypeMapping(opt.Src, opt.Dst)
	}

	errors := &SpecError{}

	return &reflector31{reflector: reflector, logger: cfg.Logger, errors: errors}
}

type reflector31 struct {
	reflector *openapi31.Reflector
	logger    openapi.Logger
	errors    *SpecError
}

func (r *reflector31) Add(method, path string, opts ...option.OperationOption) {
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

func (r *reflector31) Spec() spec {
	return r.reflector.Spec
}

func (r *reflector31) Validate() error {
	if r.errors.HasErrors() {
		return r.errors
	}
	return nil
}

func (r *reflector31) addOperation(oc operationContext) error {
	if oc == nil {
		return nil
	}
	openapiOC := oc.build()
	if openapiOC == nil {
		return nil
	}
	return r.reflector.AddOperation(openapiOC)
}

func (r *reflector31) newOperationContext(method, path string) (operationContext, error) {
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
