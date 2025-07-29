package option

import (
	openapiwrapper "github.com/faizlabs/openapi-wrapper"
	"github.com/faizlabs/openapi-wrapper/internal/util"
)

type securityConfig struct {
	Description *string
	APIKey      *openapiwrapper.SecuritySchemeAPIKey
	HTTPBearer  *openapiwrapper.SecuritySchemeHTTPBearer
	Oauth2      *openapiwrapper.SecuritySchemeOAuth2
}

// SecurityOption is a function that applies configuration to a securityConfig.
type SecurityOption func(*securityConfig)

// SecurityDescription sets the description for the security scheme.
func SecurityDescription(description string) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.Description = &description
	}
}

// SecurityAPIKey creates a security scheme for API key authentication.
func SecurityAPIKey(name string, in openapiwrapper.SecuritySchemeAPIKeyIn) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.APIKey = &openapiwrapper.SecuritySchemeAPIKey{
			Name: name,
			In:   in,
		}
	}
}

// SecurityHTTPBearer creates a security scheme for HTTP Bearer authentication.
func SecurityHTTPBearer(scheme ...string) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.HTTPBearer = &openapiwrapper.SecuritySchemeHTTPBearer{
			Scheme: util.Optional("Bearer", scheme...),
		}
	}
}

// SecurityOAuth2 creates a security scheme for OAuth 2.0 authentication.
func SecurityOAuth2(flows openapiwrapper.OAuthFlows) SecurityOption {
	return func(cfg *securityConfig) {
		cfg.Oauth2 = &openapiwrapper.SecuritySchemeOAuth2{
			Flows: flows,
		}
	}
}
