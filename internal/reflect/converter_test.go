package reflect_test

import (
	"mime/multipart"
	std_reflect "reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oaswrap/spec"
	"github.com/oaswrap/spec/internal/reflect"
	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
)

// TextMarshaler test types.

type textStatus int

func (s textStatus) MarshalText() ([]byte, error) { return []byte("ok"), nil }
func (s *textStatus) UnmarshalText([]byte) error  { return nil }

type textStatusJSONOverride int

func (s textStatusJSONOverride) MarshalText() ([]byte, error) { return []byte("ok"), nil }
func (s *textStatusJSONOverride) UnmarshalText([]byte) error  { return nil }
func (s textStatusJSONOverride) MarshalJSON() ([]byte, error) { return []byte(`"ok"`), nil }

// EmbedReferencer test types.

type embedBase struct {
	BaseField string `json:"base_field"`
}

type embedRefViaTag struct {
	embedBase `refer:"true"`

	Name string `json:"name"`
}

type embedBaseViaInterface struct {
	OtherField string `json:"other_field"`
}

func (embedBaseViaInterface) ReferEmbedded() {}

type embedRefViaInterface struct {
	embedBaseViaInterface

	ID int `json:"id"`
}

type CustomSlug string

func (*CustomSlug) OpenAPISchema(version string) *openapi.Schema {
	if strings.HasPrefix(version, "3.0.") {
		return &openapi.Schema{Type: "string", Format: "slug"}
	}
	return &openapi.Schema{Type: []string{"string", "null"}, Format: "slug"}
}

type CustomSchemaPayload struct {
	ID CustomSlug `json:"id" description:"Stable identifier"`
}

type SchemaExposerType struct{}

func (SchemaExposerType) OpenAPISchema(_ string) *openapi.Schema {
	return &openapi.Schema{Type: "string", Description: "Exposed"}
}

func TestConverter_SchemaForType(t *testing.T) {
	cfg := &openapi.Config{OpenAPIVersion: openapi.Version312}
	r := reflect.NewReflector(cfg)

	tests := []struct {
		name     string
		val      any
		expected string
	}{
		{"Int64", int64(0), "integer"},
		{"Uint32", uint32(0), "integer"},
		{"Float32", float32(0), "number"},
		{"ByteSlice", []byte{0}, "string"},
		{"Map", map[string]int{"a": 1}, "object"},
		{"Interface", any(nil), ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typ := std_reflect.TypeOf(tt.val)
			if tt.name == "Interface" {
				typ = std_reflect.TypeFor[any]()
			}
			schema, err := r.SchemaForType(typ, reflect.SchemaInline, nil)
			require.NoError(t, err)
			if tt.expected != "" {
				assert.Equal(t, tt.expected, schema.Type)
			}
		})
	}
}

func TestConverter_ApplyNullable_EdgeCases(t *testing.T) {
	t.Run("OpenAPI304_Ref", func(t *testing.T) {
		r := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version304})
		schema := &openapi.Schema{Ref: "#/components/schemas/User"}
		r.ApplyNullable(schema, true)
		assert.Len(t, schema.AllOf, 1)
		assert.True(t, schema.Nullable)
	})

	t.Run("OpenAPI312_Ref", func(t *testing.T) {
		r := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version312})
		schema := &openapi.Schema{Ref: "#/components/schemas/User"}
		r.ApplyNullable(schema, true)
		assert.Len(t, schema.AnyOf, 2)
		assert.Equal(t, "null", schema.AnyOf[1].Type)
	})
}

func TestConverter_CustomSchemaExposer(t *testing.T) {
	r := spec.NewRouter(option.WithTitle("Custom Schema"), option.WithVersion("1.0.0"))
	r.Post("/payload", option.Request(new(CustomSchemaPayload)), option.Response(204, nil))

	raw, err := r.GenerateSchema("json")
	require.NoError(t, err)

	id := generatedComponentProperty(t, raw, "ReflectTestCustomSchemaPayload", "id")
	assert.Equal(t, "slug", id["format"])
	assert.Equal(t, "Stable identifier", id["description"])
}

