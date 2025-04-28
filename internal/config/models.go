package config

import (
	"fmt"
	"groq-cli-chat/resources"
)

// ValidateModels ensures the model list is valid (non-empty)
func ValidateModels(models []string) error {
	if len(models) == 0 {
		return fmt.Errorf(resources.ErrNoModels)
	}
	return nil
}

// IsValidModel checks if a model is in the list of available models
func IsValidModel(model string, models []string) bool {
	for _, m := range models {
		if m == model {
			return true
		}
	}
	return false
}