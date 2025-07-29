package mapper

import (
	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/swaggest/openapi-go/openapi3"
)

func Servers3(servers []openapiwrapper.Server) []openapi3.Server {
	result := make([]openapi3.Server, 0, len(servers))
	for _, server := range servers {
		result = append(result, Server3(server))
	}
	return result
}

func Server3(server openapiwrapper.Server) openapi3.Server {
	var variables map[string]openapi3.ServerVariable

	if len(server.Variables) > 0 {
		variables = make(map[string]openapi3.ServerVariable, len(server.Variables))
		for name, variable := range server.Variables {
			variables[name] = openapi3.ServerVariable{
				Default:     variable.Default,
				Description: variable.Description,
				Enum:        variable.Enum,
			}
		}
	}

	return openapi3.Server{
		URL:         server.URL,
		Description: server.Description,
		Variables:   variables,
	}
}

func SecurityScheme3(scheme *openapiwrapper.SecurityScheme) *openapi3.SecurityScheme {
	if scheme == nil {
		return nil
	}
	openapiScheme := &openapi3.SecurityScheme{}
	if scheme.APIKey != nil {
		openapiScheme.APIKeySecurityScheme = APIKey3(scheme.APIKey)
	} else if scheme.HTTPBearer != nil {
		openapiScheme.HTTPSecurityScheme = HTTPBearer3(scheme.HTTPBearer)
	} else if scheme.OAuth2 != nil {
		openapiScheme.OAuth2SecurityScheme = OAuth2SecurityScheme(scheme.OAuth2)
	} else {
		return nil // No valid security scheme found
	}
	return openapiScheme
}

func APIKey3(apiKey *openapiwrapper.SecuritySchemeAPIKey) *openapi3.APIKeySecurityScheme {
	if apiKey == nil {
		return nil
	}
	return &openapi3.APIKeySecurityScheme{
		Name: apiKey.Name,
		In:   openapi3.APIKeySecuritySchemeIn(apiKey.In),
	}
}

func HTTPBearer3(scheme *openapiwrapper.SecuritySchemeHTTPBearer) *openapi3.HTTPSecurityScheme {
	if scheme == nil {
		return nil
	}
	return &openapi3.HTTPSecurityScheme{
		Scheme:       scheme.Scheme,
		BearerFormat: scheme.BearerFormat,
	}
}

func OAuth2SecurityScheme(oauth2 *openapiwrapper.SecuritySchemeOAuth2) *openapi3.OAuth2SecurityScheme {
	if oauth2 == nil {
		return nil
	}
	return &openapi3.OAuth2SecurityScheme{
		Flows: Oauth2Flows3(oauth2.Flows),
	}
}

func Oauth2Flows3(flows openapiwrapper.OAuthFlows) openapi3.OAuthFlows {
	return openapi3.OAuthFlows{
		Implicit:          OauthFlowsDefsImplicit3(flows.Implicit),
		Password:          OauthFlowsDefsPassword3(flows.Password),
		ClientCredentials: OauthFlowsDefsClientCredentials3(flows.ClientCredentials),
		AuthorizationCode: OauthFlowsDefsAuthorizationCode3(flows.AuthorizationCode),
	}
}

func OauthFlowsDefsImplicit3(flows *openapiwrapper.OAuthFlowsDefsImplicit) *openapi3.ImplicitOAuthFlow {
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

func OauthFlowsDefsPassword3(flows *openapiwrapper.OAuthFlowsDefsPassword) *openapi3.PasswordOAuthFlow {
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

func OauthFlowsDefsClientCredentials3(flows *openapiwrapper.OAuthFlowsDefsClientCredentials) *openapi3.ClientCredentialsFlow {
	if flows == nil {
		return nil
	}
	return &openapi3.ClientCredentialsFlow{
		TokenURL:      flows.TokenURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func OauthFlowsDefsAuthorizationCode3(flows *openapiwrapper.OAuthFlowsDefsAuthorizationCode) *openapi3.AuthorizationCodeOAuthFlow {
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
