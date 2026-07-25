package validate

import (
	"encoding/json"
	"net/url"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"

	spec_reflect "github.com/oaswrap/spec/internal/reflect"
	"github.com/oaswrap/spec/openapi"
)

var (
	pathParamRe    = regexp.MustCompile(`\{((?:\\.|[^{}\\]|\{(?:\\.|[^{}\\])*\})*)\}`)
	componentRe    = regexp.MustCompile(`^[a-zA-Z0-9.\-_]+$`)
	responseCodeRe = regexp.MustCompile(`^[1-5]([0-9]{2}|XX)$`)
)

func NormalizeTemplatedPath(path string) string {
	return pathParamRe.ReplaceAllString(path, "{}")
}

func ValidateParameterSerialization(context string, param *openapi.Parameter, version string) []error {
	var errs []error
	if param.AllowEmptyValue && param.In != string(openapi.ParameterInQuery) {
		errs = append(errs, Errorf("%s allowEmptyValue is only allowed for query parameters", context))
	}
	if param.AllowEmptyValue && param.In == string(openapi.ParameterInQuery) && IsOpenAPI32(version) {
		errs = append(
			errs,
			Warningf(
				"%s allowEmptyValue is deprecated in OpenAPI 3.2.0 and will be removed in a later version",
				context,
			),
		)
	}
	if param.Style == "" {
		return errs
	}
	if !ValidParameterStyle(param.In, param.Style, version) {
		errs = append(errs, Errorf("%s.style %q is not allowed for %s parameters", context, param.Style, param.In))
	}
	return errs
}

func ValidParameterStyle(in, style, version string) bool {
	switch in {
	case string(openapi.ParameterInPath):
		return slices.Contains([]string{"matrix", "label", "simple"}, style)
	case string(openapi.ParameterInQuery):
		return slices.Contains([]string{"form", "spaceDelimited", "pipeDelimited", "deepObject"}, style)
	case string(openapi.ParameterInHeader):
		return style == "simple"
	case string(openapi.ParameterInCookie):
		if version == openapi.Version320 {
			return style == "form" || style == "cookie"
		}
		return style == "form"
	default:
		return false
	}
}

func MediaTypeAllowsNamedEncoding(mediaType string) bool {
	if mediaType == "" {
		return true
	}
	return MediaTypeIsMultipart(mediaType) || MediaTypeBase(mediaType) == "application/x-www-form-urlencoded"
}

func MediaTypeIsMultipart(mediaType string) bool {
	if mediaType == "" {
		return true
	}
	return strings.HasPrefix(MediaTypeBase(mediaType), "multipart/")
}

func MediaTypeBase(mediaType string) string {
	base, _, _ := strings.Cut(strings.ToLower(strings.TrimSpace(mediaType)), ";")
	return strings.TrimSpace(base)
}

func BodyRefHasInvalidSiblings(body *openapi.RequestBody, version string) bool {
	if spec_reflect.IsOpenAPI30(version) && (body.Summary != "" || body.Description != "") {
		return true
	}
	return len(body.Content) > 0 || body.Required || HasInvalidReferenceExtra(body.Extra, version)
}

func ResponseRefHasInvalidSiblings(response *openapi.Response, version string) bool {
	if spec_reflect.IsOpenAPI30(version) && (response.Summary != "" || response.Description != "") {
		return true
	}
	return len(response.Headers) > 0 || len(response.Content) > 0 || len(response.Links) > 0 ||
		HasInvalidReferenceExtra(response.Extra, version)
}

func HeaderRefHasInvalidSiblings(header *openapi.Header, version string) bool {
	if spec_reflect.IsOpenAPI30(version) && (header.Summary != "" || header.Description != "") {
		return true
	}
	return header.Required || header.Deprecated || header.AllowEmptyValue ||
		header.Style != "" ||
		header.Explode != nil ||
		header.AllowReserved ||
		header.Schema != nil ||
		len(header.Content) > 0 ||
		header.Example != nil ||
		len(header.Examples) > 0 ||
		HasInvalidReferenceExtra(header.Extra, version)
}

