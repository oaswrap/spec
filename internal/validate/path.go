package validate

import (
	"fmt"
	"strings"

	"github.com/oaswrap/spec/internal/reflect"
	"github.com/oaswrap/spec/openapi"
)

func ValidatePathItem(
	path string,
	item *openapi.PathItem,
	version string,
	operationIDs map[string]string,
	securitySchemes map[string]*openapi.SecurityScheme,
	componentParameters map[string]*openapi.Parameter,
) []error {
	var errs []error
	if !strings.HasPrefix(path, "/") {
		errs = append(errs, Errorf("path %q must start with /", path))
	}
	errs = append(
		errs,
		ValidatePathItemOperations(path, item, version, operationIDs, securitySchemes, componentParameters)...)
	return errs
}

func ValidatePathItemOperations(
	context string,
	item *openapi.PathItem,
	version string,
	operationIDs map[string]string,
	securitySchemes map[string]*openapi.SecurityScheme,
	componentParameters map[string]*openapi.Parameter,
) []error {
	var errs []error
	if item == nil {
		return errs
	}
	for i := range item.Servers {
		errs = append(errs, ValidateServer(fmt.Sprintf("%s.servers[%d]", context, i), &item.Servers[i], version)...)
	}
	errs = append(errs, ValidateParameters(context+".parameters", item.Parameters, version, componentParameters)...)
	if version != openapi.Version320 {
		if item.Query != nil {
			errs = append(errs, Errorf("QUERY operation at %s requires OpenAPI 3.2.0", context))
		}
		if len(item.AdditionalOperations) > 0 {
			errs = append(errs, Errorf("additionalOperations at %s requires OpenAPI 3.2.0", context))
		}
	}
	for method, op := range OperationsOf(item) {
		if op == nil {
			continue
		}
		opContext := fmt.Sprintf("%s %s", strings.ToUpper(method), context)
		errs = append(
			errs,
			ValidateOperation(opContext, op, version, operationIDs, securitySchemes, componentParameters)...)
		params := ResolveParameterRefs(append(item.Parameters, op.Parameters...), componentParameters)
		errs = append(errs, ValidatePathParams(context, method, params)...)
		errs = append(errs, ValidateQueryParameterMix(opContext, params)...)
	}
	for method := range item.AdditionalOperations {
		if IsFixedMethod(method) {
			errs = append(
				errs,
				Errorf("additionalOperations at %s must not contain fixed method %s", context, method),
			)
		}
	}
	return errs
}

func ValidatePathParams(path, method string, params []*openapi.Parameter) []error {
	var errs []error
	if !strings.HasPrefix(path, "/") {
		return nil
	}
	matches := pathParamRe.FindAllStringSubmatch(path, -1)
	templateNames := map[string]struct{}{}
	for _, match := range matches {
		name, _, _ := strings.Cut(match[1], ":")
		templateNames[name] = struct{}{}
	}
	declared := map[string]bool{}
	for _, p := range params {
		if p == nil || p.Ref != "" {
			continue
		}
		if p.In == string(openapi.ParameterInPath) {
			declared[p.Name] = p.Required
			if _, ok := templateNames[p.Name]; !ok {
				errs = append(
					errs,
					Errorf(
						"%s %s path parameter %q must match a path template",
						strings.ToUpper(method),
						path,
						p.Name,
					),
				)
			}
		}
	}
	for _, match := range matches {
		name, _, _ := strings.Cut(match[1], ":")
		if required, ok := declared[name]; !ok {
			errs = append(errs, Errorf("%s %s missing path parameter %q", strings.ToUpper(method), path, name))
		} else if !required {
			errs = append(
				errs,
				Errorf("%s %s path parameter %q must be required", strings.ToUpper(method), path, name),
			)
		}
	}
	return errs
}

func OperationsOf(item *openapi.PathItem) map[string]*openapi.Operation {
	ops := map[string]*openapi.Operation{
		"get":     item.Get,
		"put":     item.Put,
		"post":    item.Post,
		"delete":  item.Delete,
		"options": item.Options,
		"head":    item.Head,
		"patch":   item.Patch,
		"trace":   item.Trace,
		"query":   item.Query,
	}
	for method, op := range item.AdditionalOperations {
		ops[method] = op
	}
	return ops
}

func IsFixedMethod(method string) bool {
	switch strings.ToLower(method) {
	case "get", "put", "post", "delete", "options", "head", "patch", "trace", "query":
		return true
	default:
		return false
	}
}

func IsValidParameterIn(in string) bool {
	switch in {
	case "query", "header", "path", "cookie", "querystring":
		return true
	default:
		return false
	}
}

func ResolveParameterRefs(
	params []*openapi.Parameter,
	componentParameters map[string]*openapi.Parameter,
) []*openapi.Parameter {
	if len(params) == 0 {
		return nil
	}
	out := make([]*openapi.Parameter, 0, len(params))
	for _, param := range params {
		if param == nil || param.Ref == "" {
			out = append(out, param)
			continue
		}
		if resolved := ResolveParameterRef(param.Ref, componentParameters); resolved != nil {
			out = append(out, resolved)
			continue
		}
		out = append(out, param)
	}
	return out
}

func ResolveParameterRef(ref string, componentParameters map[string]*openapi.Parameter) *openapi.Parameter {
	const prefix = "#/components/parameters/"
	if !strings.HasPrefix(ref, prefix) {
		return nil
	}
	name := strings.TrimPrefix(ref, prefix)
	name = strings.ReplaceAll(strings.ReplaceAll(name, "~1", "/"), "~0", "~")
	return componentParameters[name]
}

func HasParameterRefSiblings(param *openapi.Parameter, version string) bool {
	if reflect.IsOpenAPI30(version) && (param.Summary != "" || param.Description != "") {
		return true
	}
	return param.Name != "" || param.In != "" || param.Required || param.Deprecated ||
		param.AllowEmptyValue ||
		param.Style != "" ||
		param.Explode != nil ||
		param.AllowReserved ||
		param.Schema != nil ||
		len(param.Content) > 0 ||
		param.Example != nil ||
		len(param.Examples) > 0 ||
		HasInvalidReferenceExtra(param.Extra, version)
}
