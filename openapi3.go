package spec

import (
	"github.com/swaggest/jsonschema-go"
	"github.com/swaggest/openapi-go/openapi3"
)

func newReflector3(cfg *Config, jsonSchemaOpts []func(*jsonschema.ReflectContext)) Reflector {
	reflector := openapi3.NewReflector()
	spec := reflector.Spec
	spec.Info.Title = cfg.Title
	spec.Info.Version = cfg.Version
	spec.Info.Description = cfg.Description
	spec.Servers = mapperServers3(cfg.Servers)

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

	return &reflector3{reflector: reflector}
}

type reflector3 struct {
	reflector *openapi3.Reflector
}

func (r *reflector3) AddOperation(oc OperationContext) error {
	return r.reflector.AddOperation(oc.OpenAPIOperationContext())
}

func (r *reflector3) NewOperationContext(method, path string) (OperationContext, error) {
	op, err := r.reflector.NewOperationContext(method, path)
	if err != nil {
		return nil, err
	}
	return &operationContext{OperationContext: op}, nil
}

func (r *reflector3) Spec() Spec {
	return r.reflector.Spec
}

func mapperServers3(servers []Server) []openapi3.Server {
	result := make([]openapi3.Server, 0, len(servers))
	for _, server := range servers {
		result = append(result, mapperServer3(server))
	}
	return result
}

func mapperServer3(server Server) openapi3.Server {
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

func mapperSecurityScheme3(scheme *SecurityScheme) *openapi3.SecurityScheme {
	openapiScheme := &openapi3.SecurityScheme{}
	if scheme.APIKey != nil {
		openapiScheme.APIKeySecurityScheme = mapperAPIKey3(scheme.APIKey)
	} else if scheme.HTTPBearer != nil {
		openapiScheme.HTTPSecurityScheme = mapperHTTPBearer3(scheme.HTTPBearer)
	} else if scheme.OAuth2 != nil {
		openapiScheme.OAuth2SecurityScheme = mapperOAuth2SecurityScheme(scheme.OAuth2)
	} else {
		return nil // No valid security scheme found
	}
	return openapiScheme
}

func mapperAPIKey3(apiKey *SecuritySchemeAPIKey) *openapi3.APIKeySecurityScheme {
	return &openapi3.APIKeySecurityScheme{
		Name: apiKey.Name,
		In:   openapi3.APIKeySecuritySchemeIn(apiKey.In),
	}
}

func mapperHTTPBearer3(scheme *SecuritySchemeHTTPBearer) *openapi3.HTTPSecurityScheme {
	return &openapi3.HTTPSecurityScheme{
		Scheme:       scheme.Scheme,
		BearerFormat: scheme.BearerFormat,
	}
}

func mapperOAuth2SecurityScheme(oauth2 *SecuritySchemeOAuth2) *openapi3.OAuth2SecurityScheme {
	return &openapi3.OAuth2SecurityScheme{
		Flows: mapperOauth2Flows3(oauth2.Flows),
	}
}

func mapperOauth2Flows3(flows OAuthFlows) openapi3.OAuthFlows {
	return openapi3.OAuthFlows{
		Implicit:          mapperOauthFlowsDefsImplicit3(flows.Implicit),
		Password:          mapperOauthFlowsDefsPassword3(flows.Password),
		ClientCredentials: mapperOauthFlowsDefsClientCredentials3(flows.ClientCredentials),
		AuthorizationCode: mapperOauthFlowsDefsAuthorizationCode3(flows.AuthorizationCode),
	}
}

func mapperOauthFlowsDefsImplicit3(flows *OAuthFlowsDefsImplicit) *openapi3.ImplicitOAuthFlow {
	return &openapi3.ImplicitOAuthFlow{
		AuthorizationURL: flows.AuthorizationURL,
		RefreshURL:       flows.RefreshURL,
		Scopes:           flows.Scopes,
		MapOfAnything:    flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsPassword3(flows *OAuthFlowsDefsPassword) *openapi3.PasswordOAuthFlow {
	return &openapi3.PasswordOAuthFlow{
		TokenURL:      flows.TokenURL,
		RefreshURL:    flows.RefreshURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsClientCredentials3(flows *OAuthFlowsDefsClientCredentials) *openapi3.ClientCredentialsFlow {
	return &openapi3.ClientCredentialsFlow{
		TokenURL:      flows.TokenURL,
		Scopes:        flows.Scopes,
		MapOfAnything: flows.MapOfAnything,
	}
}

func mapperOauthFlowsDefsAuthorizationCode3(flows *OAuthFlowsDefsAuthorizationCode) *openapi3.AuthorizationCodeOAuthFlow {
	return &openapi3.AuthorizationCodeOAuthFlow{
		AuthorizationURL: flows.AuthorizationURL,
		TokenURL:         flows.TokenURL,
		RefreshURL:       flows.RefreshURL,
		Scopes:           flows.Scopes,
		MapOfAnything:    flows.MapOfAnything,
	}
}
