package cmd

import (
	"fmt"
	"strings"

	"termtyper/database"

	tea "charm.land/bubbletea/v2"
	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

type UserSettingsHandler struct {
	*BaseStateHandler
	form     *huh.Form
	formData *UserSettingsFormData
	user     *database.ApplicationUser
}

type UserSettingsFormData struct {
	displayName string
}

func NewUserSettingsHandler(user *database.ApplicationUser) *UserSettingsHandler {
	data := &UserSettingsFormData{
		displayName: user.DisplayName,
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Display Name (max 20 chars)").
				Value(&data.displayName).
				Validate(func(str string) error {
					if len(str) > 20 {
						return fmt.Errorf("display name must be 20 characters or less")
					}
					if containsProfanity(str) {
						return fmt.Errorf("display name contains inappropriate language")
					}
					return nil
				}),
		),
	)

	return &UserSettingsHandler{
		BaseStateHandler: NewBaseStateHandler(StateUserSettings),
		form:             form,
		formData:         data,
		user:             user,
	}
}

func (h *UserSettingsHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	var commands []tea.Cmd

	updatedForm, formCmd := h.form.Update(msg)
	if f, ok := updatedForm.(*huh.Form); ok {
		h.form = f
		commands = append(commands, formCmd)
	}

	if h.form.State == huh.StateCompleted {
		_, err := context.model.context.UserRepository.Exec(
			"UPDATE users SET display_name = ? WHERE id = ?",
			h.formData.displayName, h.user.Id,
		)
		if err == nil {
			h.user.DisplayName = h.formData.displayName
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+q":
			return NewMainMenuHandler(h.user, context.model), nil
		}
	}

	return h, tea.Batch(commands...)
}

func (h *UserSettingsHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2

	title := style("User Settings", m.styles.themeFunc)
	title = lipgloss.NewStyle().PaddingBottom(1).Render(title)

	helpText := lipgloss.NewStyle().Faint(true).Render("enter: edit/save • esc/ctrl+q: back")

	joined := lipgloss.JoinVertical(lipgloss.Left, title, h.form.View(), "", helpText)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	return centeredText
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

// Basic profanity filter - common bad words
var profanityList = []string{
	"fuck", "shit", "ass", "bitch", "damn", "hell", "crap", "piss",
	"fucker", "motherfucker", "bastard", "douch", "douchebag",
	"asshole", "cunt", "dick", "pussy", "cock", "tit", "tits",
	"whore", "slut", "fag", "faggot", "nigger", "nigga",
}

func containsProfanity(s string) bool {
	lower := strings.ToLower(s)
	for _, word := range profanityList {
		if strings.Contains(lower, word) {
			return true
		}
	}
	return false
}
