package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	LLMProvider  string   `yaml:"llm_provider"`
	APIKey       string   `yaml:"api_key"`
	TrustedDirs  []string `yaml:"trusted_dirs"`
	Theme        string   `yaml:"theme"`
	DefaultModel string   `yaml:"default_model"`
	AutoUpdate   bool     `yaml:"auto_update"`
}

var DefaultConfig = Config{
	LLMProvider:  "anthropic",
	APIKey:       "",
	TrustedDirs:  []string{},
	Theme:        "default",
	DefaultModel: "claude-3-sonnet-20240229",
	AutoUpdate:   true,
}

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".aura.conf"), nil
}

func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return CreateDefaultConfig()
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func CreateDefaultConfig() (*Config, error) {
	config := DefaultConfig
	err := config.Save()
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (c *Config) Save() error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0600)
}

func (c *Config) IsTrustedDir(dir string) bool {
	for _, trustedDir := range c.TrustedDirs {
		if trustedDir == dir {
			return true
		}
	}
	return false
}

func (c *Config) AddTrustedDir(dir string) error {
	if !c.IsTrustedDir(dir) {
		c.TrustedDirs = append(c.TrustedDirs, dir)
		return c.Save()
	}
	return nil
}

func (c *Config) IsConfigured() bool {
	return c.APIKey != ""
}

func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key is required. Please set it in your config file at %s", getConfigPathSafe())
	}

	validProviders := []string{"anthropic", "openai"}
	isValid := false
	for _, provider := range validProviders {
		if c.LLMProvider == provider {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid LLM provider: %s. Valid providers are: %v", c.LLMProvider, validProviders)
	}

	return nil
}

func getConfigPathSafe() string {
	path, err := GetConfigPath()
	if err != nil {
		return "~/.aura.conf"
	}
	return path
}