func TestConverter_Nullable(t *testing.T) {
	t.Run("OpenAPI304", func(t *testing.T) {
		r := spec.NewRouter(option.WithOpenAPIVersion(openapi.Version304))
		r.Post("/payload", option.Request(new(ReflectionVersionPayload)), option.Response(204, nil))

		raw, err := r.GenerateSchema("json")
		require.NoError(t, err)

		owner := generatedComponentProperty(t, raw, "ReflectTestReflectionVersionPayload", "owner")
		assert.NotContains(t, owner, "$ref", "nullable component refs must not emit $ref siblings in 3.0")
		assert.Equal(t, true, owner["nullable"])
	})

	t.Run("OpenAPI312", func(t *testing.T) {
		r := spec.NewRouter(option.WithOpenAPIVersion(openapi.Version312))
		r.Post("/payload", option.Request(new(ReflectionVersionPayload312)), option.Response(204, nil))

		raw, err := r.GenerateSchema("json")
		require.NoError(t, err)

		name := generatedComponentProperty(t, raw, "ReflectTestReflectionVersionPayload312", "name")
		typ := name["type"].([]any)
		assert.Equal(t, []any{"string", "null"}, typ)
	})
}

func TestConverter_SchemaExposer(t *testing.T) {
	r := spec.NewRouter(option.WithTitle("Exposer"))
	r.Get("/exposer", option.Response(200, SchemaExposerType{}))
	_, err := r.GenerateSchema("yaml")
	require.NoError(t, err)
	doc := r.Document()
	assert.Equal(
		t,
		"Exposed",
		doc.Paths["/exposer"].Get.Responses["200"].Content["application/json"].Schema.Description,
	)
}

func TestConverter_TextMarshaler(t *testing.T) {
	r := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version312})

	t.Run("direct type produces string schema", func(t *testing.T) {
		s, err := r.SchemaForType(std_reflect.TypeFor[textStatus](), reflect.SchemaInline, nil)
		require.NoError(t, err)
		assert.Equal(t, "string", s.Type)
	})

	t.Run("pointer type preserves nullable", func(t *testing.T) {
		s, err := r.SchemaForType(std_reflect.TypeFor[*textStatus](), reflect.SchemaInline, nil)
		require.NoError(t, err)
		assert.Equal(t, "string", s.Type.([]string)[0])
		assert.Equal(t, "null", s.Type.([]string)[1])
	})

	t.Run("not a component type", func(t *testing.T) {
		assert.False(t, reflect.IsComponentType(std_reflect.TypeFor[textStatus]()))
	})

	t.Run("json.Marshaler takes precedence — not reflected as string", func(t *testing.T) {
		// textStatusJSONOverride implements TextMarshaler+TextUnmarshaler AND json.Marshaler.
		// json.Marshaler wins → falls through to normal kindSwitch → integer.
		s, err := r.SchemaForType(std_reflect.TypeFor[textStatusJSONOverride](), reflect.SchemaInline, nil)
		require.NoError(t, err)
		assert.Equal(t, "integer", s.Type, "must not receive string/text-marshaler treatment")
	})
}

