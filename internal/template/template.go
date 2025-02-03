package template

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type templateConfig struct {
	fs            embed.FS // Embedded filesystem
	basePath      string   // Template path in FS (e.g., "templates/empty")
	oldModuleName string   // Module placeholder (e.g., "example.com/app")
}

func CreateFromTemplate(config templateConfig) error {
	// Read go.mod file from current directory
	modFile, err := os.ReadFile("go.mod")
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	// Get first line and split by spaces to get module name
	lines := strings.Split(string(modFile), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty go.mod file")
	}

	parts := strings.Fields(lines[0])
	if len(parts) < 2 {
		return fmt.Errorf("invalid go.mod file: missing module name")
	}

	moduleName := parts[1]

	// Walk through template files
	err = fs.WalkDir(config.fs, config.basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == config.basePath {
			return nil
		}

		if strings.HasSuffix(path, ".keep") {
			return nil
		}

		if strings.HasSuffix(path, ".DS_Store") {
			return nil
		}

		// Calculate relative path
		relPath := strings.TrimPrefix(path, config.basePath+"/")

		if d.IsDir() {
			return os.MkdirAll(relPath, 0755)
		}

		// Read file content
		content, err := config.fs.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file %s: %w", path, err)
		}

		content = []byte(strings.ReplaceAll(string(content), config.oldModuleName, moduleName))

		// code here:
		//

		oldRootPkg := filepath.Base(config.oldModuleName)
		newRootPkg := filepath.Base(moduleName)

		content = []byte(strings.ReplaceAll(string(content),
			fmt.Sprintf("package %s", oldRootPkg),
			fmt.Sprintf("package %s", newRootPkg)))
		content = []byte(strings.ReplaceAll(string(content),
			fmt.Sprintf("%s.", oldRootPkg),
			fmt.Sprintf("%s.", newRootPkg)))

		err = os.WriteFile(relPath, content, 0644)
		if err != nil {
			return fmt.Errorf("failed to write file %s: %w", relPath, err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to create template: %w", err)
	}

	return nil
}
