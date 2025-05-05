package chat

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"groq-cli-chat/internal/config"
	"groq-cli-chat/internal/groq"
	"groq-cli-chat/resources"
)

func Run(cfg *config.Config) {
	// Only use the welcome message from resources
	fmt.Println(resources.WelcomeMessage)

	client, err := groq.NewClient(cfg.BaseURL, cfg.APIKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, resources.ErrCreateClient, err)
		os.Exit(1)
	}

	currentModel := cfg.DefaultModel
	
	// Check if default model is empty and prompt user to select one
	if currentModel == "" {
		fmt.Println("No default model found in config. Please select a model to use:")
		newModel, err := selectModel(cfg.Models, "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to select model: %v\n", err)
			os.Exit(1)
		}
		currentModel = newModel
	}
	
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf(resources.Prompt, currentModel)
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		switch input {
		case "i":
			modelInfo, err := client.GetModel(currentModel)
			if err != nil {
				fmt.Fprintf(os.Stderr, resources.ErrGetModel, err)
				fmt.Println() // Add a blank line
			} else {
				fmt.Printf(resources.InfoModelDetails,
					currentModel,
					modelInfo.OwnedBy,
					modelInfo.Active,
					modelInfo.ContextWindow)
			}
		
		case "m":
			newModel, err := selectModel(cfg.Models, currentModel)
			if err != nil {
				fmt.Printf(resources.InfoModelUnchanged, currentModel)
				fmt.Println() // Add a blank line
			} else {
				currentModel = newModel
			}
		
		case "h":
			if err := ListChatHistory(); err != nil {
				fmt.Fprintf(os.Stderr, resources.ErrReadHistoryDir, err)
				fmt.Println() // Add a blank line
			}
		
		case "u":
			if err := updateModels(cfg, client); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to update models: %v\n", err)
				fmt.Println() // Add a blank line
			}
		
		case "c":
			if err := changeConfig(cfg, client); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to change configuration: %v\n", err)
				fmt.Println() // Add a blank line
			} else {
				// Update current model to the new default model
				currentModel = cfg.DefaultModel
			}
		case "q":
			fmt.Println(resources.GoodbyeMessage)
			return
		// In the default case of your switch statement
		default:
		if input != "" {
		// Start timing the request
		startTime := time.Now()
		
		resp, err := client.Chat(currentModel, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, resources.ErrChat, err)
			fmt.Println() // Add a blank line after error message
			continue
		}
		
		// Calculate elapsed time if needed
		elapsedTime := time.Since(startTime).Seconds()
		if resp.Usage.CompletionTime <= 0 {
			resp.Usage.CompletionTime = elapsedTime
		}
		
		// Display the response content
		fmt.Println(resp.Choices[0].Message.Content)
		
		// Calculate tokens per second (avoid division by zero)
		tokensPerSecond := 0.0
		if resp.Usage.CompletionTime > 0 {
			tokensPerSecond = float64(resp.Usage.TotalTokens) / resp.Usage.CompletionTime
		}
		
		// Display statistics
		fmt.Printf(resources.StatsFormat,
			resp.Usage.TotalTokens,
			resp.Usage.CompletionTime,
			tokensPerSecond)
		fmt.Println() // Add a blank line after stats
		
		if err := saveChatHistory(input, resp, currentModel); err != nil {
			fmt.Fprintf(os.Stderr, resources.ErrSaveHistory, err)
			fmt.Println() // Add a blank line
		}
		}
		}
	}
}

func selectModel(models []string, currentModel string) (string, error) {
	// Limit to 20 models for selection
	displayModels := models
	if len(displayModels) > 20 {
		displayModels = displayModels[:20]
	}

	fmt.Println(resources.SelectModelHeader)
	for i, model := range displayModels {
		fmt.Printf("%d - %s\n", i, model)
	}
	
	// Keep prompting until valid selection or explicit cancel
	for {
		fmt.Printf(resources.SelectModelPrompt, len(displayModels)-1)

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return "", fmt.Errorf(resources.ErrReadInput)
		}
		choice := strings.TrimSpace(scanner.Text())
		
		// Allow user to cancel selection
		if choice == "" || strings.ToLower(choice) == "q" || strings.ToLower(choice) == "quit" {
			return "", fmt.Errorf("model selection cancelled")
		}
		
		index, err := parseChoice(choice, len(displayModels))
		if err != nil {
			// Show error but allow retry
			fmt.Printf("Invalid selection: %v. Please try again or press Enter/Q to cancel.\n", err)
			continue
		}
		
		return displayModels[index], nil
	}
}

