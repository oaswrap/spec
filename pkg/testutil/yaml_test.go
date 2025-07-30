package testutil_test

import (
	"testing"

	"github.com/oaswrap/spec/pkg/testutil"
)

func TestEqualYAML(t *testing.T) {
	tests := []struct {
		name      string
		want      []byte
		got       []byte
		shouldErr bool
	}{
		{
			name:      "identical YAML",
			want:      []byte("key: value\nother: 123"),
			got:       []byte("key: value\nother: 123"),
			shouldErr: false,
		},
		{
			name:      "semantically equal YAML different formatting",
			want:      []byte("key: value\nother: 123"),
			got:       []byte("other: 123\nkey: value"),
			shouldErr: false,
		},
		{
			name:      "different values",
			want:      []byte("key: value1"),
			got:       []byte("key: value2"),
			shouldErr: true,
		},
		{
			name:      "different keys",
			want:      []byte("key1: value"),
			got:       []byte("key2: value"),
			shouldErr: true,
		},
		{
			name:      "nested objects equal",
			want:      []byte("parent:\n  child: value\n  num: 42"),
			got:       []byte("parent:\n  num: 42\n  child: value"),
			shouldErr: false,
		},
		{
			name:      "arrays equal",
			want:      []byte("items:\n  - one\n  - two\n  - three"),
			got:       []byte("items:\n  - one\n  - two\n  - three"),
			shouldErr: false,
		},
		{
			name:      "arrays different order",
			want:      []byte("items:\n  - one\n  - two"),
			got:       []byte("items:\n  - two\n  - one"),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockT := &testing.T{}
			testutil.EqualYAML(mockT, tt.want, tt.got)

			if tt.shouldErr && !mockT.Failed() {
				t.Error("Expected test to fail but it passed")
			}
			if !tt.shouldErr && mockT.Failed() {
				t.Error("Expected test to pass but it failed")
			}
		})
	}
}
