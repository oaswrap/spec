package openapi

type UIProvider uint8

const (
	UIProviderStoplightElements UIProvider = iota
	UIProviderSwaggerUI
	UIProviderRedoc
)

// SwaggerUIConfig defines settings for embedding Swagger UI.
type SwaggerUIConfig struct {
	ShowTopBar bool // If true, shows the top bar in Swagger UI.
	HideCurl   bool // If true, hides curl command snippets.
	JsonEditor bool // If true, enables the JSON editor mode.

	// PreAuthorizeApiKey sets initial API key values for authorization.
	PreAuthorizeApiKey map[string]string

	// SettingsUI overrides Swagger UI configuration options.
	// See: https://swagger.io/docs/open-source-tools/swagger-ui/usage/configuration/
	SettingsUI map[string]string
}

// StoplightElementsConfig holds the configuration for the Stoplight Elements.
type StoplightElementsConfig struct {
	HideExport  bool   // Hide the "Export" button on overview section of the documentation.
	HideSchemas bool   // Hide the schemas in the Table of Contents, when using the sidebar layout.
	HideTryIt   bool   // Hide "Try it" feature.
	Layout      string // Layout type, e.g. "sidebar" or "responsive".
	Logo        string // Logo URL to an image that displays as a small square logo next to the title, above the table of contents.
	Router      string // Router type.
}

// RedocConfig holds the configuration for the Redoc.
type RedocConfig struct {
	DisableSearch    bool // Disable search functionality.
	HideDownload     bool // Hides the "Download" button for saving the API definition source file.
	HideSchemaTitles bool // Hides the schema titles in the documentation.
}
