package option

import (
	"log"

	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/pkg/util"
)

// OpenAPIOption defines a function that applies configuration to an OpenAPI Config.
type OpenAPIOption func(*openapi.Config)

// WithOpenAPIConfig creates a new OpenAPI configuration with the provided options.
// It initializes the configuration with default values and applies the provided options.
func WithOpenAPIConfig(opts ...OpenAPIOption) *openapi.Config {
	cfg := &openapi.Config{
		OpenAPIVersion: "3.0.3",
		Title:          "API Documentation",
		Description:    nil,
		Logger:         &noopLogger{},
		DocsPath:       "/docs",
		SpecPath:       "/docs/openapi.yaml",
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// WithOpenAPIVersion sets the OpenAPI version for the documentation.
//
// The default version is "3.0.3".
// Supported versions are "3.0.3" and "3.1.0".
func WithOpenAPIVersion(version string) OpenAPIOption {
	return func(c *openapi.Config) {
		c.OpenAPIVersion = version
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

// WithTermsOfService sets the terms of service URL for the OpenAPI documentation.
func WithTermsOfService(terms string) OpenAPIOption {
	return func(c *openapi.Config) {
		c.TermsOfService = &terms
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
		if c.ReflectorConfig == nil {
			c.ReflectorConfig = &openapi.ReflectorConfig{}
		}
		for _, opt := range opts {
			opt(c.ReflectorConfig)
		}
	}
}

// WithDisableDocs disables the OpenAPI documentation.
//
// If set to true, the OpenAPI documentation will not be served at the specified path.
// By default, this is false, meaning the documentation is enabled.
func WithDisableDocs(disable ...bool) OpenAPIOption {
	return func(c *openapi.Config) {
		c.DisableDocs = util.Optional(true, disable...)
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

// WithSpecPath sets the path for the OpenAPI specification.
//
// This is the path where the OpenAPI specification will be served.
// The default is "/docs/openapi.yaml"
func WithSpecPath(path string) OpenAPIOption {
	return func(c *openapi.Config) {
		c.SpecPath = path
	}
}

// WithSwaggerUI sets the UI documentation to Swagger UI.
func WithSwaggerUI(cfg ...openapi.SwaggerUIConfig) OpenAPIOption {
	return func(c *openapi.Config) {
		c.UIProvider = openapi.UIProviderSwaggerUI
		if len(cfg) > 0 {
			c.SwaggerUIConfig = &cfg[0]
		}
		if c.SwaggerUIConfig == nil {
			c.SwaggerUIConfig = &openapi.SwaggerUIConfig{}
		}
	}
}

// WithStoplightElements sets the UI documentation to Stoplight Elements.
func WithStoplightElements(cfg ...openapi.StoplightElementsConfig) OpenAPIOption {
	return func(c *openapi.Config) {
		c.UIProvider = openapi.UIProviderStoplightElements
		if len(cfg) > 0 {
			c.StoplightElementsConfig = &cfg[0]
		}
		if c.StoplightElementsConfig == nil {
			c.StoplightElementsConfig = &openapi.StoplightElementsConfig{}
		}
	}
}

// WithRedoc sets the UI documentation to Redoc.
func WithRedoc(cfg ...openapi.RedocConfig) OpenAPIOption {
	return func(c *openapi.Config) {
		c.UIProvider = openapi.UIProviderRedoc
		if len(cfg) > 0 {
			c.RedocConfig = &cfg[0]
		}
		if c.RedocConfig == nil {
			c.RedocConfig = &openapi.RedocConfig{}
		}
	}
}

// WithDebug enables or disables debug logging for OpenAPI operations.
//
// If debug is true, debug logging is enabled, otherwise it is disabled.
// By default, debug logging is disabled.
func WithDebug(debug ...bool) OpenAPIOption {
	return func(c *openapi.Config) {
		if util.Optional(true, debug...) {
			c.Logger = log.Default()
		} else {
			c.Logger = &noopLogger{}
		}
	}
}

// WithPathParser sets a custom path parser for the OpenAPI documentation.
//
// The parser must convert framework-style paths to OpenAPI-style parameter syntax.
// For example, a path like "/users/:id" should be converted to "/users/{id}".
//
// Example:
//
//	// myCustomParser implements PathParser and converts ":param" to "{param}".
//	type myCustomParser struct {
//		re *regexp.Regexp
//	}
//
//	// newMyCustomParser creates an instance with a regexp for colon-prefixed params.
//	func newMyCustomParser() *myCustomParser {
//		return &myCustomParser{
//			re: regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9_]*)`),
//		}
//	}
//
//	// Parse replaces ":param" with "{param}" to match OpenAPI path syntax.
//	func (p *myCustomParser) Parse(path string) (string, error) {
//		return p.re.ReplaceAllString(path, "{$1}"), nil
//	}
//
//	// Example usage:
//	opt := option.WithPathParser(newMyCustomParser())
func WithPathParser(parser openapi.PathParser) OpenAPIOption {
	return func(c *openapi.Config) {
		c.PathParser = parser
	}
}

type noopLogger struct{}

func (l noopLogger) Printf(format string, v ...any) {}
