package openapi

import (
	"reflect"

	"github.com/swaggest/jsonschema-go"
)

// Config defines the root configuration for OpenAPI documentation generation.
type Config struct {
	OpenAPIVersion  string                     // OpenAPI version, e.g., "3.1.0".
	Title           string                     // Title of the API.
	Version         string                     // Version of the API.
	Description     *string                    // Optional description of the API.
	Contact         *Contact                   // Contact information for the API.
	License         *License                   // License information for the API.
	TermsOfService  *string                    // Terms of service URL.
	Servers         []Server                   // List of API servers.
	SecuritySchemes map[string]*SecurityScheme // Security schemes available for the API.
	Tags            []Tag                      // Tags used to organize operations.
	ExternalDocs    *ExternalDocs              // Additional external documentation.

	ReflectorConfig *ReflectorConfig // Configuration for schema reflection.

	BaseURL       string         // Base URL of the API.
	DocsPath      string         // Path where the documentation will be served.
	DisableDocs   bool           // If true, disables serving OpenAPI docs.
	Logger        Logger         // Logger for diagnostic output.
	SwaggerConfig *SwaggerConfig // Configuration for embedded Swagger UI.
	PathParser    PathParser     // Path parser for framework-specific path conversions.
}

// ReflectorConfig holds advanced options for schema reflection.
type ReflectorConfig struct {
	InlineRefs           bool                 // If true, inline schema references instead of using components.
	RootRef              bool                 // If true, use a root reference for top-level schemas.
	RootNullable         bool                 // If true, allow root schemas to be nullable.
	StripDefNamePrefix   []string             // Prefixes to strip from generated definition names.
	InterceptDefNameFunc InterceptDefNameFunc // Function to customize definition names.
	InterceptPropFunc    InterceptPropFunc    // Function to intercept property schema generation.
	InterceptSchemaFunc  InterceptSchemaFunc  // Function to intercept full schema generation.
	TypeMappings         []TypeMapping        // Custom type mappings for schema generation.
}

// TypeMapping maps a source type to a target type in schema generation.
type TypeMapping struct {
	Src any // Source type.
	Dst any // Destination type.
}

// SwaggerConfig defines settings for embedding Swagger UI.
type SwaggerConfig struct {
	ShowTopBar bool // If true, shows the top bar in Swagger UI.
	HideCurl   bool // If true, hides curl command snippets.
	JsonEditor bool // If true, enables the JSON editor mode.

	// PreAuthorizeApiKey sets initial API key values for authorization.
	PreAuthorizeApiKey map[string]string

	// SettingsUI overrides Swagger UI configuration options.
	// See: https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/
	SettingsUI map[string]string

	// Proxy enables proxying requests through the Swagger UI handler.
	// Useful for avoiding CORS issues when the API server is not directly accessible.
	Proxy bool
}

// InterceptDefNameFunc allows customizing schema definition names.
type InterceptDefNameFunc func(t reflect.Type, defaultDefName string) string

// InterceptPropFunc allows customizing property schemas during generation.
type InterceptPropFunc func(params InterceptPropParams) error

// InterceptPropParams holds information for intercepting property generation.
type InterceptPropParams struct {
	Context        *jsonschema.ReflectContext // Reflection context.
	Path           []string                   // Path to the property.
	Name           string                     // Property name.
	Field          reflect.StructField        // Struct field being processed.
	PropertySchema *jsonschema.Schema         // Generated property schema.
	ParentSchema   *jsonschema.Schema         // Parent object schema.
	Processed      bool                       // True if the property was already processed.
}

// InterceptSchemaFunc allows intercepting schema generation for entire types.
type InterceptSchemaFunc func(params InterceptSchemaParams) (stop bool, err error)

// InterceptSchemaParams holds information for intercepting full schema generation.
type InterceptSchemaParams struct {
	Context   *jsonschema.ReflectContext // Reflection context.
	Value     reflect.Value              // Value being reflected.
	Schema    *jsonschema.Schema         // Generated schema.
	Processed bool                       // True if the schema was already processed.
}

// Logger defines an interface for logging diagnostic messages.
type Logger interface {
	Printf(format string, v ...any)
}

// PathParser defines an interface for converting router paths to OpenAPI paths.
//
// Example:
//
//	Input: "/users/:id"
//	Output: "/users/{id}"
type PathParser interface {
	// Parse converts a framework-style path to OpenAPI path syntax.
	Parse(path string) (string, error)
}
