package option_test

import (
	"testing"

	"github.com/oaswrap/spec/openapi"
	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestContentOption(t *testing.T) {
	tests := []struct {
		name       string
		httpStatus int
		opts       []option.ContentOption
		expected   openapi.ContentUnit
	}{
		{
			name:       "empty options",
			httpStatus: 0,
			opts:       []option.ContentOption{},
			expected: openapi.ContentUnit{
				HTTPStatus: 0,
			},
		},
		{
			name:       "with content type",
			httpStatus: 200,
			opts: []option.ContentOption{
				option.WithContentType("application/json"),
			},
			expected: openapi.ContentUnit{
				HTTPStatus:  200,
				ContentType: "application/json",
			},
		},
		{
			name:       "with description",
			httpStatus: 200,
			opts: []option.ContentOption{
				option.WithContentDescription("This is a response"),
			},
			expected: openapi.ContentUnit{
				HTTPStatus:  200,
				Description: "This is a response",
			},
		},
		{
			name:       "with default flag",
			httpStatus: 200,
			opts: []option.ContentOption{
				option.WithContentDefault(true),
			},
			expected: openapi.ContentUnit{
				HTTPStatus: 200,
				IsDefault:  true,
			},
		},
		{
			name:       "with multiple options",
			httpStatus: 200,
			opts: []option.ContentOption{
				option.WithContentType("application/json"),
				option.WithContentDescription("This is a response"),
				option.WithContentDefault(true),
			},
			expected: openapi.ContentUnit{
				HTTPStatus:  200,
				ContentType: "application/json",
				Description: "This is a response",
				IsDefault:   true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &openapi.ContentUnit{
				HTTPStatus: tt.httpStatus,
			}
			for _, opt := range tt.opts {
				opt(config)
			}
			assert.Equal(t, tt.expected, *config)
		})
	}
}
