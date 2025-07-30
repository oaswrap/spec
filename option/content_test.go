package option_test

import (
	"testing"

	"github.com/oaswrap/spec/option"
	"github.com/stretchr/testify/assert"
)

func TestWithContentType(t *testing.T) {
	tests := []struct {
		name        string
		contentType string
		want        string
	}{
		{
			name:        "sets application/json content type",
			contentType: "application/json",
			want:        "application/json",
		},
		{
			name:        "sets text/plain content type",
			contentType: "text/plain",
			want:        "text/plain",
		},
		{
			name:        "sets empty content type",
			contentType: "",
			want:        "",
		},
		{
			name:        "sets application/xml content type",
			contentType: "application/xml",
			want:        "application/xml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &option.ContentConfig{}
			opt := option.WithContentType(tt.contentType)
			opt(config)

			assert.Equal(t, tt.want, config.ContentType)
		})
	}
}

func TestContentOption(t *testing.T) {
	t.Run("multiple options can be applied", func(t *testing.T) {
		config := &option.ContentConfig{
			HTTPStatus: 200,
		}

		opt1 := option.WithContentType("application/json")
		opt2 := option.WithContentType("text/plain")

		opt1(config)
		assert.Equal(t, "application/json", config.ContentType)

		opt2(config)
		assert.Equal(t, "text/plain", config.ContentType)

		assert.Equal(t, 200, config.HTTPStatus)
	})
}
