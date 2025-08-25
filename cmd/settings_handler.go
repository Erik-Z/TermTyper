package cmd

import (
	"fmt"
	"termtyper/database"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SettingsHandler struct {
	*BaseStateHandler
	settingsCursor    int
	settingSelections []TestSetting
}

type TestSetting interface {
	render(styles Styles) string
}

type TimerSettings struct {
	timerCursor     int
	timerSelection  []int
	selectionCursor int
}

type WordsSettings struct {
	wordsCursor     int
	wordsSelection  []int
	selectionCursor int
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

func (t TimerSettings) render(style Styles) string {
	selections := []string{formatSettingsDuration(t.timerSelection[t.timerCursor])}
	selectionsStr := showSelections(selections, t.selectionCursor, style)
	return fmt.Sprintf("%s %s", "Timer", selectionsStr)
}

func (w WordsSettings) render(style Styles) string {
	selections := []string{fmt.Sprintf("%d", w.wordsSelection[w.wordsCursor])}
	selectionsStr := showSelections(selections, w.selectionCursor, style)
	return fmt.Sprintf("%s %s", "Words", selectionsStr)
}

func findTimerIndex(config *database.UserConfig, timerSelection []int) int {
	for i, num := range timerSelection {
		if num == config.Time {
			return i
		}
	}
	return 0
}

func showSelections(selections []string, cursor int, styles Styles) string {
	var selectionsStr string
	for i, option := range selections {
		if i+1 == cursor {
			selectionsStr += "[" + style(option, styles.magenta) + "]"
		} else {
			selectionsStr += "[" + style(option, styles.toEnter) + "]"
		}
		selectionsStr += " "
	}
	return selectionsStr
}

func findWordsIndex(config *database.UserConfig, wordCountSelection []int) int {
	for i, num := range wordCountSelection {
		if num == config.Words {
			return i
		}
	}
	return 0
}

func formatSettingsDuration(seconds int) string {
	minutes := seconds / 60
	remainingSeconds := seconds % 60
	return fmt.Sprintf("%dm%ds", minutes, remainingSeconds)
}
