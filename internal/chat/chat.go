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
	fmt.Println(resources.WelcomeMessage)

	client, err := groq.NewClient(cfg.BaseURL, cfg.APIKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, resources.ErrCreateClient, err)
		os.Exit(1)
	}

	currentModel := cfg.DefaultModel
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
			} else {
				currentModel = newModel
			}
		case "h":
			if err := ListChatHistory(); err != nil {
				fmt.Fprintf(os.Stderr, resources.ErrReadHistoryDir, err)
			}
		case "q":
			fmt.Println(resources.GoodbyeMessage)
			return
		default:
			if input != "" {
				resp, err := client.Chat(currentModel, input)
				if err != nil {
					fmt.Fprintf(os.Stderr, resources.ErrChat, err)
					continue
				}

				fmt.Println(resp.Choices[0].Message.Content)
				fmt.Printf(resources.StatsFormat,
					resp.Usage.TotalTokens,
					resp.Usage.CompletionTime,
					float64(resp.Usage.TotalTokens)/resp.Usage.CompletionTime)

				if err := saveChatHistory(input, resp, currentModel); err != nil {
					fmt.Fprintf(os.Stderr, resources.ErrSaveHistory, err)
				}
			}
		}
	}
}

func selectModel(models []string, currentModel string) (string, error) {
	// Limit to 10 models for selection
	displayModels := models
	if len(displayModels) > 10 {
		displayModels = displayModels[:10]
	}

	fmt.Println(resources.SelectModelHeader)
	for i, model := range displayModels {
		fmt.Printf("%d - %s\n", i, model)
	}
	fmt.Printf(resources.SelectModelPrompt, len(displayModels)-1)

	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		return "", fmt.Errorf(resources.ErrReadInput)
	}
	choice := strings.TrimSpace(scanner.Text())

	index, err := parseChoice(choice, len(displayModels))
	if err != nil {
		return "", err // Return error to indicate invalid selection
	}
	return displayModels[index], nil
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