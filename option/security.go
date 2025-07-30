package option

import (
	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/pkg/util"
)

type securityConfig struct {
	Description *string
	APIKey      *spec.SecuritySchemeAPIKey
	HTTPBearer  *spec.SecuritySchemeHTTPBearer
	Oauth2      *spec.SecuritySchemeOAuth2
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
func SecurityAPIKey(name string, in spec.SecuritySchemeAPIKeyIn) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.APIKey = &spec.SecuritySchemeAPIKey{
			Name: name,
			In:   in,
		}
	}
}

// SecurityHTTPBearer creates a security scheme for HTTP Bearer authentication.
func SecurityHTTPBearer(scheme ...string) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.HTTPBearer = &spec.SecuritySchemeHTTPBearer{
			Scheme: util.Optional("Bearer", scheme...),
		}
	}
}

// SecurityOAuth2 creates a security scheme for OAuth 2.0 authentication.
func SecurityOAuth2(flows spec.OAuthFlows) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.Oauth2 = &spec.SecuritySchemeOAuth2{
			Flows: flows,
		}
	}
}
