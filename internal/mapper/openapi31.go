package mapper

import (
	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/swaggest/openapi-go/openapi31"
)

func Servers31(servers []openapiwrapper.Server) []openapi31.Server {
	result := make([]openapi31.Server, 0, len(servers))
	for _, server := range servers {
		result = append(result, Server31(server))
	}
	return result
}

func Server31(server openapiwrapper.Server) openapi31.Server {
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

func SecurityScheme31(scheme *openapiwrapper.SecurityScheme) *openapi31.SecurityScheme {
	if scheme == nil {
		return nil
	}
	openapiScheme := &openapi31.SecurityScheme{
		Description:   scheme.Description,
		MapOfAnything: scheme.MapOfAnything,
	}
	if scheme.APIKey != nil {
		openapiScheme.APIKey = APIKey31(scheme.APIKey)
	} else if scheme.HTTPBearer != nil {
		openapiScheme.HTTPBearer = HTTPBearer31(scheme.HTTPBearer)
	} else if scheme.OAuth2 != nil {
		openapiScheme.Oauth2 = SecuritySchemeOauth2(scheme.OAuth2)
	} else {
		return nil // No valid security scheme found
	}
	return openapiScheme
}

func APIKey31(apiKey *openapiwrapper.SecuritySchemeAPIKey) *openapi31.SecuritySchemeAPIKey {
	if apiKey == nil {
		return nil
	}
	return &openapi31.SecuritySchemeAPIKey{
		Name: apiKey.Name,
		In:   openapi31.SecuritySchemeAPIKeyIn(apiKey.In),
	}
}

func HTTPBearer31(scheme *openapiwrapper.SecuritySchemeHTTPBearer) *openapi31.SecuritySchemeHTTPBearer {
	if scheme == nil {
		return nil
	}
	return &openapi31.SecuritySchemeHTTPBearer{
		Scheme:       scheme.Scheme,
		BearerFormat: scheme.BearerFormat,
	}
}

func SecuritySchemeOauth2(oauth2 *openapiwrapper.SecuritySchemeOAuth2) *openapi31.SecuritySchemeOauth2 {
	if oauth2 == nil {
		return nil
	}
	return &openapi31.SecuritySchemeOauth2{
		Flows: Oauth2Flows31(oauth2.Flows),
	}
}

func Oauth2Flows31(flows openapiwrapper.OAuthFlows) openapi31.OauthFlows {
	return openapi31.OauthFlows{
		Implicit:          OauthFlowsDefsImplicit31(flows.Implicit),
		Password:          OauthFlowsDefsPassword31(flows.Password),
		ClientCredentials: OauthFlowsDefsClientCredentials31(flows.ClientCredentials),
		AuthorizationCode: OauthFlowsDefsAuthorizationCode31(flows.AuthorizationCode),
	}
}

func OauthFlowsDefsImplicit31(flows *openapiwrapper.OAuthFlowsDefsImplicit) *openapi31.OauthFlowsDefsImplicit {
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

func OauthFlowsDefsPassword31(flows *openapiwrapper.OAuthFlowsDefsPassword) *openapi31.OauthFlowsDefsPassword {
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

func OauthFlowsDefsClientCredentials31(flows *openapiwrapper.OAuthFlowsDefsClientCredentials) *openapi31.OauthFlowsDefsClientCredentials {
	if flows == nil {
		return nil
	}
	return &openapi31.OauthFlowsDefsClientCredentials{
		TokenURL:      flows.TokenURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func OauthFlowsDefsAuthorizationCode31(flows *openapiwrapper.OAuthFlowsDefsAuthorizationCode) *openapi31.OauthFlowsDefsAuthorizationCode {
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
