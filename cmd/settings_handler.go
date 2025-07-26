package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
)

type SettingsHandler struct {
	*BaseStateHandler
	settings Settings
}

func NewSettingsHandler(settings Settings) *SettingsHandler {
	return &SettingsHandler{
		BaseStateHandler: NewBaseStateHandler(StateSettings),
		settings:         settings,
	}
}

func (h *SettingsHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if h.ValidateTransition(StateMainMenu, context) {
				return NewMainMenuHandler(initMainMenu(context.model.session.User)), nil
			}

		}
	}
	return h, nil
}

// Render implements StateHandler
func (h *SettingsHandler) Render(m *model) string {
	return h.settings.renderSettings(m)
}

// ValidateTransition implements StateHandler
func (h *SettingsHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateSettings]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}
