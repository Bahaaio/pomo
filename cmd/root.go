// Package cmd provides the command-line interface for the pomo timer.
package cmd

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Bahaaio/pomo/config"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/gen2brain/beeep"
	"github.com/spf13/cobra"
)

var version = "0.3.1"

var rootCmd = &cobra.Command{
	Use:     "pomo [duration]",
	Short:   "start a pomodoro work session (default: 25m)",
	Version: version,
	Long: `pomo is a simple terminal-based Pomodoro timer

Start a work session with the default duration from your config file,
or specify a custom duration. The timer shows a progress bar and sends
desktop notifications when complete.`,
	Example: `  pomo           # Start work session (default: 25m)
  pomo 1h15m     # Start 1 hour 15 minute session`,

	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("rootCmd args:", args)
		runTask(&config.C.Work, cmd)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	initLogging()
	initConfig()
	beeep.AppName = "pomo"
}

func initConfig() {
	log.Println("initializing config")

	if err := config.LoadConfig(); err != nil {
		log.Println("using default config")
	}
}

func initLogging() {
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