func parseChoice(choice string, max int) (int, error) {
	var index int
	if _, err := fmt.Sscanf(choice, "%d", &index); err != nil {
		return -1, fmt.Errorf(resources.ErrInvalidChoice, choice)
	}
	if index < 0 || index >= max {
		return -1, fmt.Errorf(resources.ErrChoiceOutOfRange, index, max)
	}
	return index, nil
}

func saveChatHistory(input string, resp *groq.ChatResponse, model string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf(resources.ErrHomeDir, err)
	}

	historyDir := filepath.Join(homeDir, ".groq-chat", "history")
	if err := os.MkdirAll(historyDir, 0755); err != nil {
		return fmt.Errorf(resources.ErrCreateHistoryDir, err)
	}

	timestamp := time.Now().Format("20060102_150405")
	filename := filepath.Join(historyDir, fmt.Sprintf("chat_%s.md", timestamp))
	content := fmt.Sprintf(resources.HistoryFormat,
		timestamp, model, input, resp.Choices[0].Message.Content,
		resp.Usage.TotalTokens, resp.Usage.CompletionTime,
		float64(resp.Usage.TotalTokens)/resp.Usage.CompletionTime)

	return os.WriteFile(filename, []byte(content), 0644)
}

func updateModels(cfg *config.Config, client *groq.Client) error {
	fmt.Printf("Fetching latest models from %s API...\n", cfg.ProviderName)
	
	// Use the existing client instead of creating a new one with a modified URL
	newModels, err := client.ListModels()
	if err != nil {
		return fmt.Errorf("failed to fetch models: %v", err)
	}
	
	// Filter out excluded models
	var filteredNewModels []string
	for _, model := range newModels {
		shouldInclude := true
		for _, excluded := range cfg.ExcludedModels {
			if strings.Contains(strings.ToLower(model), strings.ToLower(excluded)) {
				shouldInclude = false
				break
			}
		}
		if shouldInclude {
			filteredNewModels = append(filteredNewModels, model)
		}
	}
	
	// Compare old and new lists
	oldModels := cfg.Models
	
	// Check if lists are identical
	if areModelListsIdentical(oldModels, filteredNewModels) {
		fmt.Println("No updates available. Your model list is already up to date.")
		return nil
	}
	
	// Find new models (in new list but not in old list)
	var addedModels []string
	for _, newModel := range filteredNewModels {
		if !contains(oldModels, newModel) {
			addedModels = append(addedModels, newModel)
		}
	}
	
	// Find removed models (in old list but not in new list)
	var removedModels []string
	for _, oldModel := range oldModels {
		if !contains(filteredNewModels, oldModel) {
			removedModels = append(removedModels, oldModel)
		}
	}
	
	// Display changes
	fmt.Println("────────┤ Model Updates Available ├─────────")
	
	if len(addedModels) > 0 {
		fmt.Println("New models:")
		for _, model := range addedModels {
			fmt.Printf("  + %s\n", model)
		}
	}
	
	if len(removedModels) > 0 {
		fmt.Println("Removed models:")
		for _, model := range removedModels {
			fmt.Printf("  - %s\n", model)
		}
	}
	
	fmt.Printf("Old list: %d models | New list: %d models\n", len(oldModels), len(filteredNewModels))
	fmt.Println("─────────────────────────────────────")
	
	// Ask user if they want to update
	fmt.Print("Do you want to update the models list? (y/n): ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return fmt.Errorf("failed to read input")
	}
	
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if response == "y" || response == "yes" {
		// Update config
		// No need to construct a new path, use the one from the config
		configPath := cfg.ConfigPath
		
		// Update the config struct
		cfg.Models = filteredNewModels
		
		// Check if default model is still valid
		if !contains(filteredNewModels, cfg.DefaultModel) && len(filteredNewModels) > 0 {
			fmt.Printf("Warning: Your default model '%s' is no longer available. Setting default to '%s'.\n", 
				cfg.DefaultModel, filteredNewModels[0])
			cfg.DefaultModel = filteredNewModels[0]
		}
		
		// Save the updated config
		if err := config.SaveConfig(cfg, configPath); err != nil {
			return fmt.Errorf("failed to save updated config: %v", err)
		}
		
		fmt.Println("Models list updated successfully!")
	} else {
		fmt.Println("Update cancelled. Models list remains unchanged.")
	}
	
	return nil
}

// Helper function to check if two model lists are identical
func areModelListsIdentical(list1, list2 []string) bool {
	if len(list1) != len(list2) {
		return false
	}
	
	// Create maps for faster lookup
	map1 := make(map[string]bool)
	for _, item := range list1 {
		map1[item] = true
	}
	
	// Check if all items in list2 are in list1
	for _, item := range list2 {
		if !map1[item] {
			return false
		}
	}
	
	return true
}

