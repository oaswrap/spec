package option

import (
	"log"

	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/pkg/util"
)

// OpenAPIOption defines a function that modifies the OpenAPI configuration.
type OpenAPIOption func(*openapi.Config)

// WithOpenAPIConfig creates a new OpenAPI configuration with the provided options.
// It initializes the configuration with default values and applies the provided options.
func WithOpenAPIConfig(opts ...OpenAPIOption) *openapi.Config {
	cfg := &openapi.Config{
		OpenAPIVersion:  "3.1.0",
		Title:           "API Documentation",
		Version:         "1.0.0",
		Description:     nil,
		SecuritySchemes: make(map[string]*openapi.SecurityScheme),
		Logger:          &noopLogger{},
		SwaggerConfig:   &openapi.SwaggerConfig{},
		ReflectorConfig: &openapi.ReflectorConfig{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// WithOpenAPIVersion sets the OpenAPI version for the documentation.
// The default version is "3.1.0".
// Supported versions are "3.0.0" and "3.1.0".
func WithOpenAPIVersion(version string) OpenAPIOption {
	return func(c *openapi.Config) {
		c.OpenAPIVersion = version
	}
}

// WithDisableOpenAPI disables the OpenAPI documentation generation.
// By default, OpenAPI documentation generation is enabled.
// If set to true, the OpenAPI documentation will not be generated.
//
// This can be useful in production environments where you want to disable the documentation.
func WithDisableOpenAPI(disable ...bool) OpenAPIOption {
	return func(c *openapi.Config) {
		c.DisableOpenAPI = util.Optional(true, disable...)
	}
}

// WithBaseURL sets the base URL for the OpenAPI documentation.
func WithBaseURL(baseURL string) OpenAPIOption {
	return func(c *openapi.Config) {
		c.BaseURL = baseURL
	}
}

// WithTitle sets the title for the OpenAPI documentation.
func WithTitle(title string) OpenAPIOption {
	return func(c *openapi.Config) {
		c.Title = title
	}
}

// WithVersion sets the version for the OpenAPI documentation.
func WithVersion(version string) OpenAPIOption {
	return func(c *openapi.Config) {
		c.Version = version
	}
}

// WithDescription sets the description for the OpenAPI documentation.
func WithDescription(description string) OpenAPIOption {
	return func(c *openapi.Config) {
		c.Description = &description
	}
}

// WithContact sets the contact information for the OpenAPI documentation.
func WithContact(contact openapi.Contact) OpenAPIOption {
	return func(c *openapi.Config) {
		c.Contact = &contact
	}
}

// WithLicense sets the license information for the OpenAPI documentation.
func WithLicense(license openapi.License) OpenAPIOption {
	return func(c *openapi.Config) {
		c.License = &license
	}
}

// WithTags adds tags to the OpenAPI documentation.
func WithTags(tags ...openapi.Tag) OpenAPIOption {
	return func(c *openapi.Config) {
		c.Tags = append(c.Tags, tags...)
	}
}

// WithServer adds a server to the OpenAPI documentation.
func WithServer(url string, opts ...ServerOption) OpenAPIOption {
	return func(c *openapi.Config) {
		server := openapi.Server{
			URL: url,
		}
		for _, opt := range opts {
			opt(&server)
		}
		c.Servers = append(c.Servers, server)
	}
}

// WithExternalDocs sets the external documentation for the OpenAPI documentation.
func WithExternalDocs(url string, description ...string) OpenAPIOption {
	return func(c *openapi.Config) {
		externalDocs := &openapi.ExternalDocs{
			URL: url,
		}
		if len(description) > 0 {
			externalDocs.Description = description[0]
		}
		c.ExternalDocs = externalDocs
	}
}

// WithSecurity adds a security scheme to the OpenAPI documentation.
//
// It can be used to define API key or HTTP Bearer authentication schemes.
func WithSecurity(name string, opts ...SecurityOption) OpenAPIOption {
	return func(c *openapi.Config) {
		securityConfig := &securityConfig{}
		for _, opt := range opts {
			opt(securityConfig)
		}
		if c.SecuritySchemes == nil {
			c.SecuritySchemes = make(map[string]*openapi.SecurityScheme)
		}

		if securityConfig.APIKey != nil {
			c.SecuritySchemes[name] = &openapi.SecurityScheme{
				Description: securityConfig.Description,
				APIKey:      securityConfig.APIKey,
			}
		} else if securityConfig.HTTPBearer != nil {
			c.SecuritySchemes[name] = &openapi.SecurityScheme{
				Description: securityConfig.Description,
				HTTPBearer:  securityConfig.HTTPBearer,
			}
		} else if securityConfig.Oauth2 != nil {
			c.SecuritySchemes[name] = &openapi.SecurityScheme{
				Description: securityConfig.Description,
				OAuth2:      securityConfig.Oauth2,
			}
		}
	}
}

// WithReflectorConfig applies custom configurations to the OpenAPI reflector.
func WithReflectorConfig(opts ...ReflectorOption) OpenAPIOption {
	return func(c *openapi.Config) {
		for _, opt := range opts {
			opt(c.ReflectorConfig)
		}
	}
}

// WithDocsPath sets the path for the OpenAPI documentation.
//
// This is the path where the OpenAPI documentation will be served.
// The default path is "/docs".
func WithDocsPath(path string) OpenAPIOption {
	return func(c *openapi.Config) {
		c.DocsPath = path
	}
}

// WithSwagger sets the configuration for Swagger UI.
//
// This allows customization of the Swagger UI appearance and behavior.
func WithSwaggerConfig(cfg ...openapi.SwaggerConfig) OpenAPIOption {
	return func(c *openapi.Config) {
		if len(cfg) > 0 {
			c.SwaggerConfig = &cfg[0]
		}
	}
}

// WithDebug enables or disables debug logging for OpenAPI operations.
func WithDebug(debug ...bool) OpenAPIOption {
	return func(c *openapi.Config) {
		if util.Optional(true, debug...) {
			c.Logger = log.Default()
		} else {
			c.Logger = &noopLogger{}
		}
	}
}

type noopLogger struct{}

func (l noopLogger) Printf(format string, v ...any) {}
