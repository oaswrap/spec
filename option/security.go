package option

import "github.com/oaswrap/spec/openapi"

// securityConfig holds configuration for defining a security scheme.
type securityConfig struct {
	Description *string
	APIKey      *openapi.SecuritySchemeAPIKey
	HTTPBearer  *openapi.SecuritySchemeHTTPBearer
	Oauth2      *openapi.SecuritySchemeOAuth2
}

// SecurityOption applies configuration to a securityConfig.
type SecurityOption func(*securityConfig)

// SecurityDescription sets the description for the security scheme.
//
// If the description is empty, it clears any existing description.
func SecurityDescription(description string) SecurityOption {
	return func(cfg *securityConfig) {
		if description != "" {
			cfg.Description = &description
		} else {
			cfg.Description = nil
		}
	}
}

// SecurityAPIKey defines an API key security scheme.
//
// Example:
//
//	option.WithSecurity("apiKey",
//	    option.SecurityAPIKey("x-api-key", openapi.SecuritySchemeAPIKeyInHeader),
//	)
func SecurityAPIKey(name string, in openapi.SecuritySchemeAPIKeyIn) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.APIKey = &openapi.SecuritySchemeAPIKey{
			Name: name,
			In:   in,
		}
	}
}

// SecurityHTTPBearer defines an HTTP Bearer security scheme.
//
// Optionally, you can provide a bearer format.
func SecurityHTTPBearer(scheme string, bearerFormat ...string) SecurityOption {
	return func(cfg *securityConfig) {
		httpBearer := &openapi.SecuritySchemeHTTPBearer{
			Scheme: scheme,
		}
		if len(bearerFormat) > 0 {
			httpBearer.BearerFormat = &bearerFormat[0]
		}
		cfg.HTTPBearer = httpBearer
	}
}

// SecurityOAuth2 defines an OAuth 2.0 security scheme.
func SecurityOAuth2(flows openapi.OAuthFlows) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.Oauth2 = &openapi.SecuritySchemeOAuth2{
			Flows: flows,
		}
	}
}
