package option

import (
	"log"

	"github.com/oaswrap/spec/pkg/util"
)

// Config holds the configuration for OpenAPI documentation generation.
type OpenAPI struct {
	OpenAPIVersion string // OpenAPI version, e.g., "3.1.0"

	Title       string   // Title of the API
	Version     string   // Version of the API
	Description *string  // Optional description of the API
	Contact     *Contact // Contact information for the API
	License     *License // License information for the API

	Servers         []Server                   // List of servers for the API
	SecuritySchemes map[string]*SecurityScheme // Security schemes for the API
	Tags            []Tag                      // Tags for the API
	ExternalDocs    *ExternalDocs              // External documentation for the API

	TypeMappings []TypeMapping // Custom type mappings for OpenAPI generation

	DisableOpenAPI bool
	BaseURL        string
	DocsPath       string
	SwaggerConfig  *SwaggerConfig
	Logger         Logger
}

// TypeMapping holds a mapping between source and destination types.
type TypeMapping struct {
	Src any // Source type
	Dst any // Destination type
}

// Contact structure is generated from "#/$defs/contact".
type Contact struct {
	Name          string
	URL           string         // Format: uri.
	Email         string         // Format: email.
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// License structure is generated from "#/$defs/license".
type License struct {
	Name          string // Required.
	Identifier    string
	URL           string         // Format: uri.
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// Tag structure is generated from "#/definitions/Tag".
type Tag struct {
	Name          string // Required.
	Description   string
	ExternalDocs  *ExternalDocs
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// ExternalDocs structure is generated from "#/$defs/external-documentation".
type ExternalDocs struct {
	Description string
	// Format: uri.
	// Required.
	URL           string
	MapOfAnything map[string]any // Key must match pattern: `^x-`.
}

// SwaggerConfig holds the configuration for Swagger UI.
type SwaggerConfig struct {
	ShowTopBar         bool
	HideCurl           bool
	JsonEditor         bool
	PreAuthorizeApiKey map[string]string

	// SettingsUI contains keys and plain javascript values of SwaggerUIBundle configuration.
	// Overrides default values.
	// See https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/ for available options.
	SettingsUI map[string]string

	// Proxy enables proxying requests through swgui handler.
	// Can be useful if API is not directly available due to CORS policy.
	Proxy bool
}

type Logger interface {
	Printf(format string, v ...any)
}

type NoopLogger struct{}

func (l NoopLogger) Printf(format string, v ...any) {}

type OpenAPIOption func(*OpenAPI)

// WithOpenAPIVersion sets the OpenAPI version for the documentation.
// The default version is "3.1.0".
// Supported versions are "3.0.0" and "3.1.0".
func WithOpenAPIVersion(version string) OpenAPIOption {
	return func(c *OpenAPI) {
		c.OpenAPIVersion = version
	}
}

// WithDisableOpenAPI disables the OpenAPI documentation generation.
// By default, OpenAPI documentation generation is enabled.
// If set to true, the OpenAPI documentation will not be generated.
//
// This can be useful in production environments where you want to disable the documentation.
func WithDisableOpenAPI(disable ...bool) OpenAPIOption {
	return func(c *OpenAPI) {
		c.DisableOpenAPI = util.Optional(true, disable...)
	}
}

// WithBaseURL sets the base URL for the OpenAPI documentation.
func WithBaseURL(baseURL string) OpenAPIOption {
	return func(c *OpenAPI) {
		c.BaseURL = baseURL
	}
}

// WithTitle sets the title for the OpenAPI documentation.
func WithTitle(title string) OpenAPIOption {
	return func(c *OpenAPI) {
		c.Title = title
	}
}

// WithVersion sets the version for the OpenAPI documentation.
func WithVersion(version string) OpenAPIOption {
	return func(c *OpenAPI) {
		c.Version = version
	}
}

// WithDescription sets the description for the OpenAPI documentation.
func WithDescription(description string) OpenAPIOption {
	return func(c *OpenAPI) {
		c.Description = &description
	}
}

// WithContact sets the contact information for the OpenAPI documentation.
func WithContact(contact Contact) OpenAPIOption {
	return func(c *OpenAPI) {
		c.Contact = &contact
	}
}

// WithLicense sets the license information for the OpenAPI documentation.
func WithLicense(license License) OpenAPIOption {
	return func(c *OpenAPI) {
		c.License = &license
	}
}

// WithTags adds tags to the OpenAPI documentation.
func WithTags(tags ...Tag) OpenAPIOption {
	return func(c *OpenAPI) {
		c.Tags = append(c.Tags, tags...)
	}
}

// WithServer adds a server to the OpenAPI documentation.
func WithServer(url string, opts ...ServerOption) OpenAPIOption {
	return func(c *OpenAPI) {
		server := Server{
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
	return func(c *OpenAPI) {
		externalDocs := &ExternalDocs{
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
	return func(c *OpenAPI) {
		securityConfig := &securityConfig{}
		for _, opt := range opts {
			opt(securityConfig)
		}
		if c.SecuritySchemes == nil {
			c.SecuritySchemes = make(map[string]*SecurityScheme)
		}

		if securityConfig.APIKey != nil {
			c.SecuritySchemes[name] = &SecurityScheme{
				Description: securityConfig.Description,
				APIKey:      securityConfig.APIKey,
			}
		} else if securityConfig.HTTPBearer != nil {
			c.SecuritySchemes[name] = &SecurityScheme{
				Description: securityConfig.Description,
				HTTPBearer:  securityConfig.HTTPBearer,
			}
		} else if securityConfig.Oauth2 != nil {
			c.SecuritySchemes[name] = &SecurityScheme{
				Description: securityConfig.Description,
				OAuth2:      securityConfig.Oauth2,
			}
		} else {
			panic("At least one security scheme must be defined (APIKey, HTTPBearer, or Oauth2)")
		}
	}
}

// WithTypeMapping adds a type mapping for OpenAPI generation.
//
// Example usage:
//
//	option.WithTypeMapping(types.NullString{}, new(string))
func WithTypeMapping(src, dst any) OpenAPIOption {
	return func(c *OpenAPI) {
		c.TypeMappings = append(c.TypeMappings, TypeMapping{
			Src: src,
			Dst: dst,
		})
	}
}

// WithDocsPath sets the path for the OpenAPI documentation.
func WithDocsPath(path string) OpenAPIOption {
	return func(c *OpenAPI) {
		c.DocsPath = path
	}
}

// WithSwagger sets the configuration for Swagger UI.
func WithSwaggerConfig(cfg ...SwaggerConfig) OpenAPIOption {
	return func(c *OpenAPI) {
		if len(cfg) > 0 {
			c.SwaggerConfig = &cfg[0]
		}
	}
}

// WithDebug enables or disables debug logging for OpenAPI operations.
func WithDebug(debug ...bool) OpenAPIOption {
	return func(c *OpenAPI) {
		if util.Optional(true, debug...) {
			c.Logger = log.Default()
		} else {
			c.Logger = &NoopLogger{}
		}
	}
}
