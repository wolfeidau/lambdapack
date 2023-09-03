package binaries

import (
	"fmt"
	"os"
	"strings"
)

func Filter(files ...string) ([]string, error) {

	filtered := make([]string, 0)

	for _, file := range files {
		fin, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("failed to stat file %s: %w", file, err)
		}

		if fin.IsDir() {
			// skip directories
			continue
		}

		if strings.HasSuffix(fin.Name(), ".zip") {
			// skip existing zip files
			continue
		}

		filtered = append(filtered, file)
	}

	return filtered, nil
}
