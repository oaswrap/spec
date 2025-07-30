package mapper

import (
	"github.com/oaswrap/spec/openapi"
	"github.com/swaggest/openapi-go/openapi31"
)

func OAS31Contact(contact *openapi.Contact) *openapi31.Contact {
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

func OAS31License(license *openapi.License) *openapi31.License {
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

func OAS31ExternalDocs(externalDocs *openapi.ExternalDocs) *openapi31.ExternalDocumentation {
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

func OAS31Tags(tags []openapi.Tag) []openapi31.Tag {
	result := make([]openapi31.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, OAS31Tag(tag))
	}
	return result
}

func OAS31Tag(tag openapi.Tag) openapi31.Tag {
	result := openapi31.Tag{
		Name:          tag.Name,
		MapOfAnything: tag.MapOfAnything,
	}
	if tag.Description != "" {
		result.Description = &tag.Description
	}
	if tag.ExternalDocs != nil {
		result.ExternalDocs = OAS31ExternalDocs(tag.ExternalDocs)
	}
	return result
}

func OAS31Servers(servers []openapi.Server) []openapi31.Server {
	result := make([]openapi31.Server, 0, len(servers))
	for _, server := range servers {
		result = append(result, OAS31Server(server))
	}
	return result
}

func OAS31Server(server openapi.Server) openapi31.Server {
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

func OAS31SecurityScheme(scheme *openapi.SecurityScheme) *openapi31.SecurityScheme {
	if scheme == nil {
		return nil
	}
	openapiScheme := &openapi31.SecurityScheme{
		Description:   scheme.Description,
		MapOfAnything: scheme.MapOfAnything,
		APIKey:        OAS31APIKey(scheme.APIKey),
		HTTPBearer:    OAS31HTTPBearer(scheme.HTTPBearer),
		Oauth2:        OAS31SecuritySchemeOauth2(scheme.OAuth2),
	}
	if openapiScheme.APIKey == nil && openapiScheme.HTTPBearer == nil && openapiScheme.Oauth2 == nil {
		return nil // No valid security scheme defined
	}
	return openapiScheme
}

func OAS31APIKey(apiKey *openapi.SecuritySchemeAPIKey) *openapi31.SecuritySchemeAPIKey {
	if apiKey == nil {
		return nil
	}
	return &openapi31.SecuritySchemeAPIKey{
		Name: apiKey.Name,
		In:   openapi31.SecuritySchemeAPIKeyIn(apiKey.In),
	}
}

func OAS31HTTPBearer(scheme *openapi.SecuritySchemeHTTPBearer) *openapi31.SecuritySchemeHTTPBearer {
	if scheme == nil {
		return nil
	}
	return &openapi31.SecuritySchemeHTTPBearer{
		Scheme:       scheme.Scheme,
		BearerFormat: scheme.BearerFormat,
	}
}

func OAS31SecuritySchemeOauth2(oauth2 *openapi.SecuritySchemeOAuth2) *openapi31.SecuritySchemeOauth2 {
	if oauth2 == nil {
		return nil
	}
	return &openapi31.SecuritySchemeOauth2{
		Flows: OAS31Oauth2Flows(oauth2.Flows),
	}
}

func OAS31Oauth2Flows(flows openapi.OAuthFlows) openapi31.OauthFlows {
	return openapi31.OauthFlows{
		Implicit:          OAS31OauthFlowsDefsImplicit(flows.Implicit),
		Password:          OAS31OauthFlowsDefsPassword(flows.Password),
		ClientCredentials: OAS31OauthFlowsDefsClientCredentials(flows.ClientCredentials),
		AuthorizationCode: OAS31OauthFlowsDefsAuthorizationCode(flows.AuthorizationCode),
	}
}

func OAS31OauthFlowsDefsImplicit(flows *openapi.OAuthFlowsDefsImplicit) *openapi31.OauthFlowsDefsImplicit {
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

func OAS31OauthFlowsDefsPassword(flows *openapi.OAuthFlowsDefsPassword) *openapi31.OauthFlowsDefsPassword {
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

func OAS31OauthFlowsDefsClientCredentials(flows *openapi.OAuthFlowsDefsClientCredentials) *openapi31.OauthFlowsDefsClientCredentials {
	if flows == nil {
		return nil
	}
	return &openapi31.OauthFlowsDefsClientCredentials{
		TokenURL:      flows.TokenURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func OAS31OauthFlowsDefsAuthorizationCode(flows *openapi.OAuthFlowsDefsAuthorizationCode) *openapi31.OauthFlowsDefsAuthorizationCode {
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
