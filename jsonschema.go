package spec

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/oaswrap/spec/internal/debuglog"
	"github.com/oaswrap/spec/openapi"
	"github.com/swaggest/jsonschema-go"
)

var genericInstRe = regexp.MustCompile(`^(\w+)\[(.+)\]$`)

func getJSONSchemaOpts(cfg *openapi.ReflectorConfig, logger *debuglog.Logger) []func(*jsonschema.ReflectContext) {
	var opts []func(*jsonschema.ReflectContext)

	if cfg == nil {
		opts = append(opts, jsonschema.InterceptDefName(shortenGenericName))
		return opts
	}

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
	opts = append(opts, jsonschema.InterceptDefName(shortenGenericName))
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
		opts = append(
			opts,
			jsonschema.InterceptSchema(func(params jsonschema.InterceptSchemaParams) (bool, error) {
				stop, err := cfg.InterceptSchemaFunc(openapi.InterceptSchemaParams{
					Context:   params.Context,
					Value:     params.Value,
					Schema:    params.Schema,
					Processed: params.Processed,
				})
				return stop, err
			}),
		)
		logger.Printf("set custom intercept schema function")
	}

	return opts
}

// shortenGenericName converts "Page[some/pkg.Item]" to "PageItem".
func shortenGenericName(t reflect.Type, defaultDefName string) string {
	m := genericInstRe.FindStringSubmatch(t.Name())
	if m == nil {
		return defaultDefName
	}
	// Use the container name from defaultDefName, which already has the package
	// prefix applied and StripDefinitionNamePrefix already run — so the result
	// is consistent with how non-generic struct names are generated.
	containerName := m[1]
	if before, _, found := strings.Cut(defaultDefName, "["); found {
		containerName = before
	}
	args := strings.Split(m[2], ", ")
	result := containerName
	var sb strings.Builder
	for _, arg := range args {
		arg = strings.TrimPrefix(arg, "*")
		var suffixSb strings.Builder
		for strings.HasPrefix(arg, "[]") {
			suffixSb.WriteString("List")
			arg = arg[2:]
		}
		arg = strings.TrimPrefix(arg, "*")
		if i := strings.LastIndex(arg, "."); i >= 0 {
			arg = arg[i+1:]
		}
		sb.WriteString(arg + suffixSb.String())
	}
	result += sb.String()
	return result
}
