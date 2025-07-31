package spec

import (
	"fmt"
	"regexp"

	"github.com/oaswrap/spec/internal/debuglog"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/jsonschema-go"
)

var (
	re3  = regexp.MustCompile(`^3\.0\.\d(-.+)?$`)
	re31 = regexp.MustCompile(`^3\.1\.\d+(-.+)?$`)
)

func newReflector(cfg *openapi.Config) reflector {
	logger := debuglog.NewLogger("spec", cfg.Logger)

	if re3.MatchString(cfg.OpenAPIVersion) {
		return newReflector3(cfg, logger)
	} else if re31.MatchString(cfg.OpenAPIVersion) {
		return newReflector31(cfg, logger)
	}

	logger.Printf("Unsupported OpenAPI version: %s", cfg.OpenAPIVersion)
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

func getJSONSchemaOpts(cfg *openapi.ReflectorConfig, logger *debuglog.Logger) []func(*jsonschema.ReflectContext) {
	var opts []func(*jsonschema.ReflectContext)

	if cfg.InlineRefs {
		opts = append(opts, jsonschema.InlineRefs)
		logger.Printf("set inline references to true")
	}
	if cfg.RootRef {
		opts = append(opts, jsonschema.RootRef)
		logger.Printf("set root reference to true")
	}
	if cfg.RootNullable {
		opts = append(opts, jsonschema.RootNullable)
		logger.Printf("set root nullable to true")
	}
	if len(cfg.StripDefNamePrefix) > 0 {
		opts = append(opts, jsonschema.StripDefinitionNamePrefix(cfg.StripDefNamePrefix...))
		logger.LogAction("set strip definition name prefix", fmt.Sprintf("%v", cfg.StripDefNamePrefix))
	}
	if cfg.InterceptDefNameFunc != nil {
		opts = append(opts, jsonschema.InterceptDefName(cfg.InterceptDefNameFunc))
		logger.Printf("set custom intercept definition name function")
	}
	if cfg.InterceptPropFunc != nil {
		opts = append(opts, jsonschema.InterceptProp(func(params jsonschema.InterceptPropParams) error {
			return cfg.InterceptPropFunc(openapi.InterceptPropParams{
				Context:        params.Context,
				Path:           params.Path,
				Name:           params.Name,
				Field:          params.Field,
				PropertySchema: params.PropertySchema,
				ParentSchema:   params.ParentSchema,
				Processed:      params.Processed,
			})
		}))
		logger.Printf("set custom intercept property function")
	}
	if cfg.InterceptSchemaFunc != nil {
		opts = append(opts, jsonschema.InterceptSchema(func(params jsonschema.InterceptSchemaParams) (stop bool, err error) {
			stop, err = cfg.InterceptSchemaFunc(openapi.InterceptSchemaParams{
				Context:   params.Context,
				Value:     params.Value,
				Schema:    params.Schema,
				Processed: params.Processed,
			})
			return stop, err
		}))
		logger.Printf("set custom intercept schema function")
	}

	return opts
}
