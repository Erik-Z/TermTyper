package cmd

import (
	"termtyper/database"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

func createStyles(themeColor string) Styles {
	return Styles{
		correct: func(str string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("2")).Render(str)
		},
		toEnter: func(str string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("7")).Faint(true).Render(str)
		},
		mistake: func(str string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Render(str)
		},
		cursor: func(str string) string {
			return lipgloss.NewStyle().Reverse(true).Bold(true).Render(str)
		},
		themeFunc: func(str string) string {
			return lipgloss.NewStyle().Foreground(lipgloss.Color(themeColor)).Render(str)
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func initModel(width, height int, sess *Session) model {
	databaseContext := database.InitDB()

	themeName := "Magenta"
	if sess.User != nil && sess.User.Config != nil && sess.User.Config.Theme != "" {
		themeName = sess.User.Config.Theme
	}
	themeColor := GetThemeColor(themeName)

	m := model{
		width:   width,
		height:  height,
		context: databaseContext,
		session: sess,
		styles:  createStyles(themeColor),
	}

	// Initialize state machine with pre-authentication state
	m.stateMachine = NewStateMachine(&m)
	m.stateMachine.SetCurrentState(StatePreAuth)
	m.stateMachine.handlers[StatePreAuth] = NewPreAuthHandler(&databaseContext)

	return m
}
