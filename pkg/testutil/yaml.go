package testutil

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

// YAMLToInterface parses a YAML blob into an interface{}.
// Use it before comparing.
func YAMLToInterface(t *testing.T, data []byte) interface{} {
	t.Helper()
	var v interface{}
	dec := yaml.NewDecoder(bytes.NewReader(data))
	err := dec.Decode(&v)
	assert.NoError(t, err)
	return v
}

// EqualYAML asserts that two YAML documents are semantically equal.
// It returns a cmp.Diff if they are not.
func EqualYAML(t *testing.T, want []byte, got []byte) {
	wantObj := YAMLToInterface(t, want)
	gotObj := YAMLToInterface(t, got)

	if diff := cmp.Diff(wantObj, gotObj); diff != "" {
		t.Errorf("YAML mismatch (-want +got):\n%s", diff)
	}
}
