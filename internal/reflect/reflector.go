package reflect

import (
	"errors"
	"io"
	"log/slog"
	"reflect"
	"slices"
	"strings"
	"time"

	"github.com/oaswrap/spec/openapi"
)

type SchemaMode int

const (
	SchemaInline SchemaMode = iota
	SchemaUseComponent
)

type Reflector struct {
	Config      *openapi.Config
	Components  map[string]*openapi.Schema
	Names       map[reflect.Type]string
	Generating  map[reflect.Type]bool
	TypeMapping map[reflect.Type]reflect.Type
}

func NewReflector(cfg *openapi.Config) *Reflector {
	if cfg == nil {
		cfg = &openapi.Config{}
	}
	if cfg.Logger == nil {
		cfg.Logger = slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	r := &Reflector{
		Config:      cfg,
		Components:  map[string]*openapi.Schema{},
		Names:       map[reflect.Type]string{},
		Generating:  map[reflect.Type]bool{},
		TypeMapping: map[reflect.Type]reflect.Type{},
	}
	if cfg.ReflectorConfig != nil {
		for _, tm := range cfg.ReflectorConfig.TypeMappings {
			src := IndirectType(reflect.TypeOf(tm.Src))
			dst := IndirectType(reflect.TypeOf(tm.Dst))
			if src != nil && dst != nil {
				r.TypeMapping[src] = dst
			}
		}
	}
	return r
}

func (r *Reflector) RequestParts(
	value any,
	ct string,
) ([]*openapi.Parameter, *openapi.Schema, error) {
	t := IndirectType(reflect.TypeOf(value))
	if t == nil {
		return nil, nil, nil
	}
	if mapped := r.TypeMapping[t]; mapped != nil {
		r.Config.Logger.Debug("applying type mapping", "src", t.String(), "dst", mapped.String())
		t = mapped
	}
	if t.Kind() != reflect.Struct || IsTime(t) {
		schema, err := r.SchemaForType(t, SchemaUseComponent, nil)
		return nil, schema, err
	}

	var params []*openapi.Parameter
	bodyTag := r.bodyNameTag(ct)
	hasBody := false
	hasParam := false
	ForEachField(t, func(field reflect.StructField) {
		if paramIn, name, ok := r.ParameterField(field); ok {
			hasParam = true
			params = append(params, r.ParameterSchema(field, paramIn, name))
		}
		if TagName(field, bodyTag) != "" || (bodyTag != "json" && TagName(field, "json") != "") {
			hasBody = true
		}
	})
	if !hasParam {
		schema, err := r.SchemaForType(t, SchemaUseComponent, nil)
		return nil, schema, err
	}
	if !hasBody {
		return params, nil, nil
	}
	body, err := r.StructSchema(t, bodyTag, true, SchemaInline)
	if err != nil {
		return nil, nil, err
	}
	if len(body.Properties) == 0 {
		body = nil
	}
	return params, body, nil
}

func (r *Reflector) ParameterField(field reflect.StructField) (string, string, bool) {
	tagPairs := []struct {
		in  openapi.ParameterIn
		tag string
	}{
		{openapi.ParameterInPath, "path"},
		{openapi.ParameterInQuery, "query"},
		{openapi.ParameterInHeader, "header"},
		{openapi.ParameterInCookie, "cookie"},
	}
	if r.Config.OpenAPIVersion == openapi.Version320 {
		tagPairs = append(tagPairs, struct {
			in  openapi.ParameterIn
			tag string
		}{openapi.ParameterInQueryString, "querystring"})
	}
	custom := make(map[openapi.ParameterIn]string)
	if r.Config.ReflectorConfig != nil {
		for in, tag := range r.Config.ReflectorConfig.ParameterTagMapping {
			// ParameterInBody and ParameterInForm override body tags, not parameter locations.
			if in == openapi.ParameterInBody || in == openapi.ParameterInForm {
				continue
			}
			custom[in] = tag
		}
	}
	for _, pair := range tagPairs {
		if tag, ok := custom[pair.in]; ok {
			delete(custom, pair.in)
			if tag != "" && tag != pair.tag {
				tagPairs = append(tagPairs, struct {
					in  openapi.ParameterIn
					tag string
				}{pair.in, tag})
			}
		}
	}
	for in, tag := range custom {
		tagPairs = append(tagPairs, struct {
			in  openapi.ParameterIn
			tag string
		}{in, tag})
	}
	for _, pair := range tagPairs {
		in, tag := pair.in, pair.tag
		if name := TagName(field, tag); name != "" {
			return string(in), name, true
		}
	}
	return "", "", false
}

func (r *Reflector) ParameterSchema(field reflect.StructField, in, name string) *openapi.Parameter {
	schema, _ := r.SchemaForType(field.Type, SchemaInline, &field)
	param := &openapi.Parameter{
		Name:        name,
		In:          in,
		Description: field.Tag.Get("description"),
		Required:    in == string(openapi.ParameterInPath) || BoolTag(field.Tag.Get("required")),
		Deprecated:  BoolTag(field.Tag.Get("deprecated")),
	}
	if in == string(openapi.ParameterInQueryString) {
		mediaType := field.Tag.Get("mediaType")
		if mediaType == "" {
			mediaType = "application/x-www-form-urlencoded"
		}
		param.Content = map[string]*openapi.MediaType{
			mediaType: {
				Schema: schema,
			},
		}
	} else {
		param.Schema = schema
	}
	return param
}

func (r *Reflector) SchemaForValue(value any, mode SchemaMode) (*openapi.Schema, error) {
	if ov, ok := value.(OneOfValue); ok {
		values := ov.GetValues()
		schemas := make([]*openapi.Schema, 0, len(values))
		for _, item := range values {
			s, err := r.SchemaForValue(item, mode)
			if err != nil {
				return nil, err
			}
			schemas = append(schemas, s)
		}
		return &openapi.Schema{OneOf: schemas}, nil
	}
	if schema, ok := value.(*openapi.Schema); ok {
		return schema, nil
	}
	//nolint:nestif // exposer path needs pre+post hook
	if schema := r.SchemaFromValueExposer(value); schema != nil {
		t := IndirectType(reflect.TypeOf(value))
		r.Config.Logger.Debug("using SchemaExposer bypass", "type", t.String())
		interceptSchema := r.interceptSchemaFn()
		if interceptSchema != nil {
			preSchema := &openapi.Schema{}
			stop, err := interceptSchema(openapi.InterceptSchemaParams{Type: t, Schema: preSchema})
			if err != nil {
				return nil, err
			}
			if stop {
				r.Config.Logger.Debug("interceptSchema: pre-build stopped", "type", t.String())
				return preSchema, nil
			}
			params := openapi.InterceptSchemaParams{Type: t, Schema: schema, Processed: true}
			r.Config.Logger.Debug("interceptSchema: post-build called", "type", t.String())
			if _, err := interceptSchema(params); err != nil {
				return nil, err
			}
		}
		return schema, nil
	}
	return r.SchemaForType(IndirectType(reflect.TypeOf(value)), mode, nil)
}

func (r *Reflector) RefSchema(t reflect.Type) (*openapi.Schema, error) {
	name := r.TypeName(t)
	if _, ok := r.Components[name]; ok {
		return &openapi.Schema{Ref: "#/components/schemas/" + name}, nil
	}
	if r.Generating[t] {
		return &openapi.Schema{Ref: "#/components/schemas/" + name}, nil
	}
	r.Generating[t] = true
	r.Components[name] = &openapi.Schema{}
	r.Config.Logger.Debug("generating component schema", "name", name, "type", t.String())
	interceptSchema := r.interceptSchemaFn()
	if interceptSchema != nil {
		stop, err := interceptSchema(openapi.InterceptSchemaParams{Type: t, Schema: r.Components[name]})
		if err != nil {
			delete(r.Generating, t)
			delete(r.Components, name)
			return nil, err
		}
		if stop {
			r.Config.Logger.Debug("interceptSchema: pre-build stopped", "type", t.String(), "component", name)
			delete(r.Generating, t)
			return &openapi.Schema{Ref: "#/components/schemas/" + name}, nil
		}
	}
	built, err := r.StructSchema(t, "json", false, SchemaInline)
	if err != nil {
		delete(r.Generating, t)
		delete(r.Components, name)
		return nil, err
	}
	// Assign onto the existing pointer so pre-hook customizations on non-overlapping fields survive.
	// StructSchema only sets Type, Properties, and Required.
	r.Components[name].Type = built.Type
	r.Components[name].Properties = built.Properties
	r.Components[name].Required = built.Required
	r.Components[name].AllOf = built.AllOf
	if interceptSchema != nil {
		postParams := openapi.InterceptSchemaParams{Type: t, Schema: r.Components[name], Processed: true}
		r.Config.Logger.Debug("interceptSchema: post-build called", "type", t.String(), "component", name)
		if _, err := interceptSchema(postParams); err != nil {
			delete(r.Generating, t)
			delete(r.Components, name)
			return nil, err
		}
	}
	delete(r.Generating, t)
	return &openapi.Schema{Ref: "#/components/schemas/" + name}, nil
}

//nolint:gocognit // covers full struct field inspection with parameter/body split logic.
func (r *Reflector) StructSchema(
	t reflect.Type,
	nameTag string,
	onlyTagged bool,
	mode SchemaMode,
) (*openapi.Schema, error) {
	schema := &openapi.Schema{Type: "object", Properties: map[string]*openapi.Schema{}}
	// Pre-scan: collect embedded types opted into allOf $ref (via refer:"true" tag or EmbedReferencer).
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.Anonymous || TagName(field, "json") != "" {
			continue
		}
		embType := IndirectType(field.Type)
		if embType == nil || embType.Kind() != reflect.Struct || !isEmbedRef(field) {
			continue
		}
		ref, err := r.RefSchema(embType)
		if err != nil {
			return nil, err
		}
		schema.AllOf = append(schema.AllOf, ref)
	}
	interceptProp := r.interceptPropFn()
	parentType := t
	var firstErr error
	ForEachField(t, func(field reflect.StructField) {
		if firstErr != nil {
			return
		}
		if IgnoredField(field, nameTag) {
			return
		}
		name := TagName(field, nameTag)
		if name == "" && nameTag != "json" {
			name = TagName(field, "json")
		}
		if name == "" {
			if onlyTagged {
				return
			}
			name = LowerCamel(field.Name)
		}
		if interceptProp != nil {
			r.Config.Logger.Debug("interceptProp: field hook called", "field", name, "parent", typeName(parentType))
			if err := interceptProp(openapi.InterceptPropParams{
				Name:         name,
				Field:        field,
				ParentSchema: schema,
				ParentType:   parentType,
			}); err != nil {
				if errors.Is(err, openapi.ErrSkipProperty) {
					return
				}
				firstErr = err
				return
			}
		}
		prop, err := r.SchemaForType(field.Type, mode, &field)
		if err != nil {
			firstErr = err
			return
		}
		schema.Properties[name] = prop
		if BoolTag(field.Tag.Get("required")) {
			schema.Required = append(schema.Required, name)
		}
		if interceptProp != nil {
			if err := r.runPostHook(interceptProp, schema, prop, name, field, parentType); err != nil {
				firstErr = err
				return
			}
		}
	})
	if firstErr != nil {
		return nil, firstErr
	}
	if len(schema.Properties) == 0 {
		schema.Properties = nil
	}
	schema.Required = uniqueStrings(schema.Required)
	return schema, nil
}

