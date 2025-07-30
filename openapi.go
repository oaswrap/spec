package spec

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/oaswrap/spec/option"
	"github.com/swaggest/jsonschema-go"
)

var (
	re3  = regexp.MustCompile(`^3\.0\.\d(-.+)?$`)
	re31 = regexp.MustCompile(`^3\.1\.\d+(-.+)?$`)
)

func newReflector(cfg *option.OpenAPI) (Reflector, error) {
	jsonSchemaOpts := []func(*jsonschema.ReflectContext){
		jsonschema.InterceptProp(func(params jsonschema.InterceptPropParams) error {
			if !params.Processed {
				return nil
			}
			if v, ok := params.Field.Tag.Lookup("validate"); ok {
				if strings.Contains(v, "required") {
					params.ParentSchema.Required = append(params.ParentSchema.Required, params.Name)
				}
			}
			return nil
		}),
	}

	if re3.MatchString(cfg.OpenAPIVersion) {
		return newReflector3(cfg, jsonSchemaOpts), nil
	} else if re31.MatchString(cfg.OpenAPIVersion) {
		return newReflector31(cfg, jsonSchemaOpts), nil
	}
	return nil, fmt.Errorf("unsupported OpenAPI version: %s", cfg.OpenAPIVersion)
}
