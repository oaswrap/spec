package spec

import (
	"fmt"
	"strings"

	"github.com/oaswrap/spec/internal/debuglog"
	"github.com/oaswrap/spec/internal/mapper"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/openapi-go/openapi31"
)

type reflector31 struct {
	reflector  *openapi31.Reflector
	logger     *debuglog.Logger
	pathParser openapi.PathParser
	errors     *SpecError
}

func newReflector31(cfg *openapi.Config, logger *debuglog.Logger) reflector {
	reflector := openapi31.NewReflector()
	logger.LogAction("Using OpenAPI 3.1 reflector for version", cfg.OpenAPIVersion)
	spec := reflector.Spec

	spec.Info.Title = cfg.Title
	logger.LogAction("set title", cfg.Title)

	spec.Info.Version = cfg.Version
	logger.LogAction("set version", cfg.Version)

	spec.Info.Description = cfg.Description
	if cfg.Description != nil {
		logger.LogAction("set description", *cfg.Description)
	}

	spec.Info.Contact = mapper.OAS31Contact(cfg.Contact)
	if cfg.Contact != nil {
		logger.LogContact(cfg.Contact)
	}

	spec.Info.License = mapper.OAS31License(cfg.License)
	if cfg.License != nil {
		logger.LogLicense(cfg.License)
	}

	spec.ExternalDocs = mapper.OAS31ExternalDocs(cfg.ExternalDocs)
	if cfg.ExternalDocs != nil {
		logger.LogExternalDocs(cfg.ExternalDocs)
	}

	spec.Servers = mapper.OAS31Servers(cfg.Servers)
	for _, server := range cfg.Servers {
		logger.LogServer(server)
	}

	spec.Tags = mapper.OAS31Tags(cfg.Tags)
	for _, tag := range cfg.Tags {
		logger.LogTag(tag)
	}

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
		for name, scheme := range cfg.SecuritySchemes {
			logger.LogSecurityScheme(name, scheme)
		}
	}

	// Custom options for JSON schema generation
	if cfg.ReflectorConfig != nil {
		jsonSchemaOpts := getJSONSchemaOpts(cfg.ReflectorConfig, logger)
		if len(jsonSchemaOpts) > 0 {
			reflector.DefaultOptions = append(reflector.DefaultOptions, jsonSchemaOpts...)
		}

		for _, opt := range cfg.ReflectorConfig.TypeMappings {
			reflector.AddTypeMapping(opt.Src, opt.Dst)
			logger.LogAction("add type mapping", fmt.Sprintf("%T -> %T", opt.Src, opt.Dst))
		}
	}

	return &reflector31{
		reflector:  reflector,
		logger:     logger,
		errors:     &SpecError{},
		pathParser: cfg.PathParser,
	}
}

func (r *reflector31) Add(method, path string, opts ...option.OperationOption) {
	if r.pathParser != nil {
		parsedPath, err := r.pathParser.Parse(path)
		if err != nil {
			r.errors.add(err)
			return
		}
		path = parsedPath
	}
	op, err := r.newOperationContext(method, path)
	if err != nil {
		r.errors.add(err)
		return
	}

	op.With(opts...)

	method = strings.ToUpper(method)

	if err := r.addOperation(op); err != nil {
		r.logger.LogOp(method, path, "add operation", "failed")
		r.errors.add(err)
		return
	}
	r.logger.LogOp(method, path, "add operation", "successfully registered")
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
