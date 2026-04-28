package cmd

import (
	"fmt"
	"runtime"
	"termtyper/database"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func createStyles(termProfile termenv.Profile, foregroundColor termenv.Color, themeColor string) Styles {
	return Styles{
		correct: func(str string) termenv.Style {
			return termenv.String(str).Foreground(foregroundColor)
		},
		toEnter: func(str string) termenv.Style {
			return termenv.String(str).Foreground(foregroundColor).Faint()
		},
		mistake: func(str string) termenv.Style {
			return termenv.String(str).Foreground(termProfile.Color("1"))
		},
		cursor: func(str string) termenv.Style {
			return termenv.String(str).Reverse().Bold()
		},
		themeFunc: func(str string) termenv.Style {
			return termenv.String(str).Foreground(termProfile.Color(themeColor))
		},
	}
}

func getTermProfile(m *model) termenv.Profile {
	return m.termProfile
}

func getForegroundColor(m *model) termenv.Color {
	return m.foregroundColor
}

func (m model) Init() tea.Cmd {
	return nil
}

func initModel(termProfile termenv.Profile, foregroundColor termenv.Color, width, height int, sess *Session) model {
	databaseContext := database.InitDB()

	themeName := "Magenta"
	if sess.User != nil && sess.User.Config != nil && sess.User.Config.Theme != "" {
		themeName = sess.User.Config.Theme
	}
	themeColor := GetThemeColor(themeName)

	m := model{
		width:           width,
		height:          height,
		context:         databaseContext,
		session:         sess,
		termProfile:     termProfile,
		foregroundColor: foregroundColor,
		styles:          createStyles(termProfile, foregroundColor, themeColor),
	}

	// Initialize state machine with pre-authentication state
	m.stateMachine = NewStateMachine(&m)
	m.stateMachine.SetCurrentState(StatePreAuth)
	m.stateMachine.handlers[StatePreAuth] = NewPreAuthHandler(&databaseContext)

	return m
}

func OsInit() {
	if runtime.GOOS == "windows" {
		mode, err := termenv.EnableWindowsANSIConsole()
		if err != nil {
			fmt.Println(err)
		}
		defer termenv.RestoreWindowsConsole(mode)
	}
}