func TestConverter_EmbedRef(t *testing.T) {
	t.Run("refer tag produces allOf ref, fields not inlined", func(t *testing.T) {
		r := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version312})
		schema, err := r.StructSchema(std_reflect.TypeFor[embedRefViaTag](), "json", false, reflect.SchemaUseComponent)
		require.NoError(t, err)

		require.Len(t, schema.AllOf, 1, "embedded type must appear in allOf")
		assert.Equal(t, "#/components/schemas/ReflectTestEmbedBase", schema.AllOf[0].Ref)

		_, hasBaseField := schema.Properties["base_field"]
		assert.False(t, hasBaseField, "base_field must not be inlined into parent properties")

		_, hasName := schema.Properties["name"]
		assert.True(t, hasName, "own fields must still be present")
	})

	t.Run("EmbedReferencer interface produces allOf ref", func(t *testing.T) {
		r := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version312})
		schema, err := r.StructSchema(
			std_reflect.TypeFor[embedRefViaInterface](),
			"json",
			false,
			reflect.SchemaUseComponent,
		)
		require.NoError(t, err)

		require.Len(t, schema.AllOf, 1)
		assert.Equal(t, "#/components/schemas/ReflectTestEmbedBaseViaInterface", schema.AllOf[0].Ref)

		_, hasOther := schema.Properties["other_field"]
		assert.False(t, hasOther, "interface-embedded fields must not be inlined")

		_, hasID := schema.Properties["id"]
		assert.True(t, hasID)
	})

	t.Run("plain embed without refer tag still inlines fields", func(t *testing.T) {
		type plainBase struct {
			BaseVal string `json:"base_val"`
		}
		type plainParent struct {
			plainBase

			Extra string `json:"extra"`
		}
		r := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version312})
		schema, err := r.StructSchema(std_reflect.TypeFor[plainParent](), "json", false, reflect.SchemaInline)
		require.NoError(t, err)
		assert.Empty(t, schema.AllOf, "no allOf for plain embed")
		_, hasBase := schema.Properties["base_val"]
		assert.True(t, hasBase, "plain embed fields must be inlined")
	})
}

func TestConverter_SchemaForType_Branches(t *testing.T) {
	t.Run("primitive and collection kinds", func(t *testing.T) {
		r := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version312})

		cases := []struct {
			name   string
			val    any
			typ    std_reflect.Type
			assert func(*testing.T, *openapi.Schema)
		}{
			{
				name: "bool",
				val:  true,
				assert: func(t *testing.T, s *openapi.Schema) {
					assert.Equal(t, "boolean", s.Type)
				},
			},
			{
				name: "int32",
				val:  int32(1),
				assert: func(t *testing.T, s *openapi.Schema) {
					assert.Equal(t, "integer", s.Type)
					assert.Equal(t, "int32", s.Format)
				},
			},
			{
				name: "uint16",
				val:  uint16(1),
				assert: func(t *testing.T, s *openapi.Schema) {
					assert.Equal(t, "integer", s.Type)
					require.NotNil(t, s.Minimum)
					assert.InDelta(t, 0.0, *s.Minimum, 0.0001)
				},
			},
			{
				name: "uint64",
				val:  uint64(1),
				assert: func(t *testing.T, s *openapi.Schema) {
					assert.Equal(t, "integer", s.Type)
					assert.Equal(t, "int64", s.Format)
				},
			},
			{
				name: "float64",
				val:  float64(1.25),
				assert: func(t *testing.T, s *openapi.Schema) {
					assert.Equal(t, "number", s.Type)
					assert.Equal(t, "double", s.Format)
				},
			},
			{
				name: "array",
				val:  [2]int{1, 2},
				assert: func(t *testing.T, s *openapi.Schema) {
					assert.Equal(t, "array", s.Type)
					require.NotNil(t, s.Items)
					assert.Equal(t, "integer", s.Items.Type)
				},
			},
			{
				name: "map",
				val:  map[string]bool{"ok": true},
				assert: func(t *testing.T, s *openapi.Schema) {
					assert.Equal(t, "object", s.Type)
					require.NotNil(t, s.AdditionalProperties)
					additional, ok := s.AdditionalProperties.(*openapi.Schema)
					require.True(t, ok)
					assert.Equal(t, "boolean", additional.Type)
				},
			},
			{
				name: "time",
				val:  time.Time{},
				assert: func(t *testing.T, s *openapi.Schema) {
					assert.Equal(t, "string", s.Type)
					assert.Equal(t, "date-time", s.Format)
				},
			},
			{
				name: "interface",
				typ:  std_reflect.TypeFor[any](),
				assert: func(t *testing.T, s *openapi.Schema) {
					assert.Equal(t, &openapi.Schema{}, s)
				},
			},
			{
				name: "default",
				val:  func() {},
				assert: func(t *testing.T, s *openapi.Schema) {
					assert.Equal(t, &openapi.Schema{}, s)
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				typ := tc.typ
				if typ == nil {
					typ = std_reflect.TypeOf(tc.val)
				}
				s, err := r.SchemaForType(typ, reflect.SchemaInline, nil)
				require.NoError(t, err)
				tc.assert(t, s)
			})
		}
	})

	t.Run("byte slice encoding differs by version", func(t *testing.T) {
		r304 := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version304})
		s304, err := r304.SchemaForType(std_reflect.TypeFor[[]byte](), reflect.SchemaInline, nil)
		require.NoError(t, err)
		assert.Equal(t, "string", s304.Type)
		assert.Equal(t, "byte", s304.Format)

		r312 := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version312})
		s312, err := r312.SchemaForType(std_reflect.TypeFor[[]byte](), reflect.SchemaInline, nil)
		require.NoError(t, err)
		assert.Equal(t, "string", s312.Type)
		assert.Equal(t, "base64", s312.ContentEncoding)
	})

	t.Run("nullable pointer to component schema", func(t *testing.T) {
		type User struct {
			ID string `json:"id"`
		}

		r304 := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version304})
		s304, err := r304.SchemaForType(std_reflect.TypeFor[*User](), reflect.SchemaUseComponent, nil)
		require.NoError(t, err)
		assert.True(t, s304.Nullable)
		require.Len(t, s304.AllOf, 1)
		assert.Equal(t, "#/components/schemas/ReflectTestUser", s304.AllOf[0].Ref)

		r312 := reflect.NewReflector(&openapi.Config{OpenAPIVersion: openapi.Version312})
		s312, err := r312.SchemaForType(std_reflect.TypeFor[*User](), reflect.SchemaUseComponent, nil)
		require.NoError(t, err)
		require.Len(t, s312.AnyOf, 2)
		assert.Equal(t, "#/components/schemas/ReflectTestUser", s312.AnyOf[0].Ref)
		assert.Equal(t, "null", s312.AnyOf[1].Type)
	})
}

