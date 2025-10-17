package spec

import (
	"fmt"
	"strings"

	"github.com/oaswrap/spec/internal/debuglog"
	"github.com/oaswrap/spec/internal/errs"
	"github.com/oaswrap/spec/internal/mapper"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/openapi-go/openapi3"
)

type reflector3 struct {
	reflector  *openapi3.Reflector
	logger     *debuglog.Logger
	errors     *errs.SpecError
	pathParser openapi.PathParser
}

func newReflector3(cfg *openapi.Config, logger *debuglog.Logger) reflector {
	reflector := openapi3.NewReflector()
	logger.LogAction("Using OpenAPI 3.0 reflector for version", cfg.OpenAPIVersion)
	spec := reflector.Spec

	spec.Info.Title = cfg.Title
	logger.LogAction("set title", cfg.Title)

	spec.Info.Version = cfg.Version
	logger.LogAction("set version", cfg.Version)

	spec.Info.Description = cfg.Description
	if cfg.Description != nil {
		logger.LogAction("set description", *cfg.Description)
	}

	spec.Info.TermsOfService = cfg.TermsOfService
	if cfg.TermsOfService != nil {
		logger.LogAction("set terms of service", *cfg.TermsOfService)
	}

	spec.Info.Contact = mapper.OAS3Contact(cfg.Contact)
	if cfg.Contact != nil {
		logger.LogContact(cfg.Contact)
	}

	spec.Info.License = mapper.OAS3License(cfg.License)
	if cfg.License != nil {
		logger.LogLicense(cfg.License)
	}

	spec.ExternalDocs = mapper.OAS3ExternalDocs(cfg.ExternalDocs)
	if cfg.ExternalDocs != nil {
		logger.LogExternalDocs(cfg.ExternalDocs)
	}

	spec.Servers = mapper.OAS3Servers(cfg.Servers)
	for _, server := range cfg.Servers {
		logger.LogServer(server)
	}

	spec.Tags = mapper.OAS3Tags(cfg.Tags)
	for _, tag := range cfg.Tags {
		logger.LogTag(tag)
	}

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

	return &reflector3{
		reflector:  reflector,
		logger:     logger,
		errors:     &errs.SpecError{},
		pathParser: cfg.PathParser,
	}
}

func (r *reflector3) Spec() spec {
	return r.reflector.Spec
}

func (r *reflector3) Add(method, path string, opts ...option.OperationOption) {
	if r.pathParser != nil {
		parsedPath, err := r.pathParser.Parse(path)
		if err != nil {
			r.errors.Add(fmt.Errorf("failed to parse path %q: %w", path, err))
			return
		}
		path = parsedPath
	}
	op, err := r.newOperationContext(method, path)
	if err != nil {
		r.errors.Add(err)
		return
	}

	op.With(opts...)

	method = strings.ToUpper(method)

	if err = r.addOperation(op); err != nil {
		r.logger.LogOp(method, path, "add operation", "failed")
		r.errors.Add(err)
		return
	}
	r.logger.LogOp(method, path, "add operation", "successfully registered")
}

func (r *reflector3) Validate() error {
	if r.errors.HasErrors() {
		return r.errors
	}
	return nil
}

func (r *reflector3) addOperation(oc operationContext) error {
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
