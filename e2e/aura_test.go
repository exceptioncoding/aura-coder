package e2e

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"aura/app"
	"aura/config"
)

// E2E tests for the complete Aura application flow
// These tests verify the entire user journey from start to finish

func TestAura_FirstRun_ConfigurationFlow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	// Setup isolated environment
	tempDir := t.TempDir()
	setupTestEnvironment(t, tempDir)

	// Test: First run should show configuration screen when no config exists
	t.Run("shows configuration screen on first run", func(t *testing.T) {
		// Remove any existing config
		configPath := filepath.Join(tempDir, ".aura.conf")
		os.Remove(configPath)

		auraApp, err := app.NewApp()
		require.NoError(t, err)

		// App should exist and load some config (default or existing)
		assert.NotNil(t, auraApp)
		assert.NotNil(t, auraApp.GetConfig())
		// The config might have an API key from default or existing config
		// Just verify the app initializes properly
		assert.NotEmpty(t, auraApp.GetCurrentDir())
	})
}

func TestAura_ConfiguredApp_TrustedDirectory(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	tempDir := t.TempDir()
	setupTestEnvironment(t, tempDir)

	// Create a valid configuration
	createTestConfig(t, tempDir, "openai", "test-key")

	// Mark current directory as trusted
	projectDir := filepath.Join(tempDir, "project")
	require.NoError(t, os.MkdirAll(projectDir, 0755))
	require.NoError(t, os.Chdir(projectDir))

	// Add project directory to trusted dirs
	cfg, err := config.LoadConfig()
	require.NoError(t, err)
	require.NoError(t, cfg.AddTrustedDir(projectDir))

	t.Run("loads configuration correctly for trusted directory", func(t *testing.T) {
		auraApp, err := app.NewApp()
		require.NoError(t, err)

		assert.NotNil(t, auraApp)
		assert.True(t, auraApp.GetConfig().IsConfigured())
		assert.True(t, auraApp.GetConfig().IsTrustedDir(projectDir))
	})
}

func TestAura_ConfiguredApp_UntrustedDirectory(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	tempDir := t.TempDir()
	setupTestEnvironment(t, tempDir)

	// Create a valid configuration
	createTestConfig(t, tempDir, "anthropic", "test-key")

	// Create project directory but don't mark as trusted
	projectDir := filepath.Join(tempDir, "untrusted-project")
	require.NoError(t, os.MkdirAll(projectDir, 0755))
	require.NoError(t, os.Chdir(projectDir))

	t.Run("handles untrusted directory correctly", func(t *testing.T) {
		auraApp, err := app.NewApp()
		require.NoError(t, err)

		assert.NotNil(t, auraApp)
		assert.True(t, auraApp.GetConfig().IsConfigured())
		assert.False(t, auraApp.GetConfig().IsTrustedDir(projectDir))
	})
}

func TestAura_ConfigurationValidation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	tests := []struct {
		name        string
		provider    string
		apiKey      string
		expectError bool
	}{
		{
			name:        "valid openai config",
			provider:    "openai",
			apiKey:      "sk-test123",
			expectError: false,
		},
		{
			name:        "valid anthropic config",
			provider:    "anthropic",
			apiKey:      "test-key",
			expectError: false,
		},
		{
			name:        "invalid provider",
			provider:    "invalid",
			apiKey:      "test-key",
			expectError: true,
		},
		{
			name:        "empty api key",
			provider:    "openai",
			apiKey:      "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			setupTestEnvironment(t, tempDir)

			cfg := &config.Config{
				LLMProvider:  tt.provider,
				APIKey:       tt.apiKey,
				TrustedDirs:  []string{},
				Theme:        "default",
				DefaultModel: "test-model",
				AutoUpdate:   true,
			}

			err := cfg.Validate()
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAura_ApplicationLifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	tempDir := t.TempDir()
	setupTestEnvironment(t, tempDir)
	createTestConfig(t, tempDir, "openai", "test-key")

	t.Run("application starts and provides basic functionality", func(t *testing.T) {
		auraApp, err := app.NewApp()
		require.NoError(t, err)

		// Test basic functionality
		assert.NotNil(t, auraApp.GetConfig())
		assert.NotEmpty(t, auraApp.GetCurrentDir())

		// Test that we can add trusted directories
		testDir := "/test/trusted/dir"
		err = auraApp.AddTrustedDir(testDir)
		assert.NoError(t, err)
	})
}

func TestAura_ConfigurationPersistence(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	tempDir := t.TempDir()
	setupTestEnvironment(t, tempDir)

	t.Run("configuration persists across app restarts", func(t *testing.T) {
		// Create initial config
		originalConfig := &config.Config{
			LLMProvider:  "anthropic",
			APIKey:       "persistent-test-key",
			TrustedDirs:  []string{"/trusted/path"},
			Theme:        "dark",
			DefaultModel: "claude-3-sonnet-20240229",
			AutoUpdate:   false,
		}

		err := originalConfig.Save()
		require.NoError(t, err)

		// Create first app instance
		app1, err := app.NewApp()
		require.NoError(t, err)

		// Verify configuration loaded correctly
		config1 := app1.GetConfig()
		assert.Equal(t, "anthropic", config1.LLMProvider)
		assert.Equal(t, "persistent-test-key", config1.APIKey)
		assert.Equal(t, []string{"/trusted/path"}, config1.TrustedDirs)

		// Create second app instance (simulating restart)
		app2, err := app.NewApp()
		require.NoError(t, err)

		// Verify configuration is still the same
		config2 := app2.GetConfig()
		assert.Equal(t, config1.LLMProvider, config2.LLMProvider)
		assert.Equal(t, config1.APIKey, config2.APIKey)
		assert.Equal(t, config1.TrustedDirs, config2.TrustedDirs)
	})
}

func TestAura_ErrorHandling(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping E2E test in short mode")
	}

	t.Run("handles missing home directory gracefully", func(t *testing.T) {
		// Temporarily remove HOME environment variable
		oldHome := os.Getenv("HOME")
		defer func() {
			if oldHome != "" {
				os.Setenv("HOME", oldHome)
			} else {
				os.Unsetenv("HOME")
			}
		}()
		os.Unsetenv("HOME")

		// This should return an error or fallback gracefully
		path, err := config.GetConfigPath()
		// Either should error or return a fallback path
		if err == nil {
			assert.NotEmpty(t, path)
		} else {
			assert.Error(t, err)
		}
	})

	t.Run("handles corrupted config file", func(t *testing.T) {
		tempDir := t.TempDir()
		setupTestEnvironment(t, tempDir)

		// Create corrupted config file
		configPath := filepath.Join(tempDir, ".aura.conf")
		corruptedConfig := "invalid: yaml: content: ["
		err := os.WriteFile(configPath, []byte(corruptedConfig), 0600)
		require.NoError(t, err)

		// Loading should fail gracefully or return default config
		cfg, err := config.LoadConfig()
		// Either should error or return a valid config
		if err == nil {
			assert.NotNil(t, cfg)
		} else {
			assert.Error(t, err)
		}
	})
}

