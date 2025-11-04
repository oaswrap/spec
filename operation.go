package spec

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/oaswrap/spec/internal/debuglog"
	specopenapi "github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/swaggest/openapi-go"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/openapi-go/openapi31"
)

var _ operationContext = (*operationContextImpl)(nil)

type operationContextImpl struct {
	op                  openapi.OperationContext
	cfg                 *option.OperationConfig
	logger              *debuglog.Logger
	parameterTagMapping map[specopenapi.ParameterIn]string
}

func (oc *operationContextImpl) With(opts ...option.OperationOption) operationContext {
	for _, opt := range opts {
		opt(oc.cfg)
	}
	return oc
}

func (oc *operationContextImpl) build() openapi.OperationContext {
	method := strings.ToUpper(oc.op.Method())
	path := oc.op.PathPattern()

	logger := oc.logger

	cfg := oc.cfg
	if cfg == nil {
		return nil
	}
	if cfg.Hide {
		logger.LogAction("skip operation", fmt.Sprintf("%s %s", method, path))
		return nil
	}
	if cfg.Deprecated {
		oc.op.SetIsDeprecated(true)
		logger.LogOp(method, path, "set is deprecated", "true")
	}
	if cfg.OperationID != "" {
		oc.op.SetID(cfg.OperationID)
		logger.LogOp(method, path, "set operation ID", cfg.OperationID)
	}
	if cfg.Summary != "" {
		oc.op.SetSummary(cfg.Summary)
		logger.LogOp(method, path, "set summary", cfg.Summary)
	}
	if cfg.Description != "" {
		oc.op.SetDescription(cfg.Description)
		logger.LogOp(method, path, "set description", cfg.Description)
	}
	if len(cfg.Tags) > 0 {
		oc.op.SetTags(cfg.Tags...)
		logger.LogOp(method, path, "set tags", fmt.Sprintf("%v", cfg.Tags))
	}
	if len(cfg.Security) > 0 {
		for _, sec := range cfg.Security {
			oc.op.AddSecurity(sec.Name, sec.Scopes...)
		}
		logger.LogOp(method, path, "set security", fmt.Sprintf("%v", cfg.Security))
	}

	for _, req := range cfg.Requests {
		opts, value := oc.buildRequestOpts(req)
		oc.op.AddReqStructure(oc.modifyReqStructure(req.Structure), opts...)
		logger.LogOp(method, path, "add request", value)
	}

	for _, resp := range cfg.Responses {
		opts, value := oc.buildResponseOpts(resp)
		oc.op.AddRespStructure(resp.Structure, opts...)
		logger.LogOp(method, path, "add response", value)
	}

	return oc.op
}

func stringMapToEncodingMap3(enc map[string]string) map[string]openapi3.Encoding {
	res := map[string]openapi3.Encoding{}
	for k, v := range enc {
		rv := v
		res[k] = openapi3.Encoding{
			ContentType: &rv,
		}
	}
	return res
}

func stringMapToEncodingMap31(enc map[string]string) map[string]openapi31.Encoding {
	res := map[string]openapi31.Encoding{}
	for k, v := range enc {
		rv := v
		res[k] = openapi31.Encoding{
			ContentType: &rv,
		}
	}
	return res
}

func (oc *operationContextImpl) buildRequestOpts(req *specopenapi.ContentUnit) ([]openapi.ContentOption, string) {
	log := fmt.Sprintf("%T", req.Structure)
	var opts []openapi.ContentOption
	if req.Description != "" {
		opts = append(opts, func(cu *openapi.ContentUnit) {
			cu.Description = req.Description
		})
		log += fmt.Sprintf(" (%s)", req.Description)
	}
	if req.ContentType != "" {
		opts = append(opts, openapi.WithContentType(req.ContentType))
		log += fmt.Sprintf(" (Content-Type: %s)", req.ContentType)
	}
	opts = append(opts, func(cu *openapi.ContentUnit) {
		cu.Customize = func(cor openapi.ContentOrReference) {
			switch v := cor.(type) {
			case *openapi3.RequestBodyOrRef:
				content := map[string]openapi3.MediaType{}
				for k, val := range v.RequestBody.Content {
					content[k] = *val.WithEncoding(stringMapToEncodingMap3(req.Encoding))
				}
				v.RequestBody.WithContent(content)
			case *openapi31.RequestBodyOrReference:
				content := map[string]openapi31.MediaType{}
				for k, val := range v.RequestBody.Content {
					content[k] = *val.WithEncoding(stringMapToEncodingMap31(req.Encoding))
				}
				v.RequestBody.WithContent(content)
			}
		}
	})
	return opts, log
}

