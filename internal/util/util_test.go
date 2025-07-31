package util_test

import (
	"github.com/oaswrap/spec/internal/util"
	"testing"

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
