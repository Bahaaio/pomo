// Package cmd provides the command-line interface for the pomo timer.
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Bahaaio/pomo/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pomo",
	Short: "pomo is a simple cli pomodoro",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("args:", args)
		runTask(config.C.Work, "work")
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	initializeLogging()
	initConfig()
}

func initConfig() {
	log.Println("initializing config")

	if err := config.LoadConfig(); err != nil {
		fmt.Fprintln(os.Stderr, "could not load config", err)
		os.Exit(1)
	}
}

func initializeLogging() {
	if len(os.Getenv("DEBUG")) > 0 {
		_, err := tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to setup logging:", err)
			os.Exit(1)
		}
		log.SetFlags(log.Ltime)
		log.Println("debug mode")
	} else {
		log.SetOutput(io.Discard)
	}
}
