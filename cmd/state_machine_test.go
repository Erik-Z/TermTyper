package cmd

import (
	"testing"
)

func TestNewStateMachine(t *testing.T) {
	m := &model{}
	sm := NewStateMachine(m)

	if sm == nil {
		t.Fatal("state machine should not be nil")
	}

	if sm.model != m {
		t.Error("state machine should have the correct model")
	}

	if sm.currentState != 0 {
		t.Error("initial state should be 0 (StatePreAuth)")
	}

	// Check that all handlers are initialized
	expectedHandlers := []StateType{
		StatePreAuth, StateLogin, StateRegister, StateMainMenu,
		StateTimerTest, StateZenMode, StateWordCountTest,
		StateResults, StateSettings, StateReplay,
	}

	for _, stateType := range expectedHandlers {
		if sm.handlers[stateType] == nil {
			t.Errorf("handler for state %d should not be nil", stateType)
		}
	}
}

func TestStateMachineTransition(t *testing.T) {
	m := &model{}
	sm := NewStateMachine(m)

	tests := []struct {
		name           string
		fromState      StateType
		toState        StateType
		expectedResult bool
	}{
		{
			name:           "valid transition from pre-auth to login",
			fromState:      StatePreAuth,
			toState:        StateLogin,
			expectedResult: true,
		},
		{
			name:           "valid transition from pre-auth to register",
			fromState:      StatePreAuth,
			toState:        StateRegister,
			expectedResult: true,
		},
		{
			name:           "valid transition from pre-auth to main menu",
			fromState:      StatePreAuth,
			toState:        StateMainMenu,
			expectedResult: true,
		},
		{
			name:           "invalid transition from pre-auth to timer test",
			fromState:      StatePreAuth,
			toState:        StateTimerTest,
			expectedResult: false,
		},
		{
			name:           "valid transition from main menu to timer test",
			fromState:      StateMainMenu,
			toState:        StateTimerTest,
			expectedResult: true,
		},
		{
			name:           "valid transition from timer test to results",
			fromState:      StateTimerTest,
			toState:        StateResults,
			expectedResult: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sm.SetCurrentState(tt.fromState)
			result := sm.Transition(tt.toState)

			if result != tt.expectedResult {
				t.Errorf("expected transition result %v, got %v", tt.expectedResult, result)
			}

			if result {
				if sm.GetCurrentState() != tt.toState {
					t.Errorf("state should be updated to %d, got %d", tt.toState, sm.GetCurrentState())
				}
			} else {
				if sm.GetCurrentState() != tt.fromState {
					t.Errorf("state should remain %d, got %d", tt.fromState, sm.GetCurrentState())
				}
			}
		})
	}
}

func TestStateMachineGetSetCurrentState(t *testing.T) {
	m := &model{}
	sm := NewStateMachine(m)

	if sm.GetCurrentState() != StatePreAuth {
		t.Errorf("initial state should be StatePreAuth, got %d", sm.GetCurrentState())
	}

	sm.SetCurrentState(StateMainMenu)
	if sm.GetCurrentState() != StateMainMenu {
		t.Errorf("state should be StateMainMenu, got %d", sm.GetCurrentState())
	}

	sm.SetCurrentState(StateTimerTest)
	if sm.GetCurrentState() != StateTimerTest {
		t.Errorf("state should be StateTimerTest, got %d", sm.GetCurrentState())
	}
}

func TestStateMachineHandleInput(t *testing.T) {
	m := &model{}
	sm := NewStateMachine(m)

	// Test with nil message
	model, cmd := sm.HandleInput(nil)
	if model != m {
		t.Error("should return the same model")
	}
	if cmd != nil {
		t.Error("should return nil command for nil message")
	}

	// Test with a key message - the PreAuthHandler needs authMenu to be set up
	// Since we're testing StateMachine, not the handler, we'll test with a different approach
	// Just verify that HandleInput doesn't panic with a valid message
	// The PreAuthHandler in NewStateMachine is created with empty authMenu
	// So we'll test with a message that doesn't access authMenu

	// Skip the enter key test since it requires proper authMenu initialization
	// The real initialization happens in NewPreAuthHandler which requires database context
}

func TestStateMachineRender(t *testing.T) {
	// Initialize model with required fields
	themeColor := "#FF00FF"

	m := &model{
		width:  80,
		height: 24,
		styles: createStyles(themeColor),
	}
	sm := NewStateMachine(m)

	// Test render - at StatePreAuth, it should not panic
	result := sm.Render()
	if result == "" {
		t.Error("render should not return empty string")
	}

	// Skip MainMenu render test as it requires complex handler initialization
	// The render test for MainMenu is covered in integration tests
}
