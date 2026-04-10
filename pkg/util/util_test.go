package util_test

import (
	"testing"

	"github.com/oaswrap/spec/pkg/util"
	"github.com/stretchr/testify/assert"
)

func TestOptional(t *testing.T) {
	t.Run("returns default value when no optional value provided", func(t *testing.T) {
		result := util.Optional("default")
		assert.Equal(t, "default", result)
	})

	t.Run("returns first optional value when provided", func(t *testing.T) {
		result := util.Optional("default", "provided")
		assert.Equal(t, "provided", result)
	})

	t.Run("returns first optional value when multiple values provided", func(t *testing.T) {
		result := util.Optional("default", "first", "second", "third")
		assert.Equal(t, "first", result)
	})

	t.Run("returns default when empty slice provided", func(t *testing.T) {
		var values []string
		result := util.Optional("default", values...)
		assert.Equal(t, "default", result)
	})
}

func TestPtrOf(t *testing.T) {
	t.Run("returns pointer to string value", func(t *testing.T) {
		value := "test"
		result := util.PtrOf(value)
		assert.NotNil(t, result)
		assert.Equal(t, value, *result)
	})

	t.Run("returns pointer to int value", func(t *testing.T) {
		value := 42
		result := util.PtrOf(value)
		assert.NotNil(t, result)
		assert.Equal(t, value, *result)
	})

	t.Run("returns pointer to bool value", func(t *testing.T) {
		value := true
		result := util.PtrOf(value)
		assert.NotNil(t, result)
		assert.Equal(t, value, *result)
	})

	t.Run("returns pointer to zero value", func(t *testing.T) {
		value := 0
		result := util.PtrOf(value)
		assert.NotNil(t, result)
		assert.Equal(t, value, *result)
	})

	t.Run("returns pointer to struct", func(t *testing.T) {
		type testStruct struct {
			Field string
		}
		value := testStruct{Field: "test"}
		result := util.PtrOf(value)
		assert.NotNil(t, result)
		assert.Equal(t, value, *result)
	})
}

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		expected string
	}{
		{
			name:     "empty paths",
			paths:    []string{},
			expected: "",
		},
		{
			name:     "single path without trailing slash",
			paths:    []string{"api"},
			expected: "api",
		},
		{
			name:     "single path with trailing slash",
			paths:    []string{"api/"},
			expected: "api/",
		},
		{
			name:     "two paths without trailing slash",
			paths:    []string{"api", "v1"},
			expected: "api/v1",
		},
		{
			name:     "two paths with trailing slash on last",
			paths:    []string{"api", "v1/"},
			expected: "api/v1/",
		},
		{
			name:     "multiple paths without trailing slash",
			paths:    []string{"api", "v1", "users"},
			expected: "api/v1/users",
		},
		{
			name:     "multiple paths with trailing slash on last",
			paths:    []string{"api", "v1", "users/"},
			expected: "api/v1/users/",
		},
		{
			name:     "absolute paths",
			paths:    []string{"/api", "v1"},
			expected: "/api/v1",
		},
		{
			name:     "absolute paths with trailing slash",
			paths:    []string{"/api", "v1/"},
			expected: "/api/v1/",
		},
		{
			name:     "normalize double slashes",
			paths:    []string{"api/", "/v1"},
			expected: "api/v1",
		},
		{
			name:     "normalize double slashes with trailing slash",
			paths:    []string{"api/", "/v1/"},
			expected: "api/v1/",
		},
		{
			name:     "empty string in middle paths",
			paths:    []string{"api", "", "v1"},
			expected: "api/v1",
		},
		{
			name:     "empty string on last path",
			paths:    []string{"api", "v1", ""},
			expected: "api/v1",
		},
		{
			name:     "many path segments",
			paths:    []string{"api", "v1", "users", "123", "profile"},
			expected: "api/v1/users/123/profile",
		},
		{
			name:     "many path segments with trailing slash",
			paths:    []string{"api", "v1", "users", "123", "profile/"},
			expected: "api/v1/users/123/profile/",
		},
		{
			name:     "dot paths",
			paths:    []string{".", "api", "v1"},
			expected: "api/v1",
		},
		{
			name:     "parent directory paths",
			paths:    []string{"api", "..", "v1"},
			expected: "v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.JoinPath(tt.paths...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
