package parser_test

import (
	"testing"

	"github.com/oaswrap/spec/adapters/httpopenapi/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestParseRoutePattern(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *parser.RoutePattern
		wantErr bool
	}{
		{
			name:  "path only",
			input: "/api/users",
			want: &parser.RoutePattern{
				Method: "",
				Host:   "",
				Path:   "/api/users",
			},
			wantErr: false,
		},
		{
			name:  "method and path",
			input: "GET /api/users",
			want: &parser.RoutePattern{
				Method: "GET",
				Host:   "",
				Path:   "/api/users",
			},
			wantErr: false,
		},
		{
			name:  "host and path",
			input: "example.com/api/users",
			want: &parser.RoutePattern{
				Method: "",
				Host:   "example.com",
				Path:   "/api/users",
			},
			wantErr: false,
		},
		{
			name:  "method, host and path",
			input: "POST api.example.com/users",
			want: &parser.RoutePattern{
				Method: "POST",
				Host:   "api.example.com",
				Path:   "/users",
			},
			wantErr: false,
		},
		{
			name:  "host only",
			input: "example.com",
			want: &parser.RoutePattern{
				Method: "",
				Host:   "example.com",
				Path:   "/",
			},
			wantErr: false,
		},
		{
			name:  "localhost",
			input: "localhost/api",
			want: &parser.RoutePattern{
				Method: "",
				Host:   "localhost",
				Path:   "/api",
			},
			wantErr: false,
		},
		{
			name:  "host with port",
			input: "localhost:8080/api",
			want: &parser.RoutePattern{
				Method: "",
				Host:   "localhost:8080",
				Path:   "/api",
			},
			wantErr: false,
		},
		{
			name:    "invalid host with space",
			input:   "invalid host/path",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid hostname",
			input:   "invalidhost/path",
			want:    nil,
			wantErr: true,
		},
		{
			name:  "valid single word hostname",
			input: "localhost/path",
			want: &parser.RoutePattern{
				Method: "",
				Host:   "localhost",
				Path:   "/path",
			},
			wantErr: false,
		},
		{
			name:  "all HTTP methods",
			input: "DELETE api.example.com/resource",
			want: &parser.RoutePattern{
				Method: "DELETE",
				Host:   "api.example.com",
				Path:   "/resource",
			},
			wantErr: false,
		},
		{
			name:    "non-HTTP method word",
			input:   "INVALID api.example.com/resource",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parser.ParseRoutePattern(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.Method, got.Method)
				assert.Equal(t, tt.want.Host, got.Host)
				assert.Equal(t, tt.want.Path, got.Path)
			}
		})
	}
}
