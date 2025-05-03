package chat

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"groq-cli-chat/resources"
)

// ListChatHistory retrieves and displays a list of saved chat history files
func ListChatHistory() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf(resources.ErrHomeDir, err)
	}

	historyDir := filepath.Join(homeDir, ".groq-chat", "history")
	entries, err := os.ReadDir(historyDir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No chat history found.")
			return nil
		}
		return fmt.Errorf(resources.ErrReadHistoryDir, err)
	}

	var historyFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
			historyFiles = append(historyFiles, entry.Name())
		}
	}

	if len(historyFiles) == 0 {
		fmt.Println("No chat history found.")
		return nil
	}

	sort.Strings(historyFiles)
	fmt.Println("────────┤ Chat History ├─────────")
	for i, file := range historyFiles {
		timestamp := strings.TrimPrefix(strings.TrimSuffix(file, ".md"), "chat_")
		parsedTime, _ := time.Parse("20060102_150405", timestamp)
		fmt.Printf("%d - %s (%s)\n", i, file, parsedTime.Format(time.RFC1123))
	}

	return nil
}