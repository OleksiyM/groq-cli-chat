package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/viper"
	"groq-cli-chat/internal/groq"
	"groq-cli-chat/resources"
)

type Config struct {
	AppTitle      string   `mapstructure:"app_title"`
	ProviderName  string   `mapstructure:"provider_name"`
	BaseURL       string   `mapstructure:"base_url"`
	APIKeyName    string   `mapstructure:"api_key_name"`
	DefaultModel  string   `mapstructure:"default_model"`
	Models        []string `mapstructure:"models"`
	ExcludedModels []string `mapstructure:"excluded_models"`
	APIKey        string   `mapstructure:"api_key"`
	ConfigPath    string   // Path to the loaded config file (not stored in YAML)
}

// Default excluded models - will be moved to config
var defaultExcludedModels = []string{"whisper", "playai"}

func LoadConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf(resources.ErrHomeDir, err)
	}

	configDir := filepath.Join(homeDir, ".groq-chat")
	if err := ensureConfigDir(configDir); err != nil {
		return nil, fmt.Errorf(resources.ErrConfigDir, err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	cfg := &Config{}
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			cfg, err = createDefaultConfig(configDir)
			if err != nil {
				return nil, fmt.Errorf(resources.ErrCreateConfig, err)
			}
			fmt.Println(resources.InfoConfigCreated)
			cfg.ConfigPath = filepath.Join(configDir, "config.yaml")
		} else {
			return nil, fmt.Errorf(resources.ErrReadConfig, err)
		}
	} else {
		if err := viper.Unmarshal(cfg); err != nil {
			return nil, fmt.Errorf(resources.ErrUnmarshalConfig, err)
		}
		cfg.ConfigPath = viper.ConfigFileUsed()
	}

	// Set default API key name if not specified
	if cfg.APIKeyName == "" {
		cfg.APIKeyName = "GROQ_API_KEY"
	}

	// Load API key from environment using the configured name
	cfg.APIKey = os.Getenv(cfg.APIKeyName)
	if cfg.APIKey == "" {
		return nil, fmt.Errorf(resources.ErrNoAPIKey+": %s", cfg.APIKeyName)
	}

	// Validate models
	if err := ValidateModels(cfg.Models); err != nil {
		return nil, fmt.Errorf(resources.ErrInvalidConfig, err)
	}

	// Validate default model
	if cfg.DefaultModel != "" && !IsValidModel(cfg.DefaultModel, cfg.Models) {
		return nil, fmt.Errorf(resources.ErrInvalidDefaultModel, cfg.DefaultModel)
	}

	// Set default excluded models if not specified
	if len(cfg.ExcludedModels) == 0 {
		cfg.ExcludedModels = defaultExcludedModels
	}

	return cfg, nil
}

func ensureConfigDir(configDir string) error {
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		return os.MkdirAll(configDir, 0755)
	}
	return nil
}

func createDefaultConfig(configDir string) (*Config, error) {
	client, err := groq.NewClient(resources.DefaultBaseURL, os.Getenv("GROQ_API_KEY"))
	if err != nil {
		return nil, fmt.Errorf(resources.ErrCreateClient, err)
	}

	allModels, err := client.ListModels()
	if err != nil {
		return nil, fmt.Errorf(resources.ErrListModels, err)
	}
	
	// Use the default excluded models for initial config
	excludedModelsList := defaultExcludedModels
	
	// Filter out excluded models
	var filteredModels []string
	for _, model := range allModels {
		shouldInclude := true
		for _, excluded := range excludedModelsList {
			if strings.Contains(strings.ToLower(model), strings.ToLower(excluded)) {
				shouldInclude = false
				break
			}
		}
		if shouldInclude {
			filteredModels = append(filteredModels, model)
		}
	}
	
	// Sort models alphabetically
	sort.Strings(filteredModels)

	// Validate models (ensure non-empty)
	if err := ValidateModels(filteredModels); err != nil {
		return nil, fmt.Errorf(resources.ErrListModels, err)
	}

	configPath := filepath.Join(configDir, "config.yaml")
	
	cfg := &Config{
		AppTitle:      "🍎 One-shot Groq CLI chat",
		ProviderName:  "Groq",
		BaseURL:       resources.DefaultBaseURL,
		APIKeyName:    "GROQ_API_KEY",
		DefaultModel:  filteredModels[0],
		Models:        filteredModels,
		ExcludedModels: excludedModelsList,
		ConfigPath:    configPath,
	}

	viper.Set("app_title", cfg.AppTitle)
	viper.Set("provider_name", cfg.ProviderName)
	viper.Set("base_url", cfg.BaseURL)
	viper.Set("api_key_name", cfg.APIKeyName)
	viper.Set("excluded_models", cfg.ExcludedModels)
	viper.Set("default_model", cfg.DefaultModel)
	viper.Set("models", cfg.Models)

	if err := viper.WriteConfigAs(configPath); err != nil {
		return nil, fmt.Errorf(resources.ErrWriteConfig, err)
	}

	return cfg, nil
}

// SaveConfig saves the configuration to the specified path
func SaveConfig(cfg *Config, configPath string) error {
	// Set the values in viper
	viper.Set("app_title", cfg.AppTitle)
	viper.Set("provider_name", cfg.ProviderName)
	viper.Set("base_url", cfg.BaseURL)
	viper.Set("api_key_name", cfg.APIKeyName)
	viper.Set("excluded_models", cfg.ExcludedModels)
	viper.Set("default_model", cfg.DefaultModel)
	viper.Set("models", cfg.Models)

	// Write the config to file
	if err := viper.WriteConfigAs(configPath); err != nil {
		return fmt.Errorf(resources.ErrWriteConfig, err)
	}

	return nil
}

// LoadSpecificConfig loads a configuration from a specific file path
func LoadSpecificConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)
	
	cfg := &Config{}
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	
	if err := v.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf(resources.ErrUnmarshalConfig, err)
	}
	
	// Set the config path
	cfg.ConfigPath = configPath
	
	// Set default API key name if not specified
	if cfg.APIKeyName == "" {
		cfg.APIKeyName = "GROQ_API_KEY"
	}
	
	// Load API key from environment using the configured name
	cfg.APIKey = os.Getenv(cfg.APIKeyName)
	if cfg.APIKey == "" {
		return nil, fmt.Errorf(resources.ErrNoAPIKey+": %s", cfg.APIKeyName)
	}
	
	// Validate models
	if err := ValidateModels(cfg.Models); err != nil {
		return nil, fmt.Errorf(resources.ErrInvalidConfig, err)
	}
	
	// Validate default model
	if cfg.DefaultModel != "" && !IsValidModel(cfg.DefaultModel, cfg.Models) {
		return nil, fmt.Errorf(resources.ErrInvalidDefaultModel, cfg.DefaultModel)
	}
	
	// Set default excluded models if not specified
	if len(cfg.ExcludedModels) == 0 {
		cfg.ExcludedModels = defaultExcludedModels
	}
	
	return cfg, nil
}