func (oc *operationContextImpl) buildResponseOpts(resp *specopenapi.ContentUnit) ([]openapi.ContentOption, string) {
	log := fmt.Sprintf("%T", resp.Structure)
	var opts []openapi.ContentOption
	if resp.IsDefault {
		opts = append(opts, func(cu *openapi.ContentUnit) {
			cu.IsDefault = true
		})
		log += " (default)"
	}
	if resp.HTTPStatus != 0 {
		opts = append(opts, openapi.WithHTTPStatus(resp.HTTPStatus))
		log += fmt.Sprintf(" (HTTP %d)", resp.HTTPStatus)
	}
	if resp.Description != "" {
		opts = append(opts, func(cu *openapi.ContentUnit) {
			cu.Description = resp.Description
		})
		log += fmt.Sprintf(" (%s)", resp.Description)
	}
	if resp.ContentType != "" {
		opts = append(opts, openapi.WithContentType(resp.ContentType))
		log += fmt.Sprintf(" (Content-Type: %s)", resp.ContentType)
	}
	return opts, log
}

func (oc *operationContextImpl) modifyReqStructure(structure any) any {
	if len(oc.parameterTagMapping) == 0 {
		return structure
	}

	t := reflect.TypeOf(structure)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Only structs are supported for parameter tag modification
	if t.Kind() != reflect.Struct {
		return structure
	}

	fields, modified := oc.buildModifiedFields(t)
	if !modified {
		return structure
	}

	// Create new struct type with modified fields
	newType := reflect.StructOf(fields)
	return reflect.New(newType).Interface()
}

// buildModifiedFields processes struct fields and applies parameter tag mappings.
func (oc *operationContextImpl) buildModifiedFields(t reflect.Type) ([]reflect.StructField, bool) {
	fields := make([]reflect.StructField, 0, t.NumField())
	modified := false

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		originalField := field

		// Apply parameter tag mappings
		for paramIn, sourceTag := range oc.parameterTagMapping {
			if oc.shouldApplyMapping(field, sourceTag, string(paramIn)) {
				field.Tag = oc.buildNewTag(field.Tag, sourceTag, string(paramIn))
				modified = true
			}
		}

		fields = append(fields, field)

		// Log if field was modified (for debugging)
		if field.Tag != originalField.Tag {
			oc.logger.LogAction("modified field tag",
				fmt.Sprintf("field=%s, original=%q, new=%q",
					field.Name, originalField.Tag, field.Tag))
		}
	}

	return fields, modified
}

// shouldApplyMapping determines if a parameter tag mapping should be applied to a field.
func (oc *operationContextImpl) shouldApplyMapping(field reflect.StructField, sourceTag, targetTag string) bool {
	// Only apply if source tag exists and target tag doesn't exist
	return field.Tag.Get(sourceTag) != "" && field.Tag.Get(targetTag) == ""
}

// buildNewTag constructs a new struct tag by adding the mapped parameter tag.
func (oc *operationContextImpl) buildNewTag(
	originalTag reflect.StructTag,
	sourceTag, targetTag string,
) reflect.StructTag {
	sourceValue := originalTag.Get(sourceTag)
	if sourceValue == "" {
		return originalTag
	}

	// Parse existing tag string and add new tag
	tagStr := string(originalTag)
	if tagStr != "" && !strings.HasSuffix(tagStr, " ") {
		tagStr += " "
	}

	// Escape quotes in the tag value
	escapedValue := strings.ReplaceAll(sourceValue, `"`, `\"`)
	newTag := fmt.Sprintf(`%s%s:"%s"`, tagStr, targetTag, escapedValue)

	return reflect.StructTag(newTag)
}
