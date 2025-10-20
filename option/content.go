package option

import (
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/pkg/util"
)

// ContentOption is a function that modifies the ContentUnit.
type ContentOption func(cu *openapi.ContentUnit)

// ContentType sets the content type for the OpenAPI content.
func ContentType(contentType string) ContentOption {
	return func(cu *openapi.ContentUnit) {
		cu.ContentType = contentType
	}
}

// ContentDescription sets the description for the OpenAPI content.
func ContentDescription(description string) ContentOption {
	return func(cu *openapi.ContentUnit) {
		cu.Description = description
	}
}

// ContentDefault sets whether this content unit is the default response.
func ContentDefault(isDefault ...bool) ContentOption {
	return func(cu *openapi.ContentUnit) {
		cu.IsDefault = util.Optional(true, isDefault...)
	}
}

func ContentEncoding(prop, enc string) ContentOption {
	return func(cu *openapi.ContentUnit) {
		cu.Encoding[prop] = enc
	}
}