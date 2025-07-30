package spec

import (
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
)

func newReflector3(cfg *option.OpenAPI, jsonSchemaOpts []func(*jsonschema.ReflectContext)) Reflector {
	reflector := openapi3.NewReflector()
	spec := reflector.Spec
	spec.Info.Title = cfg.Title
	spec.Info.Version = cfg.Version
	spec.Info.Description = cfg.Description
	spec.Info.Contact = mapperContact3(cfg.Contact)
	spec.Info.License = mapperLicense3(cfg.License)

	spec.ExternalDocs = mapperExternalDocs3(cfg.ExternalDocs)
	spec.Servers = mapperServers3(cfg.Servers)
	spec.Tags = mapperTags3(cfg.Tags)

	if len(cfg.SecuritySchemes) > 0 {
		spec.Components = &openapi3.Components{}
		securitySchemes := &openapi3.ComponentsSecuritySchemes{
			MapOfSecuritySchemeOrRefValues: make(map[string]openapi3.SecuritySchemeOrRef),
		}
		for name, scheme := range cfg.SecuritySchemes {
			openapiScheme := mapperSecurityScheme3(scheme)
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
	logger    option.Logger
	errors    *SpecError
	reflector *openapi3.Reflector
}

func (r *reflector3) Spec() Spec {
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

func (r *reflector3) addOperation(oc OperationContext) error {
	if oc == nil {
		return nil
	}
	openapiOC := oc.build()
	if openapiOC == nil {
		return nil
	}
	return r.reflector.AddOperation(openapiOC)
}

func (r *reflector3) newOperationContext(method, path string) (OperationContext, error) {
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

func mapperContact3(contact *option.Contact) *openapi3.Contact {
	if contact == nil {
		return nil
	}
	result := &openapi3.Contact{
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

func mapperLicense3(license *option.License) *openapi3.License {
	if license == nil {
		return nil
	}
	result := &openapi3.License{
		Name:          license.Name,
		MapOfAnything: license.MapOfAnything,
	}
	if license.URL != "" {
		result.URL = &license.URL
	}
	return result
}

func mapperExternalDocs3(docs *option.ExternalDocs) *openapi3.ExternalDocumentation {
	if docs == nil {
		return nil
	}
	result := &openapi3.ExternalDocumentation{
		URL: docs.URL,
	}
	if docs.Description != "" {
		result.Description = &docs.Description
	}
	return result
}

func mapperTags3(tags []option.Tag) []openapi3.Tag {
	result := make([]openapi3.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, mapperTag3(tag))
	}
	return result
}

func mapperTag3(tag option.Tag) openapi3.Tag {
	result := openapi3.Tag{
		Name: tag.Name,
	}
	if tag.Description != "" {
		result.Description = &tag.Description
	}
	if tag.ExternalDocs != nil {
		result.ExternalDocs = mapperExternalDocs3(tag.ExternalDocs)
	}
	return result
}

func mapperServers3(servers []option.Server) []openapi3.Server {
	result := make([]openapi3.Server, 0, len(servers))
	for _, server := range servers {
		result = append(result, mapperServer3(server))
	}
	return result
}

func mapperServer3(server option.Server) openapi3.Server {
	var variables map[string]openapi3.ServerVariable

	if len(server.Variables) > 0 {
		variables = make(map[string]openapi3.ServerVariable, len(server.Variables))
		for name, variable := range server.Variables {
			oasServerVariable := openapi3.ServerVariable{
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

	return openapi3.Server{
		URL:         server.URL,
		Description: server.Description,
		Variables:   variables,
	}
}

func mapperSecurityScheme3(scheme *option.SecurityScheme) *openapi3.SecurityScheme {
	if scheme == nil {
		return nil
	}
	oasSecurityScheme := &openapi3.SecurityScheme{
		APIKeySecurityScheme: mapperAPIKey3(scheme, scheme.APIKey),
		HTTPSecurityScheme:   mapperHTTPBearer3(scheme, scheme.HTTPBearer),
		OAuth2SecurityScheme: mapperOAuth2SecurityScheme(scheme, scheme.OAuth2),
	}
	if oasSecurityScheme.APIKeySecurityScheme == nil &&
		oasSecurityScheme.HTTPSecurityScheme == nil &&
		oasSecurityScheme.OAuth2SecurityScheme == nil {
		return nil // No valid security scheme defined
	}
	return oasSecurityScheme
}

func mapperAPIKey3(scheme *option.SecurityScheme, apiKey *option.SecuritySchemeAPIKey) *openapi3.APIKeySecurityScheme {
	if apiKey == nil {
		return nil
	}
	return &openapi3.APIKeySecurityScheme{
		Description: scheme.Description,
		Name:        apiKey.Name,
		In:          openapi3.APIKeySecuritySchemeIn(apiKey.In),
	}
}

func mapperHTTPBearer3(scheme *option.SecurityScheme, securityScheme *option.SecuritySchemeHTTPBearer) *openapi3.HTTPSecurityScheme {
	if securityScheme == nil {
		return nil
	}
	return &openapi3.HTTPSecurityScheme{
		Description:  scheme.Description,
		Scheme:       securityScheme.Scheme,
		BearerFormat: securityScheme.BearerFormat,
	}
}

func mapperOAuth2SecurityScheme(scheme *option.SecurityScheme, oauth2 *option.SecuritySchemeOAuth2) *openapi3.OAuth2SecurityScheme {
	if oauth2 == nil {
		return nil
	}
	return &openapi3.OAuth2SecurityScheme{
		Description: scheme.Description,
		Flows:       mapperOauth2Flows3(oauth2.Flows),
	}
}

func mapperOauth2Flows3(flows option.OAuthFlows) openapi3.OAuthFlows {
	return openapi3.OAuthFlows{
		Implicit:          mapperOauthFlowsDefsImplicit3(flows.Implicit),
		Password:          mapperOauthFlowsDefsPassword3(flows.Password),
		ClientCredentials: mapperOauthFlowsDefsClientCredentials3(flows.ClientCredentials),
		AuthorizationCode: mapperOauthFlowsDefsAuthorizationCode3(flows.AuthorizationCode),
	}
}

func mapperOauthFlowsDefsImplicit3(flows *option.OAuthFlowsDefsImplicit) *openapi3.ImplicitOAuthFlow {
	if flows == nil {
		return nil
	}
	return &openapi3.ImplicitOAuthFlow{
		AuthorizationURL: flows.AuthorizationURL,
		RefreshURL:       flows.RefreshURL,
		Scopes:           flows.Scopes,
		MapOfAnything:    flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsPassword3(flows *option.OAuthFlowsDefsPassword) *openapi3.PasswordOAuthFlow {
	if flows == nil {
		return nil
	}
	return &openapi3.PasswordOAuthFlow{
		TokenURL:      flows.TokenURL,
		RefreshURL:    flows.RefreshURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsClientCredentials3(flows *option.OAuthFlowsDefsClientCredentials) *openapi3.ClientCredentialsFlow {
	if flows == nil {
		return nil
	}
	return &openapi3.ClientCredentialsFlow{
		TokenURL:      flows.TokenURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsAuthorizationCode3(flows *option.OAuthFlowsDefsAuthorizationCode) *openapi3.AuthorizationCodeOAuthFlow {
	if flows == nil {
		return nil
	}
	return &openapi3.AuthorizationCodeOAuthFlow{
		AuthorizationURL: flows.AuthorizationURL,
		TokenURL:         flows.TokenURL,
		RefreshURL:       flows.RefreshURL,
		Scopes:           flows.Scopes,
		MapOfAnything:    flows.MapOfAnything,
	}
}
