package option

import (
	"strings"

	"github.com/oaswrap/spec/openapi"
)

// ReflectorOption defines a function that modifies the OpenAPI reflector configuration.
type ReflectorOption func(*openapi.ReflectorConfig)

// InlineRefs sets references to be inlined in the OpenAPI documentation.
//
// When enabled, references will be inlined instead of defined in the components section.
func InlineRefs() ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.InlineRefs = true
	}
}

// RootRef sets whether to use a root reference in the OpenAPI documentation.
//
// When enabled, the root schema will be used as a shared reference for all schemas.
func RootRef() ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.RootRef = true
	}
}

// RootNullable sets whether root schemas are allowed to be nullable.
//
// When enabled, root schemas can be nullable in the OpenAPI documentation.
func RootNullable() ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.RootNullable = true
	}
}

// StripDefNamePrefix specifies one or more prefixes to strip from schema definition names.
func StripDefNamePrefix(prefixes ...string) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.StripDefNamePrefix = append(c.StripDefNamePrefix, prefixes...)
	}
}

// InterceptDefNameFunc sets a custom function for generating schema definition names.
//
// The provided function is called with the type and the default definition name,
// and should return the desired name.
func InterceptDefNameFunc(fn openapi.InterceptDefNameFunc) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.InterceptDefNameFunc = fn
	}
}

// InterceptPropFunc sets a custom function for generating property schemas.
//
// The provided function is called with the parameters for property schema generation.
func InterceptPropFunc(fn openapi.InterceptPropFunc) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.InterceptPropFunc = fn
	}
}

// RequiredPropByValidateTag marks properties as required based on a struct validation tag.
//
// By default, it uses the `validate` tag and looks for the "required" keyword.
//
// You can override the tag name and separator by providing them:
//   - First argument: tag name (default "validate")
//   - Second argument: separator (default ",")
//
// Example:
//
//	option.WithReflectorConfig(option.RequiredPropByValidateTag())
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

// InterceptSchemaFunc sets a custom function for intercepting schema generation.
//
// The provided function is called with the schema generation parameters.
// You can use it to modify schemas before they are added to the OpenAPI output.
func InterceptSchemaFunc(fn openapi.InterceptSchemaFunc) ReflectorOption {
	return func(c *openapi.ReflectorConfig) {
		c.InterceptSchemaFunc = fn
	}
}

// TypeMapping defines a custom type mapping for OpenAPI generation.
//
// Example:
//
//	type NullString struct {
//	    sql.NullString
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
