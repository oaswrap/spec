package validate_test

import (
	"reflect"
	"testing"

	"github.com/oaswrap/spec/internal/validate"

	"github.com/stretchr/testify/assert"

	"github.com/oaswrap/spec/openapi"
)

func TestNormalizeTemplatedPath(t *testing.T) {
	assert.Equal(t, "/users/{}", validate.NormalizeTemplatedPath("/users/{id}"))
	assert.Equal(t, "/orgs/{}/repos/{}", validate.NormalizeTemplatedPath("/orgs/{org}/repos/{repo}"))
	assert.Equal(t, "/type/{}", validate.NormalizeTemplatedPath("/type/{type:foo|bar}"))
	assert.Equal(t, "/more/{}/{}", validate.NormalizeTemplatedPath("/more/{param}/{type:foo|bar}"))
	assert.Equal(t, "/complex/quantifier/{}", validate.NormalizeTemplatedPath("/complex/quantifier/{t:(\\d{4})}"))
	assert.Equal(
		t,
		"/foo/bar/{}/two/{}/8/nine/{}",
		validate.NormalizeTemplatedPath("/foo/bar/{one}/two/{three:(four|five){6,7}(eight|nine)}/8/nine/{wtf}"),
	)
	assert.Equal(t, "/static", validate.NormalizeTemplatedPath("/static"))
}

func TestMediaTypeBase(t *testing.T) {
	assert.Equal(t, "application/json", validate.MediaTypeBase("application/json"))
	assert.Equal(t, "application/json", validate.MediaTypeBase("application/json; charset=utf-8"))
	assert.Equal(t, "application/json", validate.MediaTypeBase("  APPLICATION/JSON  ; foo=bar"))
}

func TestMediaTypeIsMultipart(t *testing.T) {
	assert.True(t, validate.MediaTypeIsMultipart("multipart/form-data"))
	assert.True(t, validate.MediaTypeIsMultipart("multipart/mixed"))
	assert.False(t, validate.MediaTypeIsMultipart("application/json"))
}

func TestResolveJSONPointer(t *testing.T) {
	root := map[string]any{
		"foo": []any{"bar", "baz"},
		"qux": map[string]any{
			"a/b": 1,
			"c%d": 2,
			"e~f": 3,
			"g~h": 4,
		},
	}

	assert.Equal(t, root, validate.ResolveJSONPointer(root, ""))
	assert.Equal(t, root["foo"], validate.ResolveJSONPointer(root, "/foo"))
	assert.Equal(t, "bar", validate.ResolveJSONPointer(root, "/foo/0"))
	assert.Equal(t, "baz", validate.ResolveJSONPointer(root, "/foo/1"))
	assert.Nil(t, validate.ResolveJSONPointer(root, "/foo/2"))
	assert.Equal(t, 1, validate.ResolveJSONPointer(root, "/qux/a~1b"))
	assert.Equal(t, 3, validate.ResolveJSONPointer(root, "/qux/e~0f"))
}

func TestIsNonRelativeURI(t *testing.T) {
	assert.True(t, validate.IsNonRelativeURI("https://example.com"))
	assert.True(t, validate.IsNonRelativeURI("https://example.com#frag"))
	assert.True(t, validate.IsNonRelativeURI("mailto:foo@example.com"))
	assert.False(t, validate.IsNonRelativeURI("/local/path"))
	assert.False(t, validate.IsNonRelativeURI("relative"))
}

func TestIsHTTPSURI(t *testing.T) {
	assert.True(t, validate.IsHTTPSURI("https://example.com"))
	assert.False(t, validate.IsHTTPSURI("http://example.com"))
	assert.False(t, validate.IsHTTPSURI("ftp://example.com"))
}

func TestIsURIReference(t *testing.T) {
	assert.True(t, validate.IsURIReference("https://example.com"))
	assert.True(t, validate.IsURIReference("/path"))
	assert.False(t, validate.IsURIReference("not a uri with spaces"))
}

func TestResolveURIReference(t *testing.T) {
	tests := []struct {
		base, ref string
		expected  string
	}{
		{"", "https://example.com", "https://example.com"},
		{"https://example.com/a/b", "c", "https://example.com/a/c"},
		{"https://example.com/a/b", "/c", "https://example.com/c"},
	}
	for _, tt := range tests {
		got, ok := validate.ResolveURIReference(tt.base, tt.ref)
		assert.True(t, ok)
		assert.Equal(t, tt.expected, got)
	}
}

