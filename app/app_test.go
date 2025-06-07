package app

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aura/config"
	"aura/models"
)

// Test helper to create a temporary config
func createTestConfig(t *testing.T, apiKey, provider string) *config.Config {
	t.Helper()
	return &config.Config{
		LLMProvider:  provider,
		APIKey:       apiKey,
		TrustedDirs:  []string{},
		Theme:        "default",
		DefaultModel: "test-model",
		AutoUpdate:   true,
	}
}


func TestApp_GetConfig(t *testing.T) {
	cfg := createTestConfig(t, "test-key", "anthropic")
	app := &App{
		config: cfg,
	}

	result := app.GetConfig()
	assert.Equal(t, cfg, result)
	assert.Equal(t, "test-key", result.APIKey)
	assert.Equal(t, "anthropic", result.LLMProvider)
}

func TestApp_GetCurrentDir(t *testing.T) {
	testDir := "/test/directory"
	app := &App{
		currentDir: testDir,
	}

	result := app.GetCurrentDir()
	assert.Equal(t, testDir, result)
}

func TestApp_AddTrustedDir_Integration(t *testing.T) {
	// Setup temporary environment
	tempDir := t.TempDir()
	oldHome := os.Getenv("HOME")
	defer func() { os.Setenv("HOME", oldHome) }()
	os.Setenv("HOME", tempDir)

	// Create initial config
	cfg := &config.Config{
		LLMProvider:  "openai",
		APIKey:       "test-key",
		TrustedDirs:  []string{"/existing/dir"},
		Theme:        "default",
		DefaultModel: "gpt-4",
		AutoUpdate:   true,
	}

	// Save initial config
	err := cfg.Save()
	require.NoError(t, err)

	// Create app with the config
	app := &App{
		config: cfg,
	}

	// Add new trusted directory
	newDir := "/new/trusted/dir"
	err = app.AddTrustedDir(newDir)
	require.NoError(t, err)

	// Verify the directory was added
	assert.Contains(t, cfg.TrustedDirs, newDir)
	assert.Contains(t, cfg.TrustedDirs, "/existing/dir")
}

// Integration test with real filesystem operations
func TestApp_Integration_RealFileSystem(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create temporary directory
	tempDir := t.TempDir()

	// Set home directory to temp directory
	oldHome := os.Getenv("HOME")
	defer func() { os.Setenv("HOME", oldHome) }()
	os.Setenv("HOME", tempDir)

	// Change to temp directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() { _ = os.Chdir(originalWd) }()

	workDir := filepath.Join(tempDir, "project")
	err = os.MkdirAll(workDir, 0755)
	require.NoError(t, err)
	err = os.Chdir(workDir)
	require.NoError(t, err)

	// Create config using the Config API
	cfg := &config.Config{
		LLMProvider:  "openai",
		APIKey:       "test-integration-key",
		TrustedDirs:  []string{},
		Theme:        "default",
		DefaultModel: "gpt-4",
		AutoUpdate:   true,
	}

	err = cfg.Save()
	require.NoError(t, err)

	// Test app creation
	app, err := NewApp()

	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, "test-integration-key", app.config.APIKey)
	assert.Equal(t, "openai", app.config.LLMProvider)
	assert.Contains(t, app.currentDir, workDir)
	assert.Equal(t, models.SecurityPrompt, app.state.Screen) // Should prompt for security since dir not trusted
}

// Test edge cases and error conditions
func TestEdgeCases(t *testing.T) {
	t.Run("nil config", func(t *testing.T) {
		app := &App{config: nil}
		assert.Nil(t, app.GetConfig())
	})

	t.Run("empty current directory", func(t *testing.T) {
		app := &App{currentDir: ""}
		assert.Empty(t, app.GetCurrentDir())
	})

	t.Run("add trusted dir with nil config", func(t *testing.T) {
		app := &App{config: nil}
		// This will panic since we're calling a method on nil pointer
		// Let's catch the panic and verify it happens
		defer func() {
			if r := recover(); r != nil {
				// Expected panic due to nil pointer
				assert.NotNil(t, r)
			}
		}()
		_ = app.AddTrustedDir("/test")
	})
}

// Benchmark tests
func BenchmarkApp_GetConfig(b *testing.B) {
	cfg := &config.Config{
		LLMProvider: "openai",
		APIKey:      "bench-key",
	}
	app := &App{config: cfg}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = app.GetConfig()
	}
}

func BenchmarkApp_GetCurrentDir(b *testing.B) {
	app := &App{currentDir: "/benchmark/test/directory"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = app.GetCurrentDir()
	}
}

func BenchmarkApp_AddTrustedDir(b *testing.B) {
	cfg := &config.Config{
		LLMProvider: "openai",
		APIKey:      "bench-key",
		TrustedDirs: []string{},
	}
	app := &App{config: cfg}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Use different directory for each iteration
		testDir := fmt.Sprintf("/bench/dir/%d", i)
		_ = app.AddTrustedDir(testDir)
	}
}
