package cmd

import (
	"github.com/Bahaaio/pomo/config"
	"github.com/spf13/cobra"
)

var workCmd = &cobra.Command{
	Use:   "work",
	Short: "start a pomodoro work session",
	Run: func(cmd *cobra.Command, args []string) {
		runTask(config.C.Work)
	},
}

func init() {
	rootCmd.AddCommand(workCmd)
}
