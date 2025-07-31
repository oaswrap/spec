package option

import (
	"strings"

	"github.com/oaswrap/spec/openapi"
)

// ReflectorOption defines a function that modifies the OpenAPI reflector configuration.
type ReflectorOption func(*openapi.ReflectorConfig)

// InlineRefs sets whether to inline references in the OpenAPI documentation.
//
// If set to true, references will be inlined instead of being stored in the components section.
func InlineRefs() ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.InlineRefs = true
	}
}

// RootRef sets whether to use a root reference in the OpenAPI documentation.
//
// If set to true, the root schema will be used as a reference for all schemas.
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

// InterceptPropFunc sets a function to intercept property schema generation.
//
// This function will be called with the parameters needed to generate the property schema.
func InterceptPropFunc(fn openapi.InterceptPropFunc) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.InterceptPropFunc = fn
	}
}

// RequiredPropByValidateTag sets a function to mark properties as required based on the "validate" tag.
//
// It checks if the "validate" tag contains "required" and adds the property to the required list.
//
// This is useful for automatically marking properties as required based on validation tags because default it use "required:true" tag.
func RequiredPropByValidateTag(tags ...string) ReflectorOption {
	return InterceptPropFunc(func(params openapi.InterceptPropParams) error {
		if !params.Processed {
			return nil
		}
		validateTag := "validate"
		sep := ","
		if len(tags) > 0 {
			validateTag = tags[0]
		}
		if len(tags) > 1 {
			sep = tags[1]
		}
		if v, ok := params.Field.Tag.Lookup(validateTag); ok {
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
//	type NullString struct {
//	     sql.NullString
//	}
//
//	option.WithReflectorConfig(option.TypeMapping(NullString{}, new(string)))
func TypeMapping(src, dst any) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.TypeMappings = append(c.TypeMappings, openapi.TypeMapping{
			Src: src,
			Dst: dst,
		})
	}
}
