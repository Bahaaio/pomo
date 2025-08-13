package cmd

import (
	"log"

	"github.com/Bahaaio/pomo/config"
	"github.com/spf13/cobra"
)

var breakCmd = &cobra.Command{
	Use:   "break",
	Short: "start a pomodoro break session",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("breakCmd args:", args)
		runTask(&config.C.Break, cmd)
	},
}

func init() {
	rootCmd.AddCommand(breakCmd)
}
