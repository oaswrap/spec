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

func TestJoinURL(t *testing.T) {
	tests := []struct {
		name     string
		base     string
		segments []string
		expected string
	}{
		{
			name:     "base without trailing slash",
			base:     "https://example.com",
			segments: []string{"api", "v1", "users"},
			expected: "https://example.com/api/v1/users",
		},
		{
			name:     "base with trailing slash",
			base:     "https://example.com/",
			segments: []string{"api", "v1", "users"},
			expected: "https://example.com/api/v1/users",
		},
		{
			name:     "base with multiple trailing slashes",
			base:     "https://example.com///",
			segments: []string{"api", "v1", "users"},
			expected: "https://example.com/api/v1/users",
		},
		{
			name:     "empty segments",
			base:     "https://example.com",
			segments: []string{},
			expected: "https://example.com",
		},
		{
			name:     "single segment",
			base:     "https://example.com",
			segments: []string{"api"},
			expected: "https://example.com/api",
		},
		{
			name:     "segments with slashes",
			base:     "https://example.com",
			segments: []string{"api/v1", "users"},
			expected: "https://example.com/api/v1/users",
		},
		{
			name:     "segments with leading slashes",
			base:     "https://example.com",
			segments: []string{"/api/v1", "/users"},
			expected: "https://example.com/api/v1/users",
		},
		{
			name:     "empty base",
			base:     "",
			segments: []string{"api", "v1"},
			expected: "/api/v1",
		},
		{
			name:     "base with only slashes",
			base:     "///",
			segments: []string{"api", "v1"},
			expected: "/api/v1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.JoinURL(tt.base, tt.segments...)
			assert.Equal(t, tt.expected, result)
		})
	}
}
