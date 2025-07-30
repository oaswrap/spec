package spec

import (
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi31"
)

func newReflector31(cfg *option.OpenAPI, jsonSchemaOpts []func(*jsonschema.ReflectContext)) Reflector {
	reflector := openapi31.NewReflector()
	spec := reflector.Spec
	spec.Info.Title = cfg.Title
	spec.Info.Version = cfg.Version
	spec.Info.Description = cfg.Description
	spec.Info.Contact = mapperContact31(cfg.Contact)
	spec.Info.License = mapperLicense31(cfg.License)

	spec.ExternalDocs = mapperExternalDocs31(cfg.ExternalDocs)
	spec.Servers = mapperServers31(cfg.Servers)
	spec.Tags = mapperTags31(cfg.Tags)

	if len(cfg.SecuritySchemes) > 0 {
		spec.Components = &openapi31.Components{}
		securitySchemes := make(map[string]openapi31.SecuritySchemeOrReference)
		for name, scheme := range cfg.SecuritySchemes {
			openapiScheme := mapperSecurityScheme31(scheme)
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
	logger    option.Logger
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

func (r *reflector31) Spec() Spec {
	return r.reflector.Spec
}

func (r *reflector31) Validate() error {
	if r.errors.HasErrors() {
		return r.errors
	}
	return nil
}

func (r *reflector31) addOperation(oc OperationContext) error {
	if oc == nil {
		return nil
	}
	openapiOC := oc.build()
	if openapiOC == nil {
		return nil
	}
	return r.reflector.AddOperation(openapiOC)
}

func (r *reflector31) newOperationContext(method, path string) (OperationContext, error) {
	op, err := r.reflector.NewOperationContext(method, path)
	if err != nil {
		return nil, err
	}
	return &operationContext{
		op:     op,
		logger: r.logger,
		cfg:    &option.OperationConfig{},
	}, nil
}

func mapperContact31(contact *option.Contact) *openapi31.Contact {
	if contact == nil {
		return nil
	}
	result := &openapi31.Contact{
		MapOfAnything: contact.MapOfAnything,
	}
	if contact.Name != "" {
		result.Name = &contact.Name
	}
	if contact.URL != "" {
		result.URL = &contact.URL
	}
	if contact.Email != "" {
		result.Email = &contact.Email
	}
	return result
}

func mapperLicense31(license *option.License) *openapi31.License {
	if license == nil {
		return nil
	}
	result := &openapi31.License{
		Name:          license.Name,
		MapOfAnything: license.MapOfAnything,
	}
	if license.URL != "" {
		result.URL = &license.URL
	}
	return result
}

func mapperExternalDocs31(externalDocs *option.ExternalDocs) *openapi31.ExternalDocumentation {
	if externalDocs == nil {
		return nil
	}
	result := &openapi31.ExternalDocumentation{
		URL:           externalDocs.URL,
		MapOfAnything: externalDocs.MapOfAnything,
	}
	if externalDocs.Description != "" {
		result.Description = &externalDocs.Description
	}
	return result
}

func mapperTags31(tags []option.Tag) []openapi31.Tag {
	result := make([]openapi31.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, mapperTag31(tag))
	}
	return result
}

func mapperTag31(tag option.Tag) openapi31.Tag {
	result := openapi31.Tag{
		Name:          tag.Name,
		MapOfAnything: tag.MapOfAnything,
	}
	if tag.Description != "" {
		result.Description = &tag.Description
	}
	if tag.ExternalDocs != nil {
		result.ExternalDocs = mapperExternalDocs31(tag.ExternalDocs)
	}
	return result
}

func mapperServers31(servers []option.Server) []openapi31.Server {
	result := make([]openapi31.Server, 0, len(servers))
	for _, server := range servers {
		result = append(result, mapperServer31(server))
	}
	return result
}

func mapperServer31(server option.Server) openapi31.Server {
	var variables map[string]openapi31.ServerVariable

	if len(server.Variables) > 0 {
		variables = make(map[string]openapi31.ServerVariable, len(server.Variables))
		for name, variable := range server.Variables {
			oasServerVariable := openapi31.ServerVariable{
				Default:       variable.Default,
				Enum:          variable.Enum,
				MapOfAnything: variable.MapOfAnything,
			}
			if variable.Description != "" {
				oasServerVariable.Description = &variable.Description
			}
			variables[name] = oasServerVariable
		}
	}

	return openapi31.Server{
		URL:         server.URL,
		Description: server.Description,
		Variables:   variables,
	}
}

func mapperSecurityScheme31(scheme *option.SecurityScheme) *openapi31.SecurityScheme {
	if scheme == nil {
		return nil
	}
	openapiScheme := &openapi31.SecurityScheme{
		Description:   scheme.Description,
		MapOfAnything: scheme.MapOfAnything,
		APIKey:        mapperAPIKey31(scheme.APIKey),
		HTTPBearer:    mapperHTTPBearer31(scheme.HTTPBearer),
		Oauth2:        mapperSecuritySchemeOauth2(scheme.OAuth2),
	}
	if openapiScheme.APIKey == nil && openapiScheme.HTTPBearer == nil && openapiScheme.Oauth2 == nil {
		return nil // No valid security scheme defined
	}
	return openapiScheme
}

func mapperAPIKey31(apiKey *option.SecuritySchemeAPIKey) *openapi31.SecuritySchemeAPIKey {
	if apiKey == nil {
		return nil
	}
	return &openapi31.SecuritySchemeAPIKey{
		Name: apiKey.Name,
		In:   openapi31.SecuritySchemeAPIKeyIn(apiKey.In),
	}
}

func mapperHTTPBearer31(scheme *option.SecuritySchemeHTTPBearer) *openapi31.SecuritySchemeHTTPBearer {
	if scheme == nil {
		return nil
	}
	return &openapi31.SecuritySchemeHTTPBearer{
		Scheme:       scheme.Scheme,
		BearerFormat: scheme.BearerFormat,
	}
}

func mapperSecuritySchemeOauth2(oauth2 *option.SecuritySchemeOAuth2) *openapi31.SecuritySchemeOauth2 {
	if oauth2 == nil {
		return nil
	}
	return &openapi31.SecuritySchemeOauth2{
		Flows: mapperOauth2Flows31(oauth2.Flows),
	}
}

func mapperOauth2Flows31(flows option.OAuthFlows) openapi31.OauthFlows {
	return openapi31.OauthFlows{
		Implicit:          mapperOauthFlowsDefsImplicit31(flows.Implicit),
		Password:          mapperOauthFlowsDefsPassword31(flows.Password),
		ClientCredentials: mapperOauthFlowsDefsClientCredentials31(flows.ClientCredentials),
		AuthorizationCode: mapperOauthFlowsDefsAuthorizationCode31(flows.AuthorizationCode),
	}
}

func mapperOauthFlowsDefsImplicit31(flows *option.OAuthFlowsDefsImplicit) *openapi31.OauthFlowsDefsImplicit {
	if flows == nil {
		return nil
	}
	return &openapi31.OauthFlowsDefsImplicit{
		AuthorizationURL: flows.AuthorizationURL,
		RefreshURL:       flows.RefreshURL,
		Scopes:           flows.Scopes,
		MapOfAnything:    flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsPassword31(flows *option.OAuthFlowsDefsPassword) *openapi31.OauthFlowsDefsPassword {
	if flows == nil {
		return nil
	}
	return &openapi31.OauthFlowsDefsPassword{
		TokenURL:      flows.TokenURL,
		RefreshURL:    flows.RefreshURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsClientCredentials31(flows *option.OAuthFlowsDefsClientCredentials) *openapi31.OauthFlowsDefsClientCredentials {
	if flows == nil {
		return nil
	}
	return &openapi31.OauthFlowsDefsClientCredentials{
		TokenURL:      flows.TokenURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsAuthorizationCode31(flows *option.OAuthFlowsDefsAuthorizationCode) *openapi31.OauthFlowsDefsAuthorizationCode {
	if flows == nil {
		return nil
	}
	return &openapi31.OauthFlowsDefsAuthorizationCode{
		AuthorizationURL: flows.AuthorizationURL,
		TokenURL:         flows.TokenURL,
		RefreshURL:       flows.RefreshURL,
		Scopes:           flows.Scopes,
		MapOfAnything:    flows.MapOfAnything,
	}
}
