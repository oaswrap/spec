package testutil

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
)

// Update is a test flag for updating golden files.
var Update = flag.Bool("update", false, "update golden files")

// AssertGolden compares the generated schema with a golden file.
func AssertGolden(t *testing.T, schema []byte, goldenFile string) {
	t.Helper()

	if *Update {
		err := os.MkdirAll(filepath.Dir(goldenFile), 0750)
		require.NoError(t, err, "failed to create golden file directory")
		err = os.WriteFile(goldenFile, schema, 0600)
		require.NoError(t, err, "failed to write golden file")
		t.Logf("Updated golden file: %s", goldenFile)
	}

	// #nosec G304 -- goldenFile is a test-controlled fixture path (always from testdata).
	want, err := os.ReadFile(goldenFile)
	require.NoError(t, err, "failed to read golden file %s", goldenFile)

	diff := cmp.Diff(string(want), string(schema))
	if diff != "" {
		t.Errorf("OpenAPI schema mismatch (-want +got):\n%s", diff)
	}
}
