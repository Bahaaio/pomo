package cmd

import (
	"github.com/Bahaaio/pomo/config"
	"github.com/spf13/cobra"
)

var breakCmd = &cobra.Command{
	Use:   "break",
	Short: "start a pomodoro break session",
	Run: func(cmd *cobra.Command, args []string) {
		runTask(config.C.Break)
	},
}

func init() {
	rootCmd.AddCommand(breakCmd)
}
