/*
Copyright © 2025 2xhamzeh
*/
package template

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type templateConfig struct {
	fs            embed.FS // Embedded filesystem
	basePath      string   // Template path in FS (e.g., "templates/empty")
	oldModuleName string   // Module placeholder (e.g., "example.com/app")
}

func CreateFromTemplate(config templateConfig, moduleName string) error {

	// create go.mod
	cmd := exec.Command("go", "mod", "init", moduleName)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create go.mod, %s", strings.TrimSpace(stderr.String()))
	}

	// Walk through template files
	err := fs.WalkDir(config.fs, config.basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// skip root directory
		if path == config.basePath {
			return nil
		}

		// skip .keep files
		if strings.HasSuffix(path, ".keep") {
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

	cmd = exec.Command("go", "mod", "tidy")
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create go.mod, %s", strings.TrimSpace(stderr.String()))
	}

	return nil

}
