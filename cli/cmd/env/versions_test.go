package env

import (
	"testing"
)

func TestVersionParsing(t *testing.T) {
	// Test that version functions don't panic when tools are missing
	// These return empty string on error, which is the expected behavior

	t.Run("getPythonVersion handles missing python gracefully", func(t *testing.T) {
		// This test just ensures no panic - actual version depends on environment
		_ = getPythonVersion()
	})

	t.Run("getUVVersion handles missing uv gracefully", func(t *testing.T) {
		_ = getUVVersion()
	})

	t.Run("getDagsterVersion handles missing dagster gracefully", func(t *testing.T) {
		_ = getDagsterVersion()
	})

	t.Run("getKubectlVersion handles missing kubectl gracefully", func(t *testing.T) {
		_ = getKubectlVersion()
	})

	t.Run("getHelmVersion handles missing helm gracefully", func(t *testing.T) {
		_ = getHelmVersion()
	})
}

func TestShowVersionsFunctions(t *testing.T) {
	// Ensure the show functions don't panic
	t.Run("ShowVersionsCompact doesn't panic", func(t *testing.T) {
		ShowVersionsCompact()
	})

	t.Run("ShowVersionsFull doesn't panic", func(t *testing.T) {
		ShowVersionsFull()
	})
}
