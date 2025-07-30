package spec

import (
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi31"
)

func newReflector31(cfg *Config, jsonSchemaOpts []func(*jsonschema.ReflectContext)) Reflector {
	reflector := openapi31.NewReflector()
	spec := reflector.Spec
	spec.Info.Title = cfg.Title
	spec.Info.Version = cfg.Version
	spec.Info.Description = cfg.Description
	spec.Servers = mapperServers31(cfg.Servers)

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

	return &reflector31{reflector: reflector}
}

type reflector31 struct {
	reflector *openapi31.Reflector
}

func (r *reflector31) AddOperation(oc OperationContext) error {
	return r.reflector.AddOperation(oc.OpenAPIOperationContext())
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

func mapperServers31(servers []Server) []openapi31.Server {
	result := make([]openapi31.Server, 0, len(servers))
	for _, server := range servers {
		result = append(result, mapperServer31(server))
	}
	return result
}

func mapperServer31(server Server) openapi31.Server {
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

func mapperSecurityScheme31(scheme *SecurityScheme) *openapi31.SecurityScheme {
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

func mapperAPIKey31(apiKey *SecuritySchemeAPIKey) *openapi31.SecuritySchemeAPIKey {
	return &openapi31.SecuritySchemeAPIKey{
		Name: apiKey.Name,
		In:   openapi31.SecuritySchemeAPIKeyIn(apiKey.In),
	}
}

func mapperHTTPBearer31(scheme *SecuritySchemeHTTPBearer) *openapi31.SecuritySchemeHTTPBearer {
	return &openapi31.SecuritySchemeHTTPBearer{
		Scheme:       scheme.Scheme,
		BearerFormat: scheme.BearerFormat,
	}
}

func mapperSecuritySchemeOauth2(oauth2 *SecuritySchemeOAuth2) *openapi31.SecuritySchemeOauth2 {
	return &openapi31.SecuritySchemeOauth2{
		Flows: mapperOauth2Flows31(oauth2.Flows),
	}
}

func mapperOauth2Flows31(flows OAuthFlows) openapi31.OauthFlows {
	return openapi31.OauthFlows{
		Implicit:          mapperOauthFlowsDefsImplicit31(flows.Implicit),
		Password:          mapperOauthFlowsDefsPassword31(flows.Password),
		ClientCredentials: mapperOauthFlowsDefsClientCredentials31(flows.ClientCredentials),
		AuthorizationCode: mapperOauthFlowsDefsAuthorizationCode31(flows.AuthorizationCode),
	}
}

func mapperOauthFlowsDefsImplicit31(flows *OAuthFlowsDefsImplicit) *openapi31.OauthFlowsDefsImplicit {
	return &openapi31.OauthFlowsDefsImplicit{
		AuthorizationURL: flows.AuthorizationURL,
		RefreshURL:       flows.RefreshURL,
		Scopes:           flows.Scopes,
		MapOfAnything:    flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsPassword31(flows *OAuthFlowsDefsPassword) *openapi31.OauthFlowsDefsPassword {
	return &openapi31.OauthFlowsDefsPassword{
		TokenURL:      flows.TokenURL,
		RefreshURL:    flows.RefreshURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsClientCredentials31(flows *OAuthFlowsDefsClientCredentials) *openapi31.OauthFlowsDefsClientCredentials {
	return &openapi31.OauthFlowsDefsClientCredentials{
		TokenURL:      flows.TokenURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsAuthorizationCode31(flows *OAuthFlowsDefsAuthorizationCode) *openapi31.OauthFlowsDefsAuthorizationCode {
	return &openapi31.OauthFlowsDefsAuthorizationCode{
		AuthorizationURL: flows.AuthorizationURL,
		TokenURL:         flows.TokenURL,
		RefreshURL:       flows.RefreshURL,
		Scopes:           flows.Scopes,
		MapOfAnything:    flows.MapOfAnything,
	}
}
