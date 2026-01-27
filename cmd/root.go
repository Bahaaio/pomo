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

var version = "1.0.5"

var rootCmd = &cobra.Command{
	Use:     "pomo [work duration] [break duration]",
	Short:   "start a pomodoro work session (default: 25m)",
	Version: version,
	Long: `pomo is a simple terminal-based Pomodoro timer

Start a work session with the default duration from your config file,
or specify a custom duration. The timer shows a progress bar and sends
desktop notifications when complete.`,
	Example: `  pomo           # Start work session (default: 25m)
  pomo 1h15m     # Start 1 hour 15 minute session
  pomo 45m 15m   # Start 45 minute work session with 15 minute break`,

	Args: cobra.MaximumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("rootCmd args:", args)
		runTask(config.WorkTask, cmd)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	initLogging()
	initConfig()
	beeep.AppName = config.AppName
}

func initConfig() {
	log.Println("initializing config")

	config.Setup()
	if err := config.LoadConfig(); err != nil {
		die(fmt.Errorf("could not load config: %w", err))
	}
}

func initLogging() {
	debugEnv := os.Getenv("DEBUG")
	if debugEnv == "" || debugEnv == "0" {
		log.SetOutput(io.Discard)
		return
	}

	_, err := tea.LogToFile("debug.log", "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to setup logging:", err)
		os.Exit(1)
	}

	log.SetFlags(log.Ltime)
}

func die(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
	os.Exit(1)
}