func ExampleRefHasInvalidSiblings(example *openapi.Example, version string) bool {
	if spec_reflect.IsOpenAPI30(version) && (example.Summary != "" || example.Description != "") {
		return true
	}
	return example.DataValue != nil ||
		example.Value != nil ||
		example.ExternalValue != "" ||
		example.SerializedValue != "" ||
		HasSerializedExample(example) ||
		HasInvalidReferenceExtra(example.Extra, version)
}

func HasSerializedExample(example *openapi.Example) bool {
	//nolint:staticcheck // Accepted only to detect deprecated pre-fix API usage and report a validation error.
	return example.SerializedExample != nil
}

func LinkRefHasInvalidSiblings(link *openapi.Link, version string) bool {
	if spec_reflect.IsOpenAPI30(version) && (link.Summary != "" || link.Description != "") {
		return true
	}
	return link.OperationRef != "" ||
		link.OperationID != "" ||
		len(link.Parameters) > 0 ||
		link.RequestBody != nil ||
		link.Server != nil ||
		HasInvalidReferenceExtra(link.Extra, version)
}

func CallbackRefHasInvalidSiblings(callback *openapi.Callback, version string) bool {
	if spec_reflect.IsOpenAPI30(version) && (callback.Summary != "" || callback.Description != "") {
		return true
	}
	return len(callback.Expressions) > 0 ||
		HasInvalidReferenceExtra(callback.Extra, version)
}

func MediaTypeRefHasInvalidSiblings(mediaType *openapi.MediaType, version string) bool {
	if spec_reflect.IsOpenAPI30(version) && (mediaType.Summary != "" || mediaType.Description != "") {
		return true
	}
	return mediaType.Schema != nil ||
		mediaType.ItemSchema != nil ||
		mediaType.Example != nil ||
		len(mediaType.Examples) > 0 ||
		len(mediaType.Encoding) > 0 ||
		len(mediaType.PrefixEncoding) > 0 ||
		mediaType.ItemEncoding != nil ||
		HasInvalidReferenceExtra(mediaType.Extra, version)
}

func SecuritySchemeRefHasInvalidSiblings(scheme *openapi.SecurityScheme, version string) bool {
	if spec_reflect.IsOpenAPI30(version) && (scheme.Summary != "" || scheme.Description != nil) {
		return true
	}
	return scheme.Type != "" ||
		scheme.Name != "" ||
		scheme.In != "" ||
		scheme.Scheme != "" ||
		scheme.BearerFormat != nil ||
		scheme.Flows != nil ||
		scheme.OpenIDConnectURL != "" ||
		scheme.OAuth2MetadataURL != "" ||
		scheme.Deprecated ||
		HasInvalidReferenceExtra(scheme.Extra, version)
}

func HasInvalidReferenceExtra(extra map[string]any, version string) bool {
	for key := range extra {
		if !spec_reflect.IsOpenAPI30(version) && (key == "summary" || key == "description") {
			continue
		}
		if !strings.HasPrefix(key, "x-") {
			return true
		}
	}
	return false
}

func SecurityRequirementMayUseURI(name, version string) bool {
	return version == openapi.Version320 && strings.ContainsAny(name, ":/.#?")
}

