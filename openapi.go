package spec

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/jsonschema-go"
)

var (
	re3  = regexp.MustCompile(`^3\.0\.\d(-.+)?$`)
	re31 = regexp.MustCompile(`^3\.1\.\d+(-.+)?$`)
)

func newReflector(cfg *openapi.Config) reflector {
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
		return newReflector3(cfg, jsonSchemaOpts)
	} else if re31.MatchString(cfg.OpenAPIVersion) {
		return newReflector31(cfg, jsonSchemaOpts)
	}

	return newNoopReflector(fmt.Errorf("unsupported OpenAPI version: %s", cfg.OpenAPIVersion))
}

type noopReflector struct {
	spec   *noopSpec
	errors *SpecError
}

var _ reflector = (*noopReflector)(nil)

func (r *noopReflector) Spec() spec {
	return r.spec
}

func (r *noopReflector) Add(method, path string, opts ...option.OperationOption) {}

func (r *noopReflector) Validate() error {
	if len(r.errors.errors) > 0 {
		return r.errors
	}
	return nil
}

func newNoopReflector(err error) reflector {
	return &noopReflector{
		errors: &SpecError{
			errors: []error{err},
		},
		spec: &noopSpec{},
	}
}

type noopSpec struct{}

func (s *noopSpec) MarshalYAML() ([]byte, error) {
	return nil, nil
}

func (s *noopSpec) MarshalJSON() ([]byte, error) {
	return nil, nil
}
