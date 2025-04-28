package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"groq-cli-chat/internal/chat"
	"groq-cli-chat/internal/config"
	"groq-cli-chat/resources"
)

func main() {
	// Initialize configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, resources.ErrLoadConfig, err)
		os.Exit(1)
	}

	// Initialize root command
	rootCmd := &cobra.Command{
		Use:   "groq-cli-chat",
		Short: "A CLI tool to chat with Groq AI models",
		Run: func(cmd *cobra.Command, args []string) {
			chat.Run(cfg)
		},
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, resources.ErrExecuteCmd, err)
		os.Exit(1)
	}
}