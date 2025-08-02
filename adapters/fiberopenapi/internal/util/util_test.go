package util_test

import (
	"testing"

	"github.com/oaswrap/spec/adapters/fiberopenapi/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestConvertPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single parameter",
			input:    "/users/:id",
			expected: "/users/{id}",
		},
		{
			name:     "multiple parameters",
			input:    "/users/:id/posts/:postId",
			expected: "/users/{id}/posts/{postId}",
		},
		{
			name:     "no parameters",
			input:    "/users/all",
			expected: "/users/all",
		},
		{
			name:     "parameter with numbers",
			input:    "/items/:item123",
			expected: "/items/{item123}",
		},
		{
			name:     "parameter with underscore",
			input:    "/data/:user_id",
			expected: "/data/{user_id}",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "parameter at start",
			input:    ":id/details",
			expected: "{id}/details",
		},
		{
			name:     "multiple consecutive parameters",
			input:    "/:id/:name/:type",
			expected: "/{id}/{name}/{type}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := util.ConvertPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
