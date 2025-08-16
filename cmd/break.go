package cmd

import (
	"log"

	"github.com/Bahaaio/pomo/config"
	"github.com/spf13/cobra"
)

var breakCmd = &cobra.Command{
	Use:   "break [duration]",
	Short: "start a pomodoro break session",
	Example: `  pomo break       # Start a break session
  pomo break 15m   # Start 15 minute break session`,

	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("breakCmd args:", args)
		runTask(&config.C.Break, cmd)
	},
}

func init() {
	rootCmd.AddCommand(breakCmd)
}
