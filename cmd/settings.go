package cmd

import (
	"termtyper/database"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Settings struct {
	wordCountSelection []int
	wordCountCursor    int
	timerSelection     []int
	timerCursor        int
	settingsSelection  []string
	settingsCursor     int
}

func initSettings(user *database.ApplicationUser) Settings {
	wordCountSelections := []int{15, 30, 45, 60}
	timerSelection := []int{15, 30, 60, 120}
	settingsSelection := []string{"Words", "Time"}

	return Settings{
		wordCountSelection: wordCountSelections,
		wordCountCursor:    findWordsIndex(user.Config, wordCountSelections),
		timerSelection:     timerSelection,
		timerCursor:        findTimerIndex(user.Config, timerSelection),
		settingsSelection:  settingsSelection,
		settingsCursor:     0,
	}
}

func (s Settings) renderSettings(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	settings := style("Settings", m.styles.magenta)
	settings = lipgloss.NewStyle().PaddingBottom(1).Render(settings)

	menuItemsStyle := lipgloss.NewStyle().PaddingTop(1)
	var settingSelection []string
	for i, choice := range s.settingsSelection {
		choiceShow := style(choice, m.styles.toEnter)

		choiceShow = wrapWithCursor(s.settingsCursor == i, choiceShow, m.styles.toEnter)
		choiceShow = menuItemsStyle.Render(choiceShow)
		settingSelection = append(settingSelection, choiceShow)
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{settings}, settingSelection...)...)
	renderString := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, renderString)

	return centeredText
}

func (state *Settings) updateSettings(msg tea.Msg) State {

	newCursor := state.settingsCursor
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":

		case "up", "k":
			if state.settingsCursor == 0 {
				newCursor = len(state.settingsSelection) - 1
			} else {
				newCursor--
			}

		case "down", "j":
			if state.settingsCursor == len(state.settingsSelection)-1 {
				newCursor = 0
			} else {
				newCursor++
			}
		}
	}
	state.settingsCursor = newCursor
	return *state
}

func findWordsIndex(config *database.UserConfig, wordCountSelection []int) int {
	for i, num := range wordCountSelection {
		if num == config.Words {
			return i
		}
	}
	return 0
}

func findTimerIndex(config *database.UserConfig, timerSelection []int) int {
	for i, num := range timerSelection {
		if num == config.Time {
			return i
		}
	}
	return 0
}
