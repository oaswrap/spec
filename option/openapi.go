package option

import (
	"log"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/pkg/util"
)

type OpenAPIOption func(*spec.Config)

// WithOpenAPIVersion sets the OpenAPI version for the documentation.
// The default version is "3.1.0".
// Supported versions are "3.0.0" and "3.1.0".
func WithOpenAPIVersion(version string) OpenAPIOption {
	return func(c *spec.Config) {
		c.OpenAPIVersion = version
	}
}

// WithDisableOpenAPI disables the OpenAPI documentation generation.
// By default, OpenAPI documentation generation is enabled.
// If set to true, the OpenAPI documentation will not be generated.
//
// This can be useful in production environments where you want to disable the documentation.
func WithDisableOpenAPI(disable ...bool) OpenAPIOption {
	return func(c *spec.Config) {
		c.DisableOpenAPI = util.Optional(true, disable...)
	}
}

// WithBaseURL sets the base URL for the OpenAPI documentation.
func WithBaseURL(baseURL string) OpenAPIOption {
	return func(c *spec.Config) {
		c.BaseURL = baseURL
	}
}

// WithTitle sets the title for the OpenAPI documentation.
func WithTitle(title string) OpenAPIOption {
	return func(c *spec.Config) {
		c.Title = title
	}
}

// WithVersion sets the version for the OpenAPI documentation.
func WithVersion(version string) OpenAPIOption {
	return func(c *spec.Config) {
		c.Version = version
	}
}

// WithDescription sets the description for the OpenAPI documentation.
func WithDescription(description string) OpenAPIOption {
	return func(c *spec.Config) {
		c.Description = &description
	}
}

// WithServer adds a server to the OpenAPI documentation.
func WithServer(url string, opts ...ServerOption) OpenAPIOption {
	return func(c *spec.Config) {
		server := spec.Server{
			URL: url,
		}
		for _, opt := range opts {
			opt(&server)
		}
		c.Servers = append(c.Servers, server)
	}
}

// WithDocsPath sets the path for the OpenAPI documentation.
func WithDocsPath(path string) OpenAPIOption {
	return func(c *spec.Config) {
		c.DocsPath = path
	}
}

// WithSecurity adds a security scheme to the OpenAPI documentation.
//
// It can be used to define API key or HTTP Bearer authentication schemes.
func WithSecurity(name string, opts ...SecurityOption) OpenAPIOption {
	return func(c *spec.Config) {
		securityConfig := &securityConfig{}
		for _, opt := range opts {
			opt(securityConfig)
		}
		if c.SecuritySchemes == nil {
			c.SecuritySchemes = make(map[string]*spec.SecurityScheme)
		}

		if securityConfig.APIKey != nil {
			c.SecuritySchemes[name] = &spec.SecurityScheme{
				Description: securityConfig.Description,
				APIKey:      securityConfig.APIKey,
			}
		} else if securityConfig.HTTPBearer != nil {
			c.SecuritySchemes[name] = &spec.SecurityScheme{
				Description: securityConfig.Description,
				HTTPBearer:  securityConfig.HTTPBearer,
			}
		} else if securityConfig.Oauth2 != nil {
			c.SecuritySchemes[name] = &spec.SecurityScheme{
				Description: securityConfig.Description,
				OAuth2:      securityConfig.Oauth2,
			}
		} else {
			panic("At least one security scheme must be defined (APIKey, HTTPBearer, or Oauth2)")
		}
	}
}

// WithSwagger sets the configuration for Swagger UI.
func WithSwaggerConfig(cfg ...*spec.SwaggerConfig) OpenAPIOption {
	return func(c *spec.Config) {
		if len(cfg) > 0 && cfg[0] != nil {
			c.SwaggerConfig = cfg[0]
		}
	}
}

// WithDebug enables or disables debug logging for OpenAPI operations.
func WithDebug(debug ...bool) OpenAPIOption {
	return func(c *spec.Config) {
		if util.Optional(true, debug...) {
			c.Logger = log.Default()
		} else {
			c.Logger = &spec.NoopLogger{}
		}
	}
}
