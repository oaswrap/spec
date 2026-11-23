package reflect

import (
	"fmt"
	"mime/multipart"
	"path"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	"github.com/oaswrap/spec/openapi"
)

func (r *Reflector) TypeName(t reflect.Type) string {
	if name, ok := r.Names[t]; ok {
		return name
	}
	name := SanitizeTypeName(t.Name())
	name = prefixWithPkg(t, name)
	for _, prefix := range r.StripPrefixes() {
		name = strings.TrimPrefix(name, prefix)
	}
	if r.Config.ReflectorConfig != nil && r.Config.ReflectorConfig.InterceptDefName != nil {
		name = r.Config.ReflectorConfig.InterceptDefName(t, name)
	}
	if name == "" {
		name = "Schema"
	}
	base := name
	i := 2
	for usedType, usedName := range r.Names {
		if usedName == name && usedType != t {
			name = fmt.Sprintf("%s%d", base, i)
			i++
		}
	}
	r.Names[t] = name
	return name
}

func prefixWithPkg(t reflect.Type, defName string) string {
	if defName == "" || t == nil || t.PkgPath() == "" || t.PkgPath() == "main" {
		return defName
	}
	pkgName := path.Base(t.PkgPath())
	if pkgName == "" || pkgName == "main" {
		return defName
	}
	pkgName = sanitizePkgName(pkgName)
	if pkgName == "" {
		return defName
	}
	return pkgName + strings.ToUpper(defName[:1]) + defName[1:]
}

func sanitizePkgName(name string) string {
	parts := strings.FieldsFunc(name, func(r rune) bool {
		return r == '_' || r == '-' || r == '.'
	})
	var b strings.Builder
	for _, p := range parts {
		if len(p) == 0 {
			continue
		}
		b.WriteString(strings.ToUpper(p[:1]) + p[1:])
	}
	return b.String()
}

func (r *Reflector) StripPrefixes() []string {
	if r.Config.ReflectorConfig == nil {
		return nil
	}
	return r.Config.ReflectorConfig.StripDefNamePrefix
}

func (r *Reflector) InlineRefs() bool {
	return r.Config.ReflectorConfig != nil && r.Config.ReflectorConfig.InlineRefs
}

// bodyNameTag returns the struct tag used to name body fields for the given content type.
// Defaults to "json" for JSON bodies and "form" for form bodies.
// Can be overridden via ParameterTagMapping with openapi.ParameterInBody or openapi.ParameterInForm.
func (r *Reflector) bodyNameTag(contentType string) string {
	base := BodyNameTag(contentType)
	if r.Config == nil || r.Config.ReflectorConfig == nil || r.Config.ReflectorConfig.ParameterTagMapping == nil {
		return base
	}
	var key openapi.ParameterIn
	switch base {
	case "json":
		key = openapi.ParameterInBody
	case "form":
		key = openapi.ParameterInForm
	}
	if tag, ok := r.Config.ReflectorConfig.ParameterTagMapping[key]; ok && tag != "" {
		return tag
	}
	return base
}

func (r *Reflector) interceptPropFn() openapi.InterceptPropFunc {
	if r.Config == nil || r.Config.ReflectorConfig == nil {
		return nil
	}
	return r.Config.ReflectorConfig.InterceptProp
}

func (r *Reflector) interceptSchemaFn() openapi.InterceptSchemaFunc {
	if r.Config == nil || r.Config.ReflectorConfig == nil {
		return nil
	}
	return r.Config.ReflectorConfig.InterceptSchema
}

func IndirectType(t reflect.Type) reflect.Type {
	for t != nil && t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return t
}

var typeOfEmbedReferencer = reflect.TypeFor[openapi.EmbedReferencer]()

func isEmbedRef(field reflect.StructField) bool {
	if field.Tag.Get("refer") == "true" {
		return true
	}
	t := IndirectType(field.Type)
	if t == nil {
		return false
	}
	return field.Type.Implements(typeOfEmbedReferencer) || reflect.PointerTo(t).Implements(typeOfEmbedReferencer)
}

func ForEachField(t reflect.Type, fn func(reflect.StructField)) {
	for i := range t.NumField() {
		field := t.Field(i)
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}
		if field.Anonymous && IndirectType(field.Type).Kind() == reflect.Struct && TagName(field, "json") == "" {
			if !isEmbedRef(field) {
				ForEachField(IndirectType(field.Type), fn)
			}
			continue
		}
		fn(field)
	}
}

func LowerCamel(value string) string {
	if value == "" {
		return value
	}
	runes := []rune(value)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

var genericNameRe = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func SanitizeTypeName(name string) string {
	if name == "" {
		return ""
	}

	// Handle slices: []Foo -> FooList
	if strings.HasPrefix(name, "[]") {
		return SanitizeTypeName(name[2:]) + "List"
	}

	// Handle pointers
	name = strings.TrimLeft(name, "*")

	// Handle generics: BaseResponse[github.com/foo.User]
	if start := strings.Index(name, "["); start != -1 && strings.HasSuffix(name, "]") {
		base := name[:start]
		inner := name[start+1 : len(name)-1]

		// Split multiple generic params: Map[string, int]
		parts := strings.Split(inner, ",")
		for i, p := range parts {
			parts[i] = SanitizeTypeName(strings.TrimSpace(p))
		}
		return SanitizeTypeName(base) + strings.Join(parts, "")
	}

	// For names with package paths: github.com/foo.User -> User
	// Note: reflect.Type.Name() usually only contains the local name for defined types,
	// but for generic instances it includes full paths for parameters.
	if lastDot := strings.LastIndex(name, "."); lastDot != -1 {
		name = name[lastDot+1:]
	}

	// Final cleanup for any remaining characters
	return genericNameRe.ReplaceAllString(name, "")
}

var (
	typeFileHeader = reflect.TypeFor[multipart.FileHeader]()
	typeFile       = reflect.TypeFor[multipart.File]()
)

// InferContentType inspects struct tags on value to determine request body content type.
// Returns "multipart/form-data" if any field tagged `form` holds a file type,
// "application/x-www-form-urlencoded" if any `form` tag is present, or "" otherwise.
func InferContentType(value any) string {
	t := IndirectType(reflect.TypeOf(value))
	if t == nil || t.Kind() != reflect.Struct || IsTime(t) {
		return ""
	}
	hasForm := false
	hasFile := false
	ForEachField(t, func(field reflect.StructField) {
		if TagName(field, "form") == "" {
			return
		}
		hasForm = true
		if isFileField(field.Type) {
			hasFile = true
		}
	})
	if !hasForm {
		return ""
	}
	if hasFile {
		return "multipart/form-data"
	}
	return "application/x-www-form-urlencoded"
}

func isFileField(t reflect.Type) bool {
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		t = t.Elem()
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
	}
	return isFileType(t)
}

// isFileType reports whether t is a multipart file type: [multipart.FileHeader],
// [multipart.File], or any type implementing [multipart.File].
func isFileType(t reflect.Type) bool {
	return t == typeFileHeader || t == typeFile || t.Implements(typeFile)
}
