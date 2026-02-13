package cmd

import (
	"fmt"
	"os"

	"github.com/Tucker-Sullivan/cockatrice-printlock/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Initialize interactive TUI",
	Long:  "Starts the TUI for a more interactive experience",
	RunE: func(cmd *cobra.Command, args []string) error {
		model, err := tui.InitModel()
		if err != nil {
			return err
		}
		p := tea.NewProgram(model)
		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}
