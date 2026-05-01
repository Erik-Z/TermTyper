package cmd

import (
	"fmt"

	"termtyper/database"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type UserSettingsHandler struct {
	*BaseStateHandler
	displayName    string
	savedValue    string
	editing       bool
	cursor         int
	user           *database.ApplicationUser
}

func NewUserSettingsHandler(user *database.ApplicationUser) *UserSettingsHandler {
	return &UserSettingsHandler{
		BaseStateHandler: NewBaseStateHandler(StateUserSettings),
		displayName:    user.DisplayName,
		savedValue:    user.DisplayName,
		editing:       false,
		cursor:         0,
		user:           user,
	}
}

func (h *UserSettingsHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+q":
			if h.displayName != h.savedValue {
				// Save before leaving
				h.saveDisplayName(context)
			}
			return NewMainMenuHandler(h.user, context.model), nil

		case "enter":
			h.editing = !h.editing

		case "backspace":
			if h.editing && h.cursor > 0 {
				h.displayName = h.displayName[:h.cursor-1] + h.displayName[h.cursor:]
				h.cursor--
			}

		case "left", "h":
			if h.editing && h.cursor > 0 {
				h.cursor--
			}

		case "right", "l":
			if h.editing && h.cursor < len(h.displayName) {
				h.cursor++
			}

		default:
			if h.editing {
				if keyMsg, ok := msg.(tea.KeyPressMsg); ok && len(keyMsg.Text) > 0 {
					runeChar := keyMsg.Text[0]
					if runeChar >= 32 && runeChar <= 126 && len(h.displayName) < 20 {
						h.displayName = h.displayName[:h.cursor] + string(runeChar) + h.displayName[h.cursor:]
						h.cursor++
					}
				}
			}
		}
	}
	return h, nil
}

func (h *UserSettingsHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2

	title := style("User Settings", m.styles.themeFunc)
	title = lipgloss.NewStyle().PaddingBottom(1).Render(title)

	value := h.displayName
	if value == "" {
		value = "(not set)"
	}

	cursorChar := ""
	if h.editing {
		cursorChar = style("|", m.styles.themeFunc)
	}

	status := fmt.Sprintf("%s%s", value, cursorChar)
	if h.displayName == h.savedValue {
		status = style(status, m.styles.themeFunc)
	} else {
		status = style(status, m.styles.toEnter)
	}

	displayRow := fmt.Sprintf("Display Name: [%s]", status)
	helpText := lipgloss.NewStyle().Faint(true).Render("enter: edit/save • esc/ctrl+q: back")

	joined := lipgloss.JoinVertical(lipgloss.Left, title, displayRow, "", helpText)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	return centeredText
}

func (h *UserSettingsHandler) saveDisplayName(context *StateContext) {
	_, err := context.model.context.UserRepository.Exec(
		"UPDATE users SET display_name = ? WHERE id = ?",
		h.displayName, h.user.Id,
	)
	if err == nil {
		h.savedValue = h.displayName
		h.user.DisplayName = h.displayName
	}
}

func (h *UserSettingsHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[h.GetStateType()]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}
