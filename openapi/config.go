package openapi

import (
	"reflect"

	"github.com/swaggest/jsonschema-go"
)

// Config holds the configuration for OpenAPI documentation generation.
type Config struct {
	OpenAPIVersion  string                     // OpenAPI version, e.g., "3.1.0"
	Title           string                     // Title of the API
	Version         string                     // Version of the API
	Description     *string                    // Optional description of the API
	Contact         *Contact                   // Contact information for the API
	License         *License                   // License information for the API
	Servers         []Server                   // List of servers for the API
	SecuritySchemes map[string]*SecurityScheme // Security schemes for the API
	Tags            []Tag                      // Tags for the API
	ExternalDocs    *ExternalDocs              // External documentation for the API

	ReflectorConfig *ReflectorConfig // Configuration for the OpenAPI reflector

	DisableOpenAPI bool
	BaseURL        string
	DocsPath       string
	SwaggerConfig  *SwaggerConfig
	Logger         Logger
}

// ReflectorConfig holds the configuration for the OpenAPI reflector.
type ReflectorConfig struct {
	InlineRefs           bool                 // Whether to inline references in schemas
	RootRef              bool                 // Whether to use a root reference
	RootNullable         bool                 // Whether to allow root schemas to be nullable
	StripDefNamePrefix   []string             // Prefixes to strip from definition names
	InterceptDefNameFunc InterceptDefNameFunc // Function to customize schema definition names
	InterceptPropFunc    InterceptPropFunc    // Function to intercept property schema generation
	InterceptSchemaFunc  InterceptSchemaFunc  // Function to intercept schema generation
	TypeMappings         []TypeMapping        // Custom type mappings for OpenAPI generation
}

// TypeMapping holds a mapping between source and destination types.
type TypeMapping struct {
	Src any // Source type
	Dst any // Destination type
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

// InterceptDefNameFunc is a function type for intercepting schema definition names.
type InterceptDefNameFunc func(t reflect.Type, defaultDefName string) string

// InterceptPropFunc is a function type for intercepting property schema generation.
type InterceptPropFunc func(params InterceptPropParams) error

// InterceptPropParams holds parameters for intercepting property schema generation.
type InterceptPropParams struct {
	Context        *jsonschema.ReflectContext
	Path           []string
	Name           string
	Field          reflect.StructField
	PropertySchema *jsonschema.Schema
	ParentSchema   *jsonschema.Schema
	Processed      bool
}

// InterceptSchemaFunc is a function type for intercepting schema generation.
type InterceptSchemaFunc func(params InterceptSchemaParams) (stop bool, err error)

// InterceptSchemaParams holds parameters for intercepting schema generation.
type InterceptSchemaParams struct {
	Context   *jsonschema.ReflectContext
	Value     reflect.Value
	Schema    *jsonschema.Schema
	Processed bool
}

// Logger is an interface for logging.
type Logger interface {
	Printf(format string, v ...any)
}