// Helper function to check if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// changeConfig allows the user to select a different configuration file
func changeConfig(cfg *config.Config, client *groq.Client) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf(resources.ErrHomeDir, err)
	}

	configDir := filepath.Join(homeDir, ".groq-chat")
	
	// List all YAML files in the config directory
	files, err := os.ReadDir(configDir)
	if err != nil {
		return fmt.Errorf("failed to read config directory: %v", err)
	}
	
	var yamlFiles []string
	for _, file := range files {
		if !file.IsDir() && (strings.HasSuffix(file.Name(), ".yaml") || strings.HasSuffix(file.Name(), ".yml")) {
			yamlFiles = append(yamlFiles, file.Name())
		}
	}
	
	if len(yamlFiles) == 0 {
		return fmt.Errorf("no configuration files found in %s", configDir)
	}
	
	// Display available configuration files
	fmt.Println("────────┤ Available Configurations ├─────────")
	for i, file := range yamlFiles {
		fmt.Printf("%d - %s\n", i, file)
	}
	
	// Prompt user to select a configuration
	var selectedConfig string
	for {
		fmt.Printf("─────────────────────────────────────\nSelect configuration (0-%d): ", len(yamlFiles)-1)
		
		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			return fmt.Errorf(resources.ErrReadInput)
		}
		choice := strings.TrimSpace(scanner.Text())
		
		// Allow user to cancel selection
		if choice == "" || strings.ToLower(choice) == "q" || strings.ToLower(choice) == "quit" {
			return fmt.Errorf("configuration selection cancelled")
		}
		
		index, err := parseChoice(choice, len(yamlFiles))
		if err != nil {
			fmt.Printf("Invalid selection: %v. Please try again or press Enter/Q to cancel.\n", err)
			continue
		}
		
		selectedConfig = yamlFiles[index]
		break
	}
	
	// Load and validate the selected configuration
	newCfg, err := config.LoadSpecificConfig(filepath.Join(configDir, selectedConfig))
	if err != nil {
		return fmt.Errorf("failed to load selected configuration: %v", err)
	}
	
	// Validate the new configuration
	if err := validateConfig(newCfg); err != nil {
		return fmt.Errorf("invalid configuration: %v", err)
	}
	
	// Ask for confirmation
	fmt.Printf("New configuration loaded from %s\n", selectedConfig)
	fmt.Printf("Base URL: %s\n", newCfg.BaseURL)
	fmt.Printf("Default Model: %s\n", newCfg.DefaultModel)
	fmt.Printf("Available Models: %d\n", len(newCfg.Models))
	
	fmt.Print("Do you want to apply this configuration? (y/n): ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return fmt.Errorf(resources.ErrReadInput)
	}
	
	response := strings.ToLower(strings.TrimSpace(scanner.Text()))
	if response != "y" && response != "yes" {
		return fmt.Errorf("configuration change cancelled")
	}
	
	// Apply the new configuration
	*cfg = *newCfg
	
	// Create a new client with the updated configuration
	newClient, err := groq.NewClient(cfg.BaseURL, cfg.APIKey)
	if err != nil {
		return fmt.Errorf("failed to create client with new configuration: %v", err)
	}
	
	// Update the client reference
	*client = *newClient
	
	// Display success message with config file name
	fmt.Printf("Configuration updated successfully from '%s'!\n", selectedConfig)
	
	// Display the app title from the new configuration
	fmt.Println("\n" + cfg.AppTitle)
	fmt.Println(resources.MenuOptions)
	fmt.Println() // Add a blank line after menu options
	
	return nil
}

// validateConfig checks if the configuration has all required fields
func validateConfig(cfg *config.Config) error {
	// Check for required fields
	if cfg.BaseURL == "" {
		return fmt.Errorf("base_url is missing")
	}
	
	if len(cfg.Models) == 0 {
		return fmt.Errorf("no models defined")
	}
	
	// If default model is specified, check if it exists in the models list
	if cfg.DefaultModel != "" {
		found := false
		for _, model := range cfg.Models {
			if model == cfg.DefaultModel {
				found = true
				break
			}
		}
		
		if !found {
			return fmt.Errorf("default model '%s' not found in models list", cfg.DefaultModel)
		}
	} else {
		// If no default model is specified, set it to the first model
		cfg.DefaultModel = cfg.Models[0]
	}
	
	// Check if API key is set
	if cfg.APIKey == "" {
		// Try to get it from environment
		cfg.APIKey = os.Getenv(cfg.APIKeyName)
		if cfg.APIKey == "" {
			return fmt.Errorf("API key not found in environment variable %s", cfg.APIKeyName)
		}
	}
	
	return nil
}
