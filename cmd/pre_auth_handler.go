package cmd

import (
	"reflect"
	"termtyper/database"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PreAuthHandler struct {
	*BaseStateHandler
	authMenu []State
	cursor   int
}

func (h *PreAuthHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	newCursor := h.cursor
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if _, ok := h.authMenu[newCursor].(GuestLogin); ok {
				guestUser := database.ApplicationUser{
					Id:       -1,
					Username: "Guest",
					Config: &database.UserConfig{
						Time:  30,
						Words: 30,
					},
				}
				return NewMainMenuHandler(&guestUser), nil
			}
			switch h.authMenu[newCursor].(type) {
			case Login:
				var commands []tea.Cmd
				newLoginHandler := NewLoginHandler("")
				initCmd := newLoginHandler.form.Init()
				commands = append(commands, initCmd)
				return newLoginHandler, tea.Batch(commands...)
			case Register:
				var commands []tea.Cmd
				newRegisterHandler := NewRegisterHandler(&context.model.context, "", "", "", "")
				initCmd := newRegisterHandler.form.Init()
				commands = append(commands, initCmd)
				return newRegisterHandler, tea.Batch(commands...)
			}
		case "up", "k":
			if h.cursor == 0 {
				newCursor = len(h.authMenu) - 1
			} else {
				newCursor--
			}
		case "down", "j":
			if h.cursor == len(h.authMenu)-1 {
				newCursor = 0
			} else {
				newCursor++
			}
		}
	}
	h.cursor = newCursor
	return h, nil
}

func (h *PreAuthHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	termtyper := style("TermTyper", m.styles.magenta)
	termtyper = lipgloss.NewStyle().PaddingBottom(1).Render(termtyper)

	var authMenu []string
	menuItemsStyle := lipgloss.NewStyle().PaddingTop(1)

	for i, choice := range h.authMenu {
		choiceShow := style(reflect.TypeOf(choice).Name(), m.styles.toEnter)

		choiceShow = wrapWithCursor(h.cursor == i, choiceShow, m.styles.toEnter)
		choiceShow = menuItemsStyle.Render(choiceShow)
		authMenu = append(authMenu, choiceShow)
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{termtyper}, authMenu...)...)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	return centeredText
}

func (h *PreAuthHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StatePreAuth]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}

func NewPreAuthHandler(context *database.Context) *PreAuthHandler {
	return &PreAuthHandler{
		BaseStateHandler: NewBaseStateHandler(StatePreAuth),
		authMenu: []State{
			Register{},
			Login{},
			GuestLogin{},
		},
		cursor: 0,
	}
}
