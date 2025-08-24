package cmd

import (
	"termtyper/database"
	"termtyper/words"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MainMenuHandler struct {
	*BaseStateHandler
	MainMenuSelection      []string
	cursor                 int
	timerTestWordGenerator words.WordGenerator
	wordTestWordGenerator  words.WordGenerator
	currentUser            *database.ApplicationUser
}

func NewMainMenuHandler(user *database.ApplicationUser) *MainMenuHandler {
	return &MainMenuHandler{
		BaseStateHandler: NewBaseStateHandler(StateMainMenu),
		MainMenuSelection: []string{
			"Timer",
			"Word Count",
			"Zen",
			"Config",
		},
		currentUser:            user,
		cursor:                 0,
		timerTestWordGenerator: words.NewGenerator(),
		wordTestWordGenerator:  words.NewGenerator(),
	}
}

func (h *MainMenuHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch h.MainMenuSelection[h.cursor] {
			case "Timer":
				if h.ValidateTransition(StateTimerTest, context) {
					return NewTimerTestHandler(*h), nil
				}
			case "Zen":
				if h.ValidateTransition(StateZenMode, context) {
					return NewZenModeHandler(*h), nil
				}
			case "Word Count":
				if h.ValidateTransition(StateWordCountTest, context) {
					return NewWordCountTestHandler(*h), nil
				}
			case "Config":
				if h.ValidateTransition(StateSettings, context) {
					return NewSettingsHandler(context.model.session.User), nil
				}
			}
		case "up", "k":
			if h.cursor == 0 {
				h.cursor = len(h.MainMenuSelection) - 1
			} else {
				h.cursor--
			}
		case "down", "j":
			if h.cursor == len(h.MainMenuSelection)-1 {
				h.cursor = 0
			} else {
				h.cursor++
			}
		}
	}
	return h, nil
}

func (h *MainMenuHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	termtyper := style("TermTyper - Welcome "+m.session.User.Username, m.styles.magenta)
	termtyper = lipgloss.NewStyle().PaddingBottom(1).Render(termtyper)

	var menuItems []string
	menuItemsStyle := lipgloss.NewStyle().PaddingTop(1)

	for i, choice := range h.MainMenuSelection {
		choiceShow := style(choice, m.styles.toEnter)
		choiceShow = wrapWithCursor(h.cursor == i, choiceShow, m.styles.toEnter)
		choiceShow = menuItemsStyle.Render(choiceShow)
		menuItems = append(menuItems, choiceShow)
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{termtyper}, menuItems...)...)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FF00FF"))

	return borderStyle.Render(centeredText)
}

func (h *MainMenuHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateMainMenu]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}
