package archiver

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/rs/zerolog"
)

func PackageFile(ctx context.Context, bootstrapTemplate *template.Template, binaryFile, archiveFile string) error {

	binaryName := filepath.Base(binaryFile)

	zerolog.Ctx(ctx).Debug().Str("name", binaryName).Str("archiveFile", archiveFile).Msg("Processing file")

	archive, err := os.Create(archiveFile)
	if err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}

	defer archive.Close()

	zipWriter := zip.NewWriter(archive)

	binf, err := os.Open(binaryFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}

	defer binf.Close()

	binw, err := zipWriter.Create(binaryName)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	if _, err := io.Copy(binw, binf); err != nil {
		return fmt.Errorf("failed to write file to archive: %w", err)
	}

	zerolog.Ctx(ctx).Debug().Msg("Writing bootstrap file")

	bootw, err := zipWriter.Create("bootstrap")
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}

	if err = bootstrapTemplate.Execute(bootw, filepath.Base(binaryName)); err != nil {
		return fmt.Errorf("failed to write bootstrap file: %w", err)
	}

	zerolog.Ctx(ctx).Debug().Msg("Closing archive")

	if err := zipWriter.Close(); err != nil {
		return fmt.Errorf("failed to close archive: %w", err)
	}

	return nil
}
