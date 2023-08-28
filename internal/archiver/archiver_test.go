package archiver

import (
	"archive/zip"
	"context"
	"fmt"
	"os"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPackageFile(t *testing.T) {

	t.Run("packages regular file", func(t *testing.T) {
		// Create temp file
		tmpFile, err := os.CreateTemp("", "test")
		require.NoError(t, err)

		fmt.Println("file", tmpFile.Name())

		// Create temp archive file
		archiveFile, err := os.CreateTemp("", "test*.zip")
		require.NoError(t, err)

		fmt.Println("archive", archiveFile.Name())

		// Create template
		tmpl, err := template.New("bootstrap").Parse("hello {{.}}")
		require.NoError(t, err)

		// Call function
		err = PackageFile(context.Background(), tmpl, tmpFile.Name(), archiveFile.Name())
		require.NoError(t, err)

		// Validate archive file created
		r, err := zip.OpenReader(archiveFile.Name())
		require.NoError(t, err)

		// Validate entries in archive
		assert.Equal(t, 2, len(r.File))

		// Clean up
		require.NoError(t, r.Close())
		require.NoError(t, os.Remove(archiveFile.Name()))
	})

}