//nolint:gocyclo,cyclop // explicit sibling checks keep $ref compatibility rules straightforward.
func HasSchemaRefSiblings(schema *openapi.Schema) bool {
	return schema.Schema != "" ||
		schema.ID != "" ||
		len(schema.Defs) > 0 ||
		schema.Anchor != "" ||
		schema.DynamicAnchor != "" ||
		schema.DynamicRef != "" ||
		len(schema.Vocabulary) > 0 ||
		schema.Comment != "" ||
		schema.Title != "" ||
		schema.Description != "" ||
		schema.Type != nil ||
		schema.Format != "" ||
		schema.Nullable ||
		schema.Default != nil ||
		schema.Example != nil ||
		len(schema.Examples) > 0 ||
		len(schema.Enum) > 0 ||
		schema.Const != nil ||
		schema.MultipleOf != nil ||
		schema.Maximum != nil ||
		schema.ExclusiveMaximum != nil ||
		schema.Minimum != nil ||
		schema.ExclusiveMinimum != nil ||
		schema.MaxLength != nil ||
		schema.MinLength != nil ||
		schema.Pattern != "" ||
		schema.MaxItems != nil ||
		schema.MinItems != nil ||
		schema.UniqueItems != nil ||
		schema.MaxProperties != nil ||
		schema.MinProperties != nil ||
		len(schema.Required) > 0 ||
		len(schema.Properties) > 0 ||
		len(schema.PatternProperties) > 0 ||
		schema.Items != nil ||
		len(schema.PrefixItems) > 0 ||
		schema.Contains != nil ||
		schema.MaxContains != nil ||
		schema.MinContains != nil ||
		schema.AdditionalProperties != nil ||
		schema.UnevaluatedProperties != nil ||
		schema.PropertyNames != nil ||
		len(schema.DependentRequired) > 0 ||
		len(schema.DependentSchemas) > 0 ||
		len(schema.AllOf) > 0 ||
		len(schema.AnyOf) > 0 ||
		len(schema.OneOf) > 0 ||
		schema.Not != nil ||
		schema.If != nil ||
		schema.Then != nil ||
		schema.Else != nil ||
		schema.Deprecated ||
		schema.ReadOnly ||
		schema.WriteOnly ||
		schema.ContentEncoding != "" ||
		schema.ContentMediaType != "" ||
		schema.ContentSchema != nil ||
		schema.Discriminator != nil ||
		schema.XML != nil ||
		schema.ExternalDocs != nil ||
		HasNonExtensionExtra(schema.Extra)
}

func SchemaTypeIncludesArray(schema *openapi.Schema) bool {
	if schema == nil {
		return false
	}
	switch value := schema.Type.(type) {
	case string:
		return value == "array"
	case []string:
		return slices.Contains(value, "array")
	case []any:
		for _, item := range value {
			if item == "array" {
				return true
			}
		}
	}
	return false
}

func ExtraHas(extra map[string]any, keys ...string) bool {
	for _, key := range keys {
		if _, ok := extra[key]; ok {
			return true
		}
	}
	return false
}

func HasNonExtensionExtra(extra map[string]any) bool {
	for key := range extra {
		if !strings.HasPrefix(key, "x-") {
			return true
		}
	}
	return false
}

func SchemaBaseURI(value reflect.Value, base string) string {
	field := value.FieldByName("ID")
	if !field.IsValid() || field.Kind() != reflect.String || field.String() == "" {
		return base
	}
	resolved, ok := ResolveURIReference(base, field.String())
	if !ok {
		return base
	}
	return WithoutFragment(resolved)
}

func RegisterSchemaResource(value reflect.Value, base string, resources map[string]any) {
	if !value.CanInterface() {
		return
	}
	if base != "" {
		if _, exists := resources[base]; !exists {
			resources[base] = MarshalAny(value.Interface())
		}
	}
	RegisterSchemaAnchor(value, base, "Anchor", resources)
	RegisterSchemaAnchor(value, base, "DynamicAnchor", resources)
}

func RegisterSchemaAnchor(value reflect.Value, base, fieldName string, resources map[string]any) {
	field := value.FieldByName(fieldName)
	if !field.IsValid() || field.Kind() != reflect.String || field.String() == "" {
		return
	}
	target := base + "#" + field.String()
	if target == "" {
		return
	}
	if _, exists := resources[target]; !exists {
		resources[target] = MarshalAny(value.Interface())
	}
}

