package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	expected := Config{
		LLMProvider:  "anthropic",
		APIKey:       "",
		TrustedDirs:  []string{},
		Theme:        "default",
		DefaultModel: "claude-3-sonnet-20240229",
		AutoUpdate:   true,
	}

	assert.Equal(t, expected.LLMProvider, DefaultConfig.LLMProvider)
	assert.Equal(t, expected.APIKey, DefaultConfig.APIKey)
	assert.Equal(t, expected.TrustedDirs, DefaultConfig.TrustedDirs)
	assert.Equal(t, expected.Theme, DefaultConfig.Theme)
	assert.Equal(t, expected.DefaultModel, DefaultConfig.DefaultModel)
	assert.Equal(t, expected.AutoUpdate, DefaultConfig.AutoUpdate)
}

func TestGetConfigPath(t *testing.T) {
	// Test with current environment
	path, err := GetConfigPath()
	assert.NoError(t, err)
	assert.Contains(t, path, ".aura.conf")
}

func TestConfig_IsTrustedDir(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		testDir     string
		expected    bool
	}{
		{
			name: "directory is trusted",
			config: Config{
				TrustedDirs: []string{"/home/user/project", "/opt/app"},
			},
			testDir:  "/home/user/project",
			expected: true,
		},
		{
			name: "directory is not trusted",
			config: Config{
				TrustedDirs: []string{"/home/user/project", "/opt/app"},
			},
			testDir:  "/tmp/untrusted",
			expected: false,
		},
		{
			name: "empty trusted dirs",
			config: Config{
				TrustedDirs: []string{},
			},
			testDir:  "/any/directory",
			expected: false,
		},
		{
			name: "exact match required",
			config: Config{
				TrustedDirs: []string{"/home/user/project"},
			},
			testDir:  "/home/user/project/subdir",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsTrustedDir(tt.testDir)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_IsConfigured(t *testing.T) {
	tests := []struct {
		name     string
		config   Config
		expected bool
	}{
		{
			name: "configured with API key",
			config: Config{
				APIKey: "test-api-key",
			},
			expected: true,
		},
		{
			name: "not configured - empty API key",
			config: Config{
				APIKey: "",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.IsConfigured()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid openai config",
			config: Config{
				LLMProvider: "openai",
				APIKey:      "sk-test123",
			},
			expectError: false,
		},
		{
			name: "valid anthropic config",
			config: Config{
				LLMProvider: "anthropic",
				APIKey:      "test-key",
			},
			expectError: false,
		},
		{
			name: "missing API key",
			config: Config{
				LLMProvider: "openai",
				APIKey:      "",
			},
			expectError: true,
			errorMsg:    "API key is required",
		},
		{
			name: "invalid provider",
			config: Config{
				LLMProvider: "invalid-provider",
				APIKey:      "test-key",
			},
			expectError: true,
			errorMsg:    "invalid LLM provider",
		},
		{
			name: "empty provider",
			config: Config{
				LLMProvider: "",
				APIKey:      "test-key",
			},
			expectError: true,
			errorMsg:    "invalid LLM provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestConfig_SaveAndLoad_Integration(t *testing.T) {
	// Use a temporary directory for this test
	tempDir := t.TempDir()
	
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer func() { os.Setenv("HOME", originalHome) }()
	
	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)

	// Create test config
	testConfig := &Config{
		LLMProvider:  "openai",
		APIKey:       "test-integration-key",
		TrustedDirs:  []string{"/test/path1", "/test/path2"},
		Theme:        "dark",
		DefaultModel: "gpt-4",
		AutoUpdate:   false,
	}

	// Save config
	err := testConfig.Save()
	require.NoError(t, err)

	// Verify file exists
	configPath, err := GetConfigPath()
	require.NoError(t, err)
	assert.FileExists(t, configPath)

	// Load config
	loadedConfig, err := LoadConfig()
	require.NoError(t, err)

	// Verify loaded config matches
	assert.Equal(t, testConfig.LLMProvider, loadedConfig.LLMProvider)
	assert.Equal(t, testConfig.APIKey, loadedConfig.APIKey)
	assert.Equal(t, testConfig.TrustedDirs, loadedConfig.TrustedDirs)
	assert.Equal(t, testConfig.Theme, loadedConfig.Theme)
	assert.Equal(t, testConfig.DefaultModel, loadedConfig.DefaultModel)
	assert.Equal(t, testConfig.AutoUpdate, loadedConfig.AutoUpdate)
}

func TestConfig_AddTrustedDir_Integration(t *testing.T) {
	// Use a temporary directory for this test
	tempDir := t.TempDir()
	
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer func() { os.Setenv("HOME", originalHome) }()
	
	// Set HOME to temp directory
	os.Setenv("HOME", tempDir)

	// Create initial config
	config := &Config{
		LLMProvider:  "openai",
		APIKey:       "test-key",
		TrustedDirs:  []string{"/existing/dir"},
		Theme:        "default",
		DefaultModel: "gpt-4",
		AutoUpdate:   true,
	}

	// Save initial config
	err := config.Save()
	require.NoError(t, err)

	// Add new trusted directory
	newDir := "/new/trusted/dir"
	err = config.AddTrustedDir(newDir)
	require.NoError(t, err)

	// Verify the directory was added
	assert.Contains(t, config.TrustedDirs, newDir)
	assert.Contains(t, config.TrustedDirs, "/existing/dir")

	// Load config from disk to verify persistence
	loadedConfig, err := LoadConfig()
	require.NoError(t, err)
	assert.Contains(t, loadedConfig.TrustedDirs, newDir)
	assert.Contains(t, loadedConfig.TrustedDirs, "/existing/dir")

	// Test adding duplicate directory
	err = config.AddTrustedDir(newDir)
	require.NoError(t, err)

	// Should still only have one copy
	count := 0
	for _, dir := range config.TrustedDirs {
		if dir == newDir {
			count++
		}
	}
	assert.Equal(t, 1, count, "Should not have duplicate trusted directories")
}

func TestGetConfigPathSafe(t *testing.T) {
	result := getConfigPathSafe()
	assert.NotEmpty(t, result)
	assert.Contains(t, result, ".aura.conf")
}

// Benchmark tests
func BenchmarkConfig_IsTrustedDir(b *testing.B) {
	config := Config{
		TrustedDirs: make([]string, 100),
	}

	// Fill with test directories
	for i := 0; i < 100; i++ {
		config.TrustedDirs[i] = filepath.Join("/path/to/dir", string(rune(i)))
	}

	testDir := "/path/to/dir/50" // Middle element

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.IsTrustedDir(testDir)
	}
}

func BenchmarkConfig_Validate(b *testing.B) {
	config := Config{
		LLMProvider: "openai",
		APIKey:      "test-api-key",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		config.Validate()
	}
}
