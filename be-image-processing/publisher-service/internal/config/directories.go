package config

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
)

func InitializeDirectories() error {
	dirs := []string{
		"./uploads",
		"./compressed",
	}

	for _, dir := range dirs {
		absPath, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("failed to get absolute path for %s: %w", dir, err)
		}

		err = os.MkdirAll(absPath, 0755)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", absPath, err)
		}
		slog.Info(fmt.Sprintf("Ensured directory exists: %s", absPath))
	}

	return nil
}
