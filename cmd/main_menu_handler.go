package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MainMenuHandler handles the main menu state
type MainMenuHandler struct {
	*BaseStateHandler
	menu MainMenu
}

// NewMainMenuHandler creates a new main menu handler
func NewMainMenuHandler(menu MainMenu) *MainMenuHandler {
	return &MainMenuHandler{
		BaseStateHandler: NewBaseStateHandler(StateMainMenu),
		menu:             menu,
	}
}

// HandleInput implements StateHandler
func (h *MainMenuHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			switch h.menu.MainMenuSelection[h.menu.cursor] {
			case "Timer":
				if h.ValidateTransition(StateTimerTest, context) {
					return NewTimerTestHandler(initTimerTest(h.menu)), nil
				}
			case "Zen":
				if h.ValidateTransition(StateZenMode, context) {
					return NewZenModeHandler(initZenMode(h.menu)), nil
				}
			case "Word Count":
				if h.ValidateTransition(StateWordCountTest, context) {
					return NewWordCountTestHandler(initWordCountTest(h.menu)), nil
				}
			case "Config":
				if h.ValidateTransition(StateSettings, context) {
					return NewSettingsHandler(initSettings(context.model.session.User)), nil
				}
			}
		case "up", "k":
			if h.menu.cursor == 0 {
				h.menu.cursor = len(h.menu.MainMenuSelection) - 1
			} else {
				h.menu.cursor--
			}
		case "down", "j":
			if h.menu.cursor == len(h.menu.MainMenuSelection)-1 {
				h.menu.cursor = 0
			} else {
				h.menu.cursor++
			}
		}
	}
	return h, nil
}

// Render implements StateHandler
func (h *MainMenuHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	termtyper := style("TermTyper - Welcome "+m.session.User.Username, m.styles.magenta)
	termtyper = lipgloss.NewStyle().PaddingBottom(1).Render(termtyper)

	var menuItems []string
	menuItemsStyle := lipgloss.NewStyle().PaddingTop(1)

	for i, choice := range h.menu.MainMenuSelection {
		choiceShow := style(choice, m.styles.toEnter)
		choiceShow = wrapWithCursor(h.menu.cursor == i, choiceShow, m.styles.toEnter)
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

// ValidateTransition implements StateHandler
func (h *MainMenuHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateMainMenu]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}
