package option

import (
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/pkg/util"
)

// ContentOption is a function that modifies the ContentUnit.
type ContentOption func(cu *openapi.ContentUnit)

// WithContentType sets the content type for the OpenAPI content.
func WithContentType(contentType string) ContentOption {
	return func(cu *openapi.ContentUnit) {
		cu.ContentType = contentType
	}
}

// WithContentDescription sets the description for the OpenAPI content.
func WithContentDescription(description string) ContentOption {
	return func(cu *openapi.ContentUnit) {
		cu.Description = description
	}
}

// WithContentDefault sets whether this content unit is the default response.
func WithContentDefault(isDefault ...bool) ContentOption {
	return func(cu *openapi.ContentUnit) {
		cu.IsDefault = util.Optional(true, isDefault...)
	}
}
