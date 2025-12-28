package cmd

import (
	"github.com/Bahaaio/pomo/ui/stats"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Args:  cobra.MaximumNArgs(0),
	Short: "Display Pomodoro statistics and productivity metrics",
	Run: func(cmd *cobra.Command, args []string) {
		m := stats.New()
		p := tea.NewProgram(m)

		_, err := p.Run()
		if err != nil {
			die(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
}