func MarshalAny(value any) any {
	raw, err := openapi.MarshalJSON(value)
	if err != nil {
		return nil
	}
	var out any
	if err = json.Unmarshal(raw, &out); err != nil {
		return nil
	}
	return out
}

func ResolveURIReference(base, ref string) (string, bool) {
	refURL, err := url.Parse(ref)
	if err != nil {
		return "", false
	}
	if base == "" {
		return refURL.String(), true
	}
	baseURL, err := url.Parse(base)
	if err != nil {
		return "", false
	}
	return baseURL.ResolveReference(refURL).String(), true
}

func IsNumber(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		return true
	default:
		return false
	}
}

func IsURIReference(value string) bool {
	if strings.ContainsAny(value, " \t\r\n") {
		return false
	}
	_, err := url.Parse(value)
	return err == nil
}

func IsNonRelativeURI(value string) bool {
	if !IsURIReference(value) {
		return false
	}
	parsed, err := url.Parse(value)
	return err == nil && parsed.IsAbs()
}

func IsServerURL(value string) bool {
	if strings.ContainsAny(value, "?#") {
		return false
	}
	return IsURIReference(pathParamRe.ReplaceAllString(value, "x"))
}

func IsHTTPSURI(value string) bool {
	if !IsNonRelativeURI(value) {
		return false
	}
	parsed, err := url.Parse(value)
	return err == nil && strings.EqualFold(parsed.Scheme, "https") && parsed.Fragment == ""
}

func IsLocalReference(ref string, resources map[string]any) bool {
	parsed, err := url.Parse(ref)
	if err != nil {
		return false
	}
	base := urlWithoutFragment(parsed)
	resource, ok := resources[base]
	if !ok {
		return false
	}
	if parsed.Fragment == "" {
		return true
	}
	if _, ok = resources[base+"#"+parsed.Fragment]; ok {
		return true
	}
	if strings.HasPrefix(parsed.Fragment, "/") {
		return ResolveJSONPointer(resource, parsed.Fragment) != nil
	}
	return false
}

func ReferenceTargetExists(ref string, resources map[string]any) bool {
	if _, ok := resources[ref]; ok {
		return true
	}
	parsed, err := url.Parse(ref)
	if err != nil {
		return false
	}
	base := urlWithoutFragment(parsed)
	fragment := parsed.Fragment
	if fragment == "" {
		_, ok := resources[base]
		return ok
	}
	if strings.HasPrefix(fragment, "/") {
		resource, ok := resources[base]
		return ok && ResolveJSONPointer(resource, fragment) != nil
	}
	_, ok := resources[base+"#"+fragment]
	return ok
}

func WithoutFragment(raw string) string {
	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}
	return urlWithoutFragment(parsed)
}

func urlWithoutFragment(parsed *url.URL) string {
	copyURL := *parsed
	copyURL.Fragment = ""
	copyURL.RawFragment = ""
	return copyURL.String()
}

func ResolveJSONPointer(root any, pointer string) any {
	if pointer == "" {
		return root
	}
	if !strings.HasPrefix(pointer, "/") {
		return nil
	}
	current := root
	for _, token := range strings.Split(pointer[1:], "/") {
		token = strings.ReplaceAll(strings.ReplaceAll(token, "~1", "/"), "~0", "~")
		switch node := current.(type) {
		case map[string]any:
			next, ok := node[token]
			if !ok {
				return nil
			}
			current = next
		case []any:
			index, err := strconv.Atoi(token)
			if err != nil || index < 0 || index >= len(node) {
				return nil
			}
			current = node[index]
		default:
			return nil
		}
	}
	return current
}

func IsOpenAPI31(version string) bool {
	return version == openapi.Version310 || version == openapi.Version311 || version == openapi.Version312
}

func IsOpenAPI32(version string) bool {
	return version == openapi.Version320
}
