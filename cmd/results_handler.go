package cmd

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ResultsHandler handles displaying test results
type ResultsHandler struct {
	*BaseStateHandler
	results Results
}

// NewResultsHandler creates a new results handler
func NewResultsHandler(results Results) *ResultsHandler {
	return &ResultsHandler{
		BaseStateHandler: NewBaseStateHandler(StateResults),
		results:          results,
	}
}

func (h *ResultsHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	//TODO: Fix this
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			return NewMainMenuHandler(initMainMenu(context.model.session.User)), nil
		}
	}
	return h, nil
}

func (h *ResultsHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2

	title := style("Test Results", m.styles.magenta)
	title = lipgloss.NewStyle().PaddingBottom(1).Render(title)

	var content []string

	content = append(content, fmt.Sprintf("WPM: %d", h.results.wpm))
	content = append(content, fmt.Sprintf("Accuracy: %.1f%%", h.results.accuracy))

	// if h.results.testType == "timer" {
	// 	content = append(content, fmt.Sprintf("Time: %s", formatDuration(h.results.duration)))
	// 	content = append(content, fmt.Sprintf("Words: %d", h.results.wordsTyped))
	// } else if h.results.testType == "wordcount" {
	// 	content = append(content, fmt.Sprintf("Words: %d/%d", h.results.wordsTyped, h.results.targetWords))
	// 	content = append(content, fmt.Sprintf("Time: %s", formatDuration(h.results.duration)))
	// }

	content = append(content, "\nPress Enter to replay the test")

	content = append(content, "\nPress ESC to return to main menu")

	joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{title}, content...)...)
	s := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, joined)

	return s
}

func (h *ResultsHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateResults]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}
