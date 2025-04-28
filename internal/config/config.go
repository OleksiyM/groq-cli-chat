package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"groq-cli-chat/internal/groq"
	"groq-cli-chat/resources"
)

type Config struct {
	APIKey      string   `mapstructure:"api_key"`
	BaseURL     string   `mapstructure:"base_url"`
	Models      []string `mapstructure:"models"`
	DefaultModel string   `mapstructure:"default_model"`
}

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
		} else {
			return nil, fmt.Errorf(resources.ErrReadConfig, err)
		}
	} else {
		if err := viper.Unmarshal(cfg); err != nil {
			return nil, fmt.Errorf(resources.ErrUnmarshalConfig, err)
		}
	}

	// Load API key from environment
	cfg.APIKey = os.Getenv("GROQ_API_KEY")
	if cfg.APIKey == "" {
		return nil, fmt.Errorf(resources.ErrNoAPIKey)
	}

	// Validate models
	if err := ValidateModels(cfg.Models); err != nil {
		return nil, fmt.Errorf(resources.ErrInvalidConfig, err)
	}

	// Validate default model
	if !IsValidModel(cfg.DefaultModel, cfg.Models) {
		return nil, fmt.Errorf(resources.ErrInvalidDefaultModel, cfg.DefaultModel)
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

	models, err := client.ListModels()
	if err != nil {
		return nil, fmt.Errorf(resources.ErrListModels, err)
	}

	// Validate models (ensure non-empty)
	if err := ValidateModels(models); err != nil {
		return nil, fmt.Errorf(resources.ErrListModels, err)
	}

	cfg := &Config{
		BaseURL:      resources.DefaultBaseURL,
		Models:       models,
		DefaultModel: models[0],
	}

	viper.Set("base_url", cfg.BaseURL)
	viper.Set("models", cfg.Models)
	viper.Set("default_model", cfg.DefaultModel)

	configPath := filepath.Join(configDir, "config.yaml")
	if err := viper.WriteConfigAs(configPath); err != nil {
		return nil, fmt.Errorf(resources.ErrWriteConfig, err)
	}

	return cfg, nil
}