func TestExtraHas(t *testing.T) {
	extra := map[string]any{"foo": 1, "bar": 2}
	assert.True(t, validate.ExtraHas(extra, "foo"))
	assert.True(t, validate.ExtraHas(extra, "baz", "bar"))
	assert.False(t, validate.ExtraHas(extra, "qux"))
}

func TestWithoutFragment(t *testing.T) {
	assert.Equal(t, "https://example.com/path", validate.WithoutFragment("https://example.com/path#frag"))
	assert.Equal(t, "https://example.com/path", validate.WithoutFragment("https://example.com/path"))
	assert.Equal(t, ":// bad", validate.WithoutFragment(":// bad"))
}

func TestMediaTypeAllowsNamedEncoding(t *testing.T) {
	assert.True(t, validate.MediaTypeAllowsNamedEncoding(""))
	assert.True(t, validate.MediaTypeAllowsNamedEncoding("multipart/form-data"))
	assert.True(t, validate.MediaTypeAllowsNamedEncoding("application/x-www-form-urlencoded"))
	assert.False(t, validate.MediaTypeAllowsNamedEncoding("application/json"))
}

func TestResolveURIReference_InvalidInput(t *testing.T) {
	_, ok := validate.ResolveURIReference("https://example.com", "://bad")
	assert.False(t, ok)

	_, ok = validate.ResolveURIReference("://bad", "x")
	assert.False(t, ok)
}

func TestRegisterSchemaResourceAndAnchor(t *testing.T) {
	schema := &openapi.Schema{
		ID:            "https://schemas.example.com/user",
		Type:          "object",
		Anchor:        "user-anchor",
		DynamicAnchor: "user-dyn",
	}

	resources := map[string]any{}
	validate.RegisterSchemaResource(reflect.ValueOf(*schema), schema.ID, resources)
	validate.RegisterSchemaResource(reflect.ValueOf(*schema), schema.ID, resources) // idempotent

	assert.Contains(t, resources, "https://schemas.example.com/user")
	assert.Contains(t, resources, "https://schemas.example.com/user#user-anchor")
	assert.Contains(t, resources, "https://schemas.example.com/user#user-dyn")
	assert.Len(t, resources, 3)

	// Empty anchor should be ignored.
	emptySchema := openapi.Schema{}
	validate.RegisterSchemaAnchor(reflect.ValueOf(emptySchema), "", "Anchor", resources)
	assert.Len(t, resources, 3)
}

func TestIsLocalReferenceAndReferenceTargetExists(t *testing.T) {
	resources := map[string]any{
		"https://example.com/schemas/user": map[string]any{
			"properties": map[string]any{
				"name": map[string]any{"type": "string"},
			},
		},
		"https://example.com/schemas/user#UserAnchor": map[string]any{"type": "object"},
		"": map[string]any{
			"components": map[string]any{
				"schemas": map[string]any{
					"User": map[string]any{"type": "object"},
				},
			},
		},
	}

	assert.True(t, validate.IsLocalReference("https://example.com/schemas/user", resources))
	assert.True(t, validate.IsLocalReference("https://example.com/schemas/user#UserAnchor", resources))
	assert.False(t, validate.IsLocalReference("https://example.com/schemas/user#Missing", resources))
	assert.False(t, validate.IsLocalReference("://bad", resources))

	assert.True(t, validate.ReferenceTargetExists("https://example.com/schemas/user", resources))
	assert.True(t, validate.ReferenceTargetExists("https://example.com/schemas/user#UserAnchor", resources))
	assert.True(t, validate.ReferenceTargetExists("https://example.com/schemas/user#/properties/name", resources))
	assert.False(t, validate.ReferenceTargetExists("https://example.com/schemas/user#/properties/missing", resources))
	assert.False(t, validate.ReferenceTargetExists("https://example.com/schemas/user#Missing", resources))
	assert.False(t, validate.ReferenceTargetExists("://bad", resources))
}

func TestMarshalAny(t *testing.T) {
	assert.Equal(t, map[string]any{"x": float64(1)}, validate.MarshalAny(map[string]int{"x": 1}))
	assert.Nil(t, validate.MarshalAny(func() {}))
}
