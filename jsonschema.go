package spec

import (
	"fmt"

	"github.com/oaswrap/spec/internal/debuglog"
	"github.com/oaswrap/spec/openapi"
	"github.com/swaggest/jsonschema-go"
)

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
