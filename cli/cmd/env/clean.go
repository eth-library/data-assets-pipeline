package env

import (
	"os"
	"path/filepath"

	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

// CleanCmd removes .venv and caches.
var CleanCmd = &cobra.Command{
	Use:     "clean",
	Short:   "Remove .venv and caches",
	Long:    "Remove the virtual environment and all Python cache directories.",
	GroupID: GroupID,
	RunE: func(cmd *cobra.Command, args []string) error {
		ui.Info("Cleaning environment...")

		// Remove .venv
		if err := os.RemoveAll(".venv"); err != nil {
			ui.Warn("Could not remove .venv", "error", err)
		} else {
			ui.Success("Removed .venv")
		}

		// Remove cache directories
		cachePatterns := []string{
			"__pycache__",
			".pytest_cache",
			".ruff_cache",
			".mypy_cache",
			"*.egg-info",
		}

		for _, pattern := range cachePatterns {
			matches, _ := filepath.Glob("**/" + pattern)
			// Also check root level
			rootMatches, _ := filepath.Glob(pattern)
			matches = append(matches, rootMatches...)

			for _, match := range matches {
				if err := os.RemoveAll(match); err != nil {
					ui.Warn("Could not remove "+match, "error", err)
				}
			}
		}

		// Walk directories to find nested caches
		filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() {
				return nil
			}
			name := info.Name()
			for _, pattern := range cachePatterns {
				if matched, _ := filepath.Match(pattern, name); matched {
					os.RemoveAll(path)
					return filepath.SkipDir
				}
			}
			return nil
		})

		ui.Success("Environment cleaned")
		return nil
	},
}