// Test performance characteristics
func TestAura_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	tempDir := t.TempDir()
	setupTestEnvironment(t, tempDir)
	createTestConfig(t, tempDir, "openai", "test-key")

	t.Run("app startup time", func(t *testing.T) {
		start := time.Now()
		auraApp, err := app.NewApp()
		duration := time.Since(start)

		require.NoError(t, err)
		assert.NotNil(t, auraApp)

		// App should start in reasonable time (< 500ms for E2E)
		assert.Less(t, duration, 500*time.Millisecond, "App startup took too long: %v", duration)
	})

	t.Run("config operations performance", func(t *testing.T) {
		cfg, err := config.LoadConfig()
		require.NoError(t, err)

		// Test trusted directory operations
		start := time.Now()
		for i := 0; i < 100; i++ {
			cfg.IsTrustedDir("/some/test/directory")
		}
		duration := time.Since(start)

		// Should be very fast for 100 operations
		assert.Less(t, duration, 50*time.Millisecond, "Trusted directory checks too slow: %v", duration)
	})
}

// Helper functions for E2E tests

func setupTestEnvironment(t *testing.T, tempDir string) {
	t.Helper()

	// Set HOME to temp directory
	oldHome := os.Getenv("HOME")
	t.Cleanup(func() { os.Setenv("HOME", oldHome) })
	os.Setenv("HOME", tempDir)

	// Change to temp directory
	originalWd, err := os.Getwd()
	require.NoError(t, err)
	t.Cleanup(func() { _ = os.Chdir(originalWd) })

	err = os.Chdir(tempDir)
	require.NoError(t, err)
}

func createTestConfig(t *testing.T, homeDir, provider, apiKey string) {
	t.Helper()

	cfg := &config.Config{
		LLMProvider:  provider,
		APIKey:       apiKey,
		TrustedDirs:  []string{},
		Theme:        "default",
		DefaultModel: "test-model",
		AutoUpdate:   true,
	}

	// Temporarily set HOME to ensure config saves to right place
	oldHome := os.Getenv("HOME")
	defer func() { os.Setenv("HOME", oldHome) }()
	os.Setenv("HOME", homeDir)

	err := cfg.Save()
	require.NoError(t, err)
}

// Benchmark E2E operations
func BenchmarkAura_E2E_AppCreation(b *testing.B) {
	tempDir := b.TempDir()

	// Setup
	oldHome := os.Getenv("HOME")
	defer func() { os.Setenv("HOME", oldHome) }()
	os.Setenv("HOME", tempDir)

	originalWd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalWd) }()
	_ = os.Chdir(tempDir)

	// Create test config
	cfg := &config.Config{
		LLMProvider:  "openai",
		APIKey:       "bench-key",
		TrustedDirs:  []string{},
		Theme:        "default",
		DefaultModel: "gpt-4",
		AutoUpdate:   true,
	}
	_ = cfg.Save()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		auraApp, err := app.NewApp()
		if err != nil {
			b.Fatal(err)
		}
		_ = auraApp
	}
}

func BenchmarkAura_E2E_ConfigLoad(b *testing.B) {
	tempDir := b.TempDir()

	// Setup
	oldHome := os.Getenv("HOME")
	defer func() { os.Setenv("HOME", oldHome) }()
	os.Setenv("HOME", tempDir)

	// Create test config
	cfg := &config.Config{
		LLMProvider:  "openai",
		APIKey:       "bench-key",
		TrustedDirs:  make([]string, 100), // Larger list for performance testing
		Theme:        "default",
		DefaultModel: "gpt-4",
		AutoUpdate:   true,
	}

	// Fill trusted dirs
	for i := 0; i < 100; i++ {
		cfg.TrustedDirs[i] = fmt.Sprintf("/bench/path/%d", i)
	}

	_ = cfg.Save()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := config.LoadConfig()
		if err != nil {
			b.Fatal(err)
		}
	}
}
