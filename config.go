package openapiwrapper

// Config holds the configuration for OpenAPI documentation generation.
type Config struct {
	OpenAPIVersion  string // OpenAPI version, e.g., "3.1.0"
	Title           string
	Version         string
	Description     *string
	Servers         []Server
	SecuritySchemes map[string]*SecurityScheme

	DisableOpenAPI bool
	BaseURL        string
	DocsPath       string
	SwaggerConfig  *SwaggerConfig

	Logger Logger
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