// runPostHook calls the post-hook and handles ErrSkipProperty by restoring the snapshot and
// removing the property from schema. Returns a non-nil error only for non-ErrSkipProperty failures.
func (r *Reflector) runPostHook(
	fn openapi.InterceptPropFunc,
	schema *openapi.Schema,
	prop *openapi.Schema,
	name string,
	field reflect.StructField,
	parentType reflect.Type,
) error {
	snap := snapshotParent(schema)
	err := fn(openapi.InterceptPropParams{
		Name:           name,
		Field:          field,
		PropertySchema: prop,
		ParentSchema:   schema,
		Processed:      true,
		ParentType:     parentType,
	})
	if err == nil {
		return nil
	}
	if errors.Is(err, openapi.ErrSkipProperty) {
		restoreParent(schema, snap)
		delete(schema.Properties, name)
		for i, req := range schema.Required {
			if req == name {
				schema.Required = append(schema.Required[:i], schema.Required[i+1:]...)
				break
			}
		}
		return nil
	}
	return err
}

type parentSnapshot struct {
	allOf      []*openapi.Schema
	anyOf      []*openapi.Schema
	oneOf      []*openapi.Schema
	extensions map[string]any
	extra      map[string]any
}

func snapshotParent(s *openapi.Schema) parentSnapshot {
	return parentSnapshot{
		allOf:      s.AllOf,
		anyOf:      s.AnyOf,
		oneOf:      s.OneOf,
		extensions: shallowCopyMap(s.Extensions),
		extra:      shallowCopyMap(s.Extra),
	}
}

