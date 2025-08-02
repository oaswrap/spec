package option_test

import (
	"testing"

	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestHidden(t *testing.T) {
	t.Run("default hidden true", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Hidden()(cfg)
		assert.True(t, cfg.Hide)
	})

	t.Run("explicit hidden true", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Hidden(true)(cfg)
		assert.True(t, cfg.Hide)
	})

	t.Run("explicit hidden false", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Hidden(false)(cfg)
		assert.False(t, cfg.Hide)
	})
}

func TestOperationID(t *testing.T) {
	cfg := &option.OperationConfig{}
	option.OperationID("test-operation")(cfg)
	assert.Equal(t, "test-operation", cfg.OperationID)
}

func TestDescription(t *testing.T) {
	cfg := &option.OperationConfig{}
	option.Description("Test description")(cfg)
	assert.Equal(t, "Test description", cfg.Description)
}

func TestSummary(t *testing.T) {
	t.Run("summary only", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Summary("Test summary")(cfg)
		assert.Equal(t, "Test summary", cfg.Summary)
		assert.Equal(t, "Test summary", cfg.Description)
	})

	t.Run("summary with existing description", func(t *testing.T) {
		cfg := &option.OperationConfig{Description: "Existing description"}
		option.Summary("Test summary")(cfg)
		assert.Equal(t, "Test summary", cfg.Summary)
		assert.Equal(t, "Existing description", cfg.Description)
	})
}

func TestDeprecated(t *testing.T) {
	t.Run("default deprecated true", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Deprecated()(cfg)
		assert.True(t, cfg.Deprecated)
	})

	t.Run("explicit deprecated true", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Deprecated(true)(cfg)
		assert.True(t, cfg.Deprecated)
	})

	t.Run("explicit deprecated false", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Deprecated(false)(cfg)
		assert.False(t, cfg.Deprecated)
	})
}

func TestTags(t *testing.T) {
	t.Run("single tag", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Tags("auth")(cfg)
		assert.Equal(t, []string{"auth"}, cfg.Tags)
	})

	t.Run("multiple tags", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Tags("auth", "users", "admin")(cfg)
		assert.Equal(t, []string{"auth", "users", "admin"}, cfg.Tags)
	})

	t.Run("append tags", func(t *testing.T) {
		cfg := &option.OperationConfig{Tags: []string{"existing"}}
		option.Tags("new")(cfg)
		assert.Equal(t, []string{"existing", "new"}, cfg.Tags)
	})
}

func TestSecurity(t *testing.T) {
	t.Run("security without scopes", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Security("bearer")(cfg)
		assert.Len(t, cfg.Security, 1)
		assert.Equal(t, "bearer", cfg.Security[0].Name)
		assert.Empty(t, cfg.Security[0].Scopes)
	})

	t.Run("security with scopes", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Security("oauth2", "read", "write")(cfg)
		assert.Len(t, cfg.Security, 1)
		assert.Equal(t, "oauth2", cfg.Security[0].Name)
		assert.Equal(t, []string{"read", "write"}, cfg.Security[0].Scopes)
	})

	t.Run("multiple security configs", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Security("bearer")(cfg)
		option.Security("oauth2", "read")(cfg)
		assert.Len(t, cfg.Security, 2)
		assert.Equal(t, "bearer", cfg.Security[0].Name)
		assert.Equal(t, "oauth2", cfg.Security[1].Name)
	})
}

func TestRequest(t *testing.T) {
	type TestStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	t.Run("request without options", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Request(TestStruct{})(cfg)
		assert.Len(t, cfg.Requests, 1)
		assert.Equal(t, TestStruct{}, cfg.Requests[0].Structure)
	})

	t.Run("multiple requests", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Request(TestStruct{})(cfg)
		option.Request("string")(cfg)
		assert.Len(t, cfg.Requests, 2)
		assert.Equal(t, TestStruct{}, cfg.Requests[0].Structure)
		assert.Equal(t, "string", cfg.Requests[1].Structure)
	})

	t.Run("request with content options", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Request(TestStruct{}, option.WithContentType("application/json"))(cfg)
		assert.Len(t, cfg.Requests, 1)
		assert.Equal(t, "application/json", cfg.Requests[0].ContentType)
		assert.Equal(t, TestStruct{}, cfg.Requests[0].Structure)
	})
}

func TestResponse(t *testing.T) {
	type TestStruct struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	t.Run("response without options", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Response(200, TestStruct{})(cfg)
		assert.Len(t, cfg.Responses, 1)
		assert.Equal(t, 200, cfg.Responses[0].HTTPStatus)
		assert.Equal(t, TestStruct{}, cfg.Responses[0].Structure)
	})

	t.Run("multiple responses", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Response(200, TestStruct{})(cfg)
		option.Response(400, "error")(cfg)
		assert.Len(t, cfg.Responses, 2)
		assert.Equal(t, 200, cfg.Responses[0].HTTPStatus)
		assert.Equal(t, 400, cfg.Responses[1].HTTPStatus)
	})

	t.Run("response with content options", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		option.Response(200, TestStruct{}, option.WithContentType("application/json"))(cfg)
		assert.Len(t, cfg.Responses, 1)
		assert.Equal(t, 200, cfg.Responses[0].HTTPStatus)
		assert.Equal(t, "application/json", cfg.Responses[0].ContentType)
		assert.Equal(t, TestStruct{}, cfg.Responses[0].Structure)
	})
}

func TestOperationConfig(t *testing.T) {
	t.Run("default values", func(t *testing.T) {
		cfg := &option.OperationConfig{}
		assert.False(t, cfg.Hide)
		assert.Empty(t, cfg.OperationID)
		assert.Empty(t, cfg.Description)
		assert.Empty(t, cfg.Summary)
		assert.False(t, cfg.Deprecated)
		assert.Nil(t, cfg.Tags)
		assert.Nil(t, cfg.Security)
		assert.Nil(t, cfg.Requests)
		assert.Nil(t, cfg.Responses)
	})
}
