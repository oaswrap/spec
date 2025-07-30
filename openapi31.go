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

	return &reflector31{reflector: reflector}
}

type reflector31 struct {
	reflector *openapi31.Reflector
}

func (r *reflector31) AddOperation(oc OperationContext) error {
	return r.reflector.AddOperation(oc.unwrap())
}

func (r *reflector31) NewOperationContext(method, path string) (OperationContext, error) {
	op, err := r.reflector.NewOperationContext(method, path)
	if err != nil {
		return nil, err
	}
	return &operationContext{OperationContext: op}, nil
}

func (r *reflector31) Spec() Spec {
	return r.reflector.Spec
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

func mapperExternalDocs31(externalDocs *option.ExternalDocumentation) *openapi31.ExternalDocumentation {
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
			variables[name] = openapi31.ServerVariable{
				Default:     variable.Default,
				Description: variable.Description,
				Enum:        variable.Enum,
			}
		}
	}

	return openapi31.Server{
		URL:         server.URL,
		Description: server.Description,
		Variables:   variables,
	}
}

func mapperSecurityScheme31(scheme *option.SecurityScheme) *openapi31.SecurityScheme {
	openapiScheme := &openapi31.SecurityScheme{
		Description:   scheme.Description,
		MapOfAnything: scheme.MapOfAnything,
	}
	if scheme.APIKey != nil {
		openapiScheme.APIKey = mapperAPIKey31(scheme.APIKey)
	} else if scheme.HTTPBearer != nil {
		openapiScheme.HTTPBearer = mapperHTTPBearer31(scheme.HTTPBearer)
	} else if scheme.OAuth2 != nil {
		openapiScheme.Oauth2 = mapperSecuritySchemeOauth2(scheme.OAuth2)
	} else {
		return nil // No valid security scheme found
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
