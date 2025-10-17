package parser_test

import (
	"testing"

	"github.com/oaswrap/spec/pkg/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewColonParamParser(t *testing.T) {
	parser := parser.NewColonParamParser()
	assert.NotNil(t, parser, "NewColonParamParser() returned nil")
}

func TestColonParamParser_Parse(t *testing.T) {
	parser := parser.NewColonParamParser()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple parameter",
			input:    "/users/:id",
			expected: "/users/{id}",
		},
		{
			name:     "multiple parameters",
			input:    "/users/:userId/posts/:postId",
			expected: "/users/{userId}/posts/{postId}",
		},
		{
			name:     "no parameters",
			input:    "/users",
			expected: "/users",
		},
		{
			name:     "parameter with underscore",
			input:    "/users/:user_id",
			expected: "/users/{user_id}",
		},
		{
			name:     "parameter with numbers",
			input:    "/api/v1/:id123",
			expected: "/api/v1/{id123}",
		},
		{
			name:     "parameter at root",
			input:    "/:id",
			expected: "/{id}",
		},
		{
			name:     "mixed parameters and static segments",
			input:    "/api/:version/users/:id/profile",
			expected: "/api/{version}/users/{id}/profile",
		},
		{
			name:     "empty path",
			input:    "",
			expected: "",
		},
		{
			name:     "parameter starting with underscore",
			input:    "/users/:_id",
			expected: "/users/{_id}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parser.Parse(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
