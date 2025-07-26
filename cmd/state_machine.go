package cmd

import (
	tea "github.com/charmbracelet/bubbletea"
)

// StateType represents all possible states in the application
type StateType int

const (
	StatePreAuth StateType = iota
	StateLogin
	StateRegister
	StateMainMenu
	StateTimerTest
	StateZenMode
	StateWordCountTest
	StateResults
	StateSettings
	StateReplay
)

type StateTransition struct {
	From StateType
	To   StateType
}

type StateMachine struct {
	currentState StateType
	model        *model
	transitions  map[StateType][]StateType
	handlers     map[StateType]StateHandler
}

func NewStateMachine(m *model) *StateMachine {
	sm := &StateMachine{
		model: m,
		transitions: map[StateType][]StateType{
			StatePreAuth: {
				StateLogin,
				StateRegister,
				StateMainMenu,
			},
			StateLogin: {
				StateMainMenu,
				StatePreAuth,
			},
			StateRegister: {
				StateMainMenu,
				StatePreAuth,
			},
			StateMainMenu: {
				StateTimerTest,
				StateZenMode,
				StateWordCountTest,
				StateSettings,
			},
			StateTimerTest: {
				StateResults,
				StateMainMenu,
			},
			StateZenMode: {
				StateMainMenu,
			},
			StateWordCountTest: {
				StateResults,
				StateMainMenu,
			},
			StateResults: {
				StateMainMenu,
				StateReplay,
				StateTimerTest,
				StateWordCountTest,
			},
			StateSettings: {
				StateMainMenu,
			},
			StateReplay: {
				StateMainMenu,
			},
		},
		handlers: make(map[StateType]StateHandler),
	}

	sm.handlers[StatePreAuth] = &PreAuthHandler{}
	sm.handlers[StateLogin] = &LoginHandler{}
	sm.handlers[StateRegister] = &RegisterHandler{}
	sm.handlers[StateMainMenu] = &MainMenuHandler{}
	sm.handlers[StateTimerTest] = &TimerTestHandler{}
	sm.handlers[StateZenMode] = &ZenModeHandler{}
	sm.handlers[StateWordCountTest] = &WordCountTestHandler{}
	sm.handlers[StateResults] = &ResultsHandler{}
	sm.handlers[StateSettings] = &SettingsHandler{}
	sm.handlers[StateReplay] = &ReplayHandler{}

	return sm
}

func (sm *StateMachine) Transition(to StateType) bool {
	validTransitions := sm.transitions[sm.currentState]
	for _, validState := range validTransitions {
		if validState == to {
			sm.currentState = to
			return true
		}
	}
	return false
}

func (sm *StateMachine) HandleInput(msg tea.Msg) (tea.Model, tea.Cmd) {
	handler := sm.handlers[sm.currentState]
	if handler == nil {
		return sm.model, nil
	}

	newHandler, cmd := handler.HandleInput(msg, &StateContext{
		model:         sm.model,
		transitionMap: sm.transitions,
	})

	if newHandler != nil {
		sm.handlers[sm.currentState] = newHandler
	}

	return sm.model, cmd
}

func (sm *StateMachine) Render() string {
	handler := sm.handlers[sm.currentState]
	if handler == nil {
		return ""
	}
	return handler.Render(sm.model)
}

func (sm *StateMachine) GetCurrentState() StateType {
	return sm.currentState
}

func (sm *StateMachine) SetCurrentState(state StateType) {
	sm.currentState = state
}
