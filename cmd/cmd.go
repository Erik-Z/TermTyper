package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	Version = "dev"
	RootCmd = &cobra.Command{
		Use:  "TermTyper",
		Long: "TermTyper - Terminal Typing Test",
		RunE: func(cmd *cobra.Command, args []string) error {
			termWidth, termHeight, err := term.GetSize(int(os.Stdout.Fd()))

			if err != nil {
				fmt.Println("Error getting terminal size", err)
				return err
			}

			p := tea.NewProgram(
				initModel(
					termenv.ColorProfile(),
					termenv.ForegroundColor(),
					termWidth, termHeight),
				tea.WithAltScreen(),
			)

			_, err = p.Run()
			return err
		},
	}
)
