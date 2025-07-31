package option_test

import (
	"reflect"
	"testing"

	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
	"github.com/swaggest/jsonschema-go"
)

func TestInlineRefs(t *testing.T) {
	config := &openapi.ReflectorConfig{}
	opt := option.InlineRefs()
	opt(config)

	assert.True(t, config.InlineRefs)
}

func TestRootRef(t *testing.T) {
	config := &openapi.ReflectorConfig{}
	opt := option.RootRef()
	opt(config)

	assert.True(t, config.RootRef)
}

func TestRootNullable(t *testing.T) {
	config := &openapi.ReflectorConfig{}
	opt := option.RootNullable()
	opt(config)

	assert.True(t, config.RootNullable)
}

func TestStripDefNamePrefix(t *testing.T) {
	config := &openapi.ReflectorConfig{}
	prefixes := []string{"Test", "Mock"}
	opt := option.StripDefNamePrefix(prefixes...)
	opt(config)

	assert.Equal(t, prefixes, config.StripDefNamePrefix)
}

func TestInterceptDefNameFunc(t *testing.T) {
	config := &openapi.ReflectorConfig{}
	mockFunc := func(t reflect.Type, defaultDefName string) string {
		return "CustomName"
	}
	opt := option.InterceptDefNameFunc(mockFunc)
	opt(config)

	assert.NotNil(t, config.InterceptDefNameFunc)
}

func TestInterceptPropFunc(t *testing.T) {
	config := &openapi.ReflectorConfig{}
	mockFunc := func(params openapi.InterceptPropParams) error {
		return nil
	}
	opt := option.InterceptPropFunc(mockFunc)
	opt(config)

	assert.NotNil(t, config.InterceptPropFunc)
}

func TestRequiredPropByValidateTag(t *testing.T) {
	config := &openapi.ReflectorConfig{}
	opt := option.RequiredPropByValidateTag()
	opt(config)

	assert.NotNil(t, config.InterceptPropFunc)

	params := openapi.InterceptPropParams{
		Name: "Field1",
		Field: reflect.StructField{
			Tag: reflect.StructTag(`validate:"required"`),
		},
		ParentSchema: &jsonschema.Schema{
			Required: []string{},
		},
		Processed: true,
	}
	err := config.InterceptPropFunc(params)
	assert.NoError(t, err)
	assert.Contains(t, params.ParentSchema.Required, "Field1")
}

func TestInterceptSchemaFunc(t *testing.T) {
	config := &openapi.ReflectorConfig{}
	mockFunc := func(params openapi.InterceptSchemaParams) (stop bool, err error) {
		return false, nil
	}
	opt := option.InterceptSchemaFunc(mockFunc)
	opt(config)

	assert.NotNil(t, config.InterceptSchemaFunc)
}

func TestTypeMapping(t *testing.T) {
	config := &openapi.ReflectorConfig{}
	src := "source"
	dst := "destination"
	opt := option.TypeMapping(src, dst)
	opt(config)

	assert.Len(t, config.TypeMappings, 1)

	mapping := config.TypeMappings[0]
	assert.Equal(t, src, mapping.Src)
	assert.Equal(t, dst, mapping.Dst)
}
