package option

import (
	"strings"

	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/pkg/util"
)

// ReflectorOption defines a function that modifies the OpenAPI reflector configuration.
type ReflectorOption func(*openapi.ReflectorConfig)

// InlineRefs sets whether to inline references in the OpenAPI documentation.
func InlineRefs() ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.InlineRefs = true
	}
}

// RootRef sets whether to use a root reference in the OpenAPI documentation.
func RootRef() ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.RootRef = true
	}
}

// RootNullable sets whether to allow root schemas to be nullable in the OpenAPI documentation.
func RootNullable() ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.RootNullable = true
	}
}

// StripDefNamePrefix sets prefixes to strip from definition names in the OpenAPI documentation.
func StripDefNamePrefix(prefixes ...string) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.StripDefNamePrefix = append(c.StripDefNamePrefix, prefixes...)
	}
}

// InterceptDefNameFunc sets a function to customize schema definition names.
//
// This function will be called with the type and the default definition name.
// It should return the desired definition name.
func InterceptDefNameFunc(fn openapi.InterceptDefNameFunc) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.InterceptDefNameFunc = fn
	}
}

// WithInterceptPropFunc sets a function to intercept property schema generation.
//
// This function will be called with the parameters needed to generate the property schema.
func WithInterceptPropFunc(fn openapi.InterceptPropFunc) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.InterceptPropFunc = fn
	}
}

// RequiredPropByValidateTag sets a function to mark properties as required based on the "validate" tag.
//
// If the "validate" tag contains "required", the property will be added to the required list of the parent schema.
// This is useful for automatically generating required properties based on validation tags.
func RequiredPropByValidateTag(seps ...string) ReflectorOption {
	return WithInterceptPropFunc(func(params openapi.InterceptPropParams) error {
		if !params.Processed {
			return nil
		}
		if v, ok := params.Field.Tag.Lookup("validate"); ok {
			sep := util.Optional(",", seps...)
			parts := strings.Split(v, sep)

			for _, part := range parts {
				if strings.TrimSpace(part) == "required" {
					params.ParentSchema.Required = append(params.ParentSchema.Required, params.Name)
					break
				}
			}
		}
		return nil
	})
}

// InterceptSchemaFunc sets a function to intercept schema generation.
//
// This function will be called with the parameters needed to generate the schema.
// It can be used to modify the schema before it is added to the OpenAPI specification.
func InterceptSchemaFunc(fn openapi.InterceptSchemaFunc) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.InterceptSchemaFunc = fn
	}
}

// TypeMapping adds a type mapping for OpenAPI generation.
//
// Example usage:
//
//	option.WithReflectorConfig(
//		option.TypeMapping(types.NullString{}, new(string)),
//	)
func TypeMapping(src, dst any) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.TypeMappings = append(c.TypeMappings, openapi.TypeMapping{
			Src: src,
			Dst: dst,
		})
	}
}
