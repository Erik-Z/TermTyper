package cmd

import (
	"termtyper/database"
	"termtyper/words"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type MainMenuHandler struct {
	*BaseStateHandler
	MainMenuSelection      []string
	cursor                 int
	timerTestWordGenerator words.WordGenerator
	wordTestWordGenerator  words.WordGenerator
	currentUser            *database.ApplicationUser
}

func NewMainMenuHandler(user *database.ApplicationUser, m *model) *MainMenuHandler {
	// Update model styles with user's theme preference
	themeName := user.Config.Theme
	if themeName == "" {
		themeName = "Magenta"
	}
	themeColor := GetThemeColor(themeName)
	m.styles = createStyles(m.termProfile, m.foregroundColor, themeColor)

	timerGen := words.NewGenerator()
	timerGen.Punctuation = user.Config.Punctuation

	wordGen := words.NewGenerator()
	wordGen.Punctuation = user.Config.Punctuation

	return &MainMenuHandler{
		BaseStateHandler: NewBaseStateHandler(StateMainMenu),
		MainMenuSelection: []string{
			"Timer",
			"Word Count",
			"Zen",
			"Config",
			"User Settings",
		},
		currentUser:            user,
		cursor:                 0,
		timerTestWordGenerator: timerGen,
		wordTestWordGenerator:  wordGen,
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
			case "User Settings":
				if h.ValidateTransition(StateUserSettings, context) {
					return NewUserSettingsHandler(context.model.session.User), nil
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

		case "ctrl+q":
			return NewPreAuthHandler(&context.model.context), nil
		}
	}
	return h, nil
}

func (h *MainMenuHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	displayName := m.session.User.Username
	if m.session.User.DisplayName != "" {
		displayName = m.session.User.DisplayName
	}
	termtyper := style("TermTyper - Welcome "+displayName, m.styles.themeFunc)
	termtyper = lipgloss.NewStyle().PaddingBottom(1).Render(termtyper)

	var menuItems []string
	menuItemsStyle := lipgloss.NewStyle().PaddingTop(1)

	for i, choice := range h.MainMenuSelection {
		choiceShow := style(choice, m.styles.toEnter)
		choiceShow = wrapWithCursor(h.cursor == i, choiceShow, m.styles.toEnter)
		choiceShow = menuItemsStyle.Render(choiceShow)
		menuItems = append(menuItems, choiceShow)
	}

	helpText := lipgloss.NewStyle().Faint(true).Render("\nenter: select, ctrl+q: logout")

	joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{termtyper}, menuItems...)...)
	joined = lipgloss.JoinVertical(lipgloss.Left, joined, helpText)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.DoubleBorder()).
		BorderForeground(lipgloss.Color("#FF00FF"))

	return borderStyle.Render(centeredText)
}
