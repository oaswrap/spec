package validate

import (
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/oaswrap/spec/openapi"
)

type refEntry struct {
	context string
	ref     string
	base    string
}

func ValidateReferenceTargets(doc *openapi.Document) []error {
	root := map[string]any{}
	raw, err := openapi.MarshalJSON(doc)
	if err != nil {
		return []error{Errorf("failed to serialize document for $ref validation: %w", err)}
	}
	err = json.Unmarshal(raw, &root)
	if err != nil {
		return []error{Errorf("failed to parse document for $ref validation: %w", err)}
	}

	resources := map[string]any{"": root}
	if doc.Self != "" {
		resources[WithoutFragment(doc.Self)] = root
	}
	entries := collectRefEntriesAndResources(doc, resources)

	var errs []error
	for _, entry := range entries {
		ref := entry.ref
		if strings.HasPrefix(ref, "/") {
			errs = append(errs, Errorf("%s $ref %q must use #/ for local references", entry.context, ref))
			continue
		}
		resolved, ok := ResolveURIReference(entry.base, ref)
		if !ok {
			errs = append(errs, Errorf("%s $ref %q must be a URI reference", entry.context, ref))
			continue
		}
		parsed, err := url.Parse(resolved)
		if err != nil {
			errs = append(errs, Errorf("%s $ref %q must be a URI reference", entry.context, ref))
			continue
		}
		if _, ok := resources[urlWithoutFragment(parsed)]; !ok {
			continue
		}
		if !ReferenceTargetExists(resolved, resources) {
			errs = append(errs, Errorf("%s $ref %q points to a missing target", entry.context, ref))
		}
	}
	return errs
}

//nolint:gocognit,funlen // validation traversal is complex by nature.
func collectRefEntriesAndResources(doc *openapi.Document, resources map[string]any) []refEntry {
	schemaType := reflect.TypeFor[openapi.Schema]()
	var out []refEntry
	var walk func(value reflect.Value, context, base string)

	walk = func(value reflect.Value, context, base string) {
		if !value.IsValid() {
			return
		}
		for value.Kind() == reflect.Interface || value.Kind() == reflect.Pointer {
			if value.IsNil() {
				return
			}
			value = value.Elem()
		}
		//nolint:exhaustive // only interested in container types.
		switch value.Kind() {
		case reflect.Struct:
			if value.Type() == schemaType {
				base = SchemaBaseURI(value, base)
				RegisterSchemaResource(value, base, resources)
			}
			typ := value.Type()
			for i := range value.NumField() {
				field := typ.Field(i)
				if !field.IsExported() {
					continue
				}
				fieldValue := value.Field(i)
				jsonName, _, _ := strings.Cut(field.Tag.Get("json"), ",")
				if jsonName == "-" {
					if field.Name != "Expressions" {
						continue
					}
					jsonName = ""
				}
				nextContext := context
				if jsonName != "" {
					if context == "" {
						nextContext = jsonName
					} else {
						nextContext = context + "." + jsonName
					}
				}
				if jsonName == "$ref" && fieldValue.Kind() == reflect.String {
					if ref := fieldValue.String(); ref != "" {
						out = append(out, refEntry{context: nextContext, ref: ref, base: base})
					}
					continue
				}
				walk(fieldValue, nextContext, base)
			}
		case reflect.Slice, reflect.Array:
			for i := range value.Len() {
				walk(value.Index(i), fmt.Sprintf("%s[%d]", context, i), base)
			}
		case reflect.Map:
			iter := value.MapRange()
			for iter.Next() {
				key := iter.Key()
				keyStr := fmt.Sprintf("%v", key.Interface())
				nextContext := keyStr
				if context != "" {
					nextContext = context + "." + keyStr
				}
				mapValue := iter.Value()
				if keyStr == "$ref" {
					for mapValue.Kind() == reflect.Interface {
						if mapValue.IsNil() {
							break
						}
						mapValue = mapValue.Elem()
					}
					if mapValue.IsValid() && mapValue.Kind() == reflect.String && mapValue.String() != "" {
						out = append(out, refEntry{context: nextContext, ref: mapValue.String(), base: base})
						continue
					}
				}
				walk(iter.Value(), nextContext, base)
			}
		}
	}

	walk(reflect.ValueOf(doc), "", doc.Self)
	return out
}
