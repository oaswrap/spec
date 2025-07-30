package openapi

// Config holds the configuration for OpenAPI documentation generation.
type Config struct {
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

// Logger is an interface for logging.
type Logger interface {
	Printf(format string, v ...any)
}
