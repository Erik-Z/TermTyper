package cmd

import (
	"termtyper/database"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsHandler struct {
	*BaseStateHandler
	settingsCursor    int
	settingSelections []TestSetting
}

func NewSettingsHandler(user *database.ApplicationUser) *SettingsHandler {
	wordCountSelection := []int{15, 30, 45, 60}
	timerSelection := []int{15, 30, 60, 120}

	timerSettings := TimerSettings{
		timerSelection:  timerSelection,
		timerCursor:     findTimerIndex(user.Config, timerSelection),
		selectionCursor: 0,
	}

	wordsSettings := WordsSettings{
		wordsSelection:  wordCountSelection,
		wordsCursor:     findWordsIndex(user.Config, wordCountSelection),
		selectionCursor: 0,
	}

	return &SettingsHandler{
		BaseStateHandler:  NewBaseStateHandler(StateSettings),
		settingsCursor:    0,
		settingSelections: []TestSetting{timerSettings, wordsSettings},
	}
}

func (h *SettingsHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	newCursor := h.settingsCursor
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if h.ValidateTransition(StateMainMenu, context) {
				return NewMainMenuHandler(context.model.session.User), nil
			}

		case "enter":

		case "up", "k":
			if h.settingsCursor == 0 {
				newCursor = len(h.settingSelections) - 1
			} else {
				newCursor--
			}

		case "down", "j":
			if h.settingsCursor == len(h.settingSelections)-1 {
				newCursor = 0
			} else {
				newCursor++
			}
		}
	}
	h.settingsCursor = newCursor
	return h, nil
}

func (h *SettingsHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	settings := style("Settings", m.styles.magenta)
	settings = lipgloss.NewStyle().PaddingBottom(1).Render(settings)

	menuItemsStyle := lipgloss.NewStyle().PaddingTop(1)
	var settingSelection []string
	for i, choice := range h.settingSelections {
		choiceShow := choice.render(m.styles)

		choiceShow = wrapWithCursor(h.settingsCursor == i, choiceShow, m.styles.toEnter)
		choiceShow = menuItemsStyle.Render(choiceShow)
		settingSelection = append(settingSelection, choiceShow)
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{settings}, settingSelection...)...)
	renderString := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, renderString)

	return centeredText
}

func (h *SettingsHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateSettings]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}

func findWordsIndex(config *database.UserConfig, wordCountSelection []int) int {
	for i, num := range wordCountSelection {
		if num == config.Words {
			return i
		}
	}
	return 0
}
