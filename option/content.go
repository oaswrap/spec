package option

// ContentConfig defines the structure for OpenAPI content configuration.
type ContentConfig struct {
	Structure   any
	ContentType string

	// HTTPStatus can have values 100-599 for single status, or 1-5 for status families (e.g. 2XX)
	HTTPStatus int
}

// ContentOption is a function that modifies the ContentConfig.
type ContentOption func(cu *ContentConfig)

// WithContentType sets the content type for the OpenAPI content.
func WithContentType(contentType string) ContentOption {
	return func(cu *ContentConfig) {
		cu.ContentType = contentType
	}
}
