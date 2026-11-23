package option

import (
	"reflect"
	"strings"

	"github.com/oaswrap/spec/openapi"
)

// ReflectorOption mutates schema reflection behavior.
type ReflectorOption func(*openapi.ReflectorConfig)

// InlineRefs toggles inlining referenced schemas.
func InlineRefs(inline ...bool) ReflectorOption {
	return func(cfg *openapi.ReflectorConfig) { cfg.InlineRefs = optional(true, inline...) }
}

// StripDefNamePrefix appends prefixes removed from reflected definition names.
func StripDefNamePrefix(prefixes ...string) ReflectorOption {
	return func(cfg *openapi.ReflectorConfig) {
		cfg.StripDefNamePrefix = append(cfg.StripDefNamePrefix, prefixes...)
	}
}

// InterceptDefName sets callback to customize reflected definition names.
func InterceptDefName(fn func(t reflect.Type, defaultDefName string) string) ReflectorOption {
	return func(cfg *openapi.ReflectorConfig) { cfg.InterceptDefName = fn }
}

// TypeMapping maps source type to destination type during reflection.
func TypeMapping(src, dst any) ReflectorOption {
	return func(cfg *openapi.ReflectorConfig) {
		cfg.TypeMappings = append(cfg.TypeMappings, openapi.TypeMapping{Src: src, Dst: dst})
	}
}

// InterceptSchema sets callback to intercept schema generation per type.
// If a previous hook exists, both are chained: the previous hook runs first.
// If the previous hook returns stop=true or an error, the next hook is not called.
func InterceptSchema(fn openapi.InterceptSchemaFunc) ReflectorOption {
	return func(cfg *openapi.ReflectorConfig) {
		if cfg.InterceptSchema == nil {
			cfg.InterceptSchema = fn
			return
		}
		prev := cfg.InterceptSchema
		cfg.InterceptSchema = func(params openapi.InterceptSchemaParams) (bool, error) {
			stop, err := prev(params)
			if err != nil || stop {
				return stop, err
			}
			return fn(params)
		}
	}
}

// InterceptProp sets callback to intercept property schema generation per field.
// If a previous hook exists, both are chained: the previous hook runs first,
// and the new hook runs only if the previous did not return an error.
func InterceptProp(fn openapi.InterceptPropFunc) ReflectorOption {
	return func(cfg *openapi.ReflectorConfig) {
		if cfg.InterceptProp == nil {
			cfg.InterceptProp = fn
			return
		}
		prev := cfg.InterceptProp
		cfg.InterceptProp = func(params openapi.InterceptPropParams) error {
			if err := prev(params); err != nil {
				return err
			}
			return fn(params)
		}
	}
}

// RequiredPropByValidateTag marks properties as required when their validate tag contains "required".
// Optional args: tags[0] overrides the tag name (default "validate"), tags[1] overrides the separator (default ",").
func RequiredPropByValidateTag(tags ...string) ReflectorOption {
	return InterceptProp(func(params openapi.InterceptPropParams) error {
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
			for _, part := range strings.Split(v, sep) {
				if strings.TrimSpace(part) == "required" {
					params.ParentSchema.Required = append(params.ParentSchema.Required, params.Name)
					break
				}
			}
		}
		return nil
	})
}

// ParameterTagMapping overrides tag source for a specific parameter location.
func ParameterTagMapping(in openapi.ParameterIn, sourceTag string) ReflectorOption {
	return func(cfg *openapi.ReflectorConfig) {
		if cfg.ParameterTagMapping == nil {
			cfg.ParameterTagMapping = make(map[openapi.ParameterIn]string)
		}
		cfg.ParameterTagMapping[in] = sourceTag
	}
}
