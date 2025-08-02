package mapper

import (
	"github.com/oaswrap/spec/openapi"
	"github.com/swaggest/openapi-go/openapi3"
)

func OAS3Contact(contact *openapi.Contact) *openapi3.Contact {
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

func OAS3License(license *openapi.License) *openapi3.License {
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

func OAS3ExternalDocs(docs *openapi.ExternalDocs) *openapi3.ExternalDocumentation {
	if docs == nil {
		return nil
	}
	result := &openapi3.ExternalDocumentation{
		URL:           docs.URL,
		MapOfAnything: docs.MapOfAnything,
	}
	if docs.Description != "" {
		result.Description = &docs.Description
	}
	return result
}

func OAS3Tags(tags []openapi.Tag) []openapi3.Tag {
	if len(tags) == 0 {
		return nil
	}
	result := make([]openapi3.Tag, 0, len(tags))
	for _, tag := range tags {
		result = append(result, OAS3Tag(tag))
	}
	return result
}

func OAS3Tag(tag openapi.Tag) openapi3.Tag {
	result := openapi3.Tag{
		Name:          tag.Name,
		MapOfAnything: tag.MapOfAnything,
	}
	if tag.Description != "" {
		result.Description = &tag.Description
	}
	if tag.ExternalDocs != nil {
		result.ExternalDocs = OAS3ExternalDocs(tag.ExternalDocs)
	}
	return result
}

func OAS3Servers(servers []openapi.Server) []openapi3.Server {
	if len(servers) == 0 {
		return nil
	}
	result := make([]openapi3.Server, 0, len(servers))
	for _, server := range servers {
		result = append(result, OAS3Server(server))
	}
	return result
}

func OAS3Server(server openapi.Server) openapi3.Server {
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

func OAS3SecurityScheme(scheme *openapi.SecurityScheme) *openapi3.SecurityScheme {
	if scheme == nil {
		return nil
	}
	oasSecurityScheme := &openapi3.SecurityScheme{
		APIKeySecurityScheme: OAS3APIKey(scheme, scheme.APIKey),
		HTTPSecurityScheme:   OAS3HTTPBearer(scheme.HTTPBearer, scheme.Description),
		OAuth2SecurityScheme: OAS3OAuth2SecurityScheme(scheme.OAuth2, scheme.Description),
	}
	if oasSecurityScheme.APIKeySecurityScheme == nil &&
		oasSecurityScheme.HTTPSecurityScheme == nil &&
		oasSecurityScheme.OAuth2SecurityScheme == nil {
		return nil // No valid security scheme defined
	}
	return oasSecurityScheme
}

func OAS3APIKey(scheme *openapi.SecurityScheme, apiKey *openapi.SecuritySchemeAPIKey) *openapi3.APIKeySecurityScheme {
	if apiKey == nil {
		return nil
	}
	return &openapi3.APIKeySecurityScheme{
		Description: scheme.Description,
		Name:        apiKey.Name,
		In:          openapi3.APIKeySecuritySchemeIn(apiKey.In),
	}
}

func OAS3HTTPBearer(securityScheme *openapi.SecuritySchemeHTTPBearer, description *string) *openapi3.HTTPSecurityScheme {
	if securityScheme == nil {
		return nil
	}
	return &openapi3.HTTPSecurityScheme{
		Description:  description,
		Scheme:       securityScheme.Scheme,
		BearerFormat: securityScheme.BearerFormat,
	}
}

func OAS3OAuth2SecurityScheme(oauth2 *openapi.SecuritySchemeOAuth2, description *string) *openapi3.OAuth2SecurityScheme {
	if oauth2 == nil {
		return nil
	}
	return &openapi3.OAuth2SecurityScheme{
		Description: description,
		Flows:       OAS3Oauth2Flows(oauth2.Flows),
	}
}

func OAS3Oauth2Flows(flows openapi.OAuthFlows) openapi3.OAuthFlows {
	return openapi3.OAuthFlows{
		Implicit:          OAS3OauthFlowsImplicit(flows.Implicit),
		Password:          OAS3OauthFlowsPassword(flows.Password),
		ClientCredentials: OAS3OauthFlowsClientCredentials(flows.ClientCredentials),
		AuthorizationCode: OAS3OauthFlowsAuthorizationCode(flows.AuthorizationCode),
	}
}

func OAS3OauthFlowsImplicit(flows *openapi.OAuthFlowsImplicit) *openapi3.ImplicitOAuthFlow {
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

func OAS3OauthFlowsPassword(flows *openapi.OAuthFlowsPassword) *openapi3.PasswordOAuthFlow {
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

func OAS3OauthFlowsClientCredentials(flows *openapi.OAuthFlowsClientCredentials) *openapi3.ClientCredentialsFlow {
	if flows == nil {
		return nil
	}
	return &openapi3.ClientCredentialsFlow{
		TokenURL:      flows.TokenURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func OAS3OauthFlowsAuthorizationCode(flows *openapi.OAuthFlowsAuthorizationCode) *openapi3.AuthorizationCodeOAuthFlow {
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
