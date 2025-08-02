package option

// ContentUnit defines the structure for OpenAPI content configuration.
type ContentUnit struct {
	Structure   any
	ContentType string

	// HTTPStatus can have values 100-599 for single status, or 1-5 for status families (e.g. 2XX)
	HTTPStatus int
}

// ContentOption is a function that modifies the ContentUnit.
type ContentOption func(cu *ContentUnit)

// WithContentType sets the content type for the OpenAPI content.
func WithContentType(contentType string) ContentOption {
	return func(cu *ContentUnit) {
		cu.ContentType = contentType
	}
}
