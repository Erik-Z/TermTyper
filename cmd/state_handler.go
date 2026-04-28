package cmd

import (
	tea "charm.land/bubbletea/v2"
)

type StateHandler interface {
	// HandleInput processes input for the current state and returns a new handler and command if needed
	HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd)

	Render(m *model) string

	GetStateType() StateType

	ValidateTransition(to StateType, context *StateContext) bool
}

type StateContext struct {
	model         *model
	transitionMap map[StateType][]StateType
}

type BaseStateHandler struct {
	stateType StateType
}

func (h *BaseStateHandler) GetStateType() StateType {
	return h.stateType
}

func (h *BaseStateHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[h.stateType]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}

func (h *BaseStateHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	return nil, nil // Base implementation does nothing, should be overridden
}

func (h *BaseStateHandler) Render(m *model) string {
	return "" // Base implementation returns empty string, should be overridden
}

func NewBaseStateHandler(stateType StateType) *BaseStateHandler {
	return &BaseStateHandler{
		stateType: stateType,
	}
}