func TestConverter_FileTypes(t *testing.T) {
	type multiFile struct {
		Files []*multipart.FileHeader `form:"files" description:"Attachments"`
	}

	versions := []string{openapi.Version304, openapi.Version312, openapi.Version320}
	for _, version := range versions {
		t.Run(version, func(t *testing.T) {
			r := reflect.NewReflector(&openapi.Config{OpenAPIVersion: version})

			for name, typ := range map[string]std_reflect.Type{
				"FileHeader":        std_reflect.TypeFor[multipart.FileHeader](),
				"FileHeaderPointer": std_reflect.TypeFor[*multipart.FileHeader](),
				"File":              std_reflect.TypeFor[multipart.File](),
			} {
				t.Run(name, func(t *testing.T) {
					schema, err := r.SchemaForType(typ, reflect.SchemaUseComponent, nil)
					require.NoError(t, err)
					assert.Equal(t, &openapi.Schema{Type: "string", Format: "binary"}, schema)
					assert.Empty(t, r.Components)
				})
			}

			t.Run("SliceOfFileHeaders", func(t *testing.T) {
				field, ok := std_reflect.TypeFor[multiFile]().FieldByName("Files")
				require.True(t, ok)
				schema, err := r.SchemaForType(field.Type, reflect.SchemaUseComponent, &field)
				require.NoError(t, err)
				assert.Equal(t, "array", schema.Type)
				assert.Equal(t, "Attachments", schema.Description)
				require.NotNil(t, schema.Items)
				assert.Equal(t, &openapi.Schema{Type: "string", Format: "binary"}, schema.Items)
				assert.Empty(t, r.Components)
			})
		})
	}

	t.Run("TypeMappingTakesPrecedence", func(t *testing.T) {
		r := reflect.NewReflector(&openapi.Config{
			OpenAPIVersion: openapi.Version312,
			ReflectorConfig: &openapi.ReflectorConfig{
				TypeMappings: []openapi.TypeMapping{{Src: multipart.FileHeader{}, Dst: 0}},
			},
		})
		schema, err := r.SchemaForType(std_reflect.TypeFor[multipart.FileHeader](), reflect.SchemaInline, nil)
		require.NoError(t, err)
		assert.Equal(t, "integer", schema.Type)
	})
}