func restoreParent(s *openapi.Schema, snap parentSnapshot) {
	s.AllOf = snap.allOf
	s.AnyOf = snap.anyOf
	s.OneOf = snap.oneOf
	s.Extensions = snap.extensions
	s.Extra = snap.extra
}

func shallowCopyMap(m map[string]any) map[string]any {
	if m == nil {
		return nil
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func uniqueStrings(s []string) []string {
	if len(s) == 0 {
		return s
	}
	seen := make(map[string]struct{}, len(s))
	out := s[:0]
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	if len(out) == 0 {
		return nil
	}
	return slices.Clip(out)
}

// SchemaExposer lets a value provide an OpenAPI schema for a specific version.
type SchemaExposer interface {
	OpenAPISchema(version string) *openapi.Schema
}

// StaticSchemaExposer lets a value provide a version-independent OpenAPI schema.
type StaticSchemaExposer interface {
	OpenAPISchema() *openapi.Schema
}

// OneOfValue represents a reflected one-of union container.
type OneOfValue interface {
	GetValues() []any
}

func (r *Reflector) SchemaFromValueExposer(value any) *openapi.Schema {
	if value == nil {
		return nil
	}
	if exposer, ok := value.(SchemaExposer); ok {
		return exposer.OpenAPISchema(r.Config.OpenAPIVersion)
	}
	if exposer, ok := value.(StaticSchemaExposer); ok {
		return exposer.OpenAPISchema()
	}
	return nil
}

func (r *Reflector) SchemaFromTypeExposer(t reflect.Type) *openapi.Schema {
	if t == nil {
		return nil
	}
	if t.Kind() == reflect.Interface {
		return nil
	}
	value := reflect.New(t).Interface()
	return r.SchemaFromValueExposer(value)
}

func IsTime(t reflect.Type) bool {
	return t == reflect.TypeFor[time.Time]()
}

func typeName(t reflect.Type) string {
	if t.Name() == "" {
		return "<anonymous>"
	}
	pkg := t.PkgPath()
	if idx := strings.LastIndex(pkg, "/"); idx >= 0 {
		pkg = pkg[idx+1:]
	}
	if pkg != "" {
		return pkg + "." + t.Name()
	}
	return t.Name()
}
