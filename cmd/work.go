package cmd

import (
	"log"

	"github.com/Bahaaio/pomo/config"
	"github.com/spf13/cobra"
)

var workCmd = &cobra.Command{
	Use:   "work",
	Short: "start a pomodoro work session",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("workCmd args:", args)
		runTask(&config.C.Work, cmd)
	},
}

func init() {
	rootCmd.AddCommand(workCmd)
}
