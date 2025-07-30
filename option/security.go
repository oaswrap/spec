package option

import "github.com/oaswrap/spec/openapi"

type securityConfig struct {
	Description *string
	APIKey      *openapi.SecuritySchemeAPIKey
	HTTPBearer  *openapi.SecuritySchemeHTTPBearer
	Oauth2      *openapi.SecuritySchemeOAuth2
}

// SecurityOption is a function that applies configuration to a securityConfig.
type SecurityOption func(*securityConfig)

// SecurityDescription sets the description for the security scheme.
func SecurityDescription(description string) SecurityOption {
	return func(cfg *securityConfig) {
		if description != "" {
			cfg.Description = &description
		} else {
			cfg.Description = nil // Clear description if empty
		}
	}
}

// SecurityAPIKey creates a security scheme for API key authentication.
func SecurityAPIKey(name string, in openapi.SecuritySchemeAPIKeyIn) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.APIKey = &openapi.SecuritySchemeAPIKey{
			Name: name,
			In:   in,
		}
	}
}

// SecurityHTTPBearer creates a security scheme for HTTP Bearer authentication.
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

// SecurityOAuth2 creates a security scheme for OAuth 2.0 authentication.
func SecurityOAuth2(flows openapi.OAuthFlows) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.Oauth2 = &openapi.SecuritySchemeOAuth2{
			Flows: flows,
		}
	}
}
