package binaries

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGlob(t *testing.T) {

	t.Run("returns matched files", func(t *testing.T) {

		temDir, _ := os.MkdirTemp("", "glob")

		// Create some temp files
		tmpFile1, _ := os.CreateTemp(temDir, "test1")
		tmpFile2, _ := os.CreateTemp(temDir, "test2")

		// Call function
		files, err := Glob(filepath.Join(temDir, "test*"))

		// Validate results
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{tmpFile1.Name(), tmpFile2.Name()}, files)
	})

	t.Run("skips directories", func(t *testing.T) {
		// Create temp dir
		tmpDir, _ := os.MkdirTemp("", "testdir")

		// Call function
		files, err := Glob(tmpDir)

		// Validate results
		require.NoError(t, err)
		assert.Empty(t, files)
	})

	// Add tests for other cases...
}
