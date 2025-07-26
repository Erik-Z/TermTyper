package cmd

import (
	"math"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// WordCountTestHandler handles the word count test state
type WordCountTestHandler struct {
	*BaseStateHandler
	test WordCountTest
}

// NewWordCountTestHandler creates a new word count test handler
func NewWordCountTestHandler(test WordCountTest) *WordCountTestHandler {
	return &WordCountTestHandler{
		BaseStateHandler: NewBaseStateHandler(StateWordCountTest),
		test:             test,
	}
}

// HandleInput implements StateHandler
func (h *WordCountTestHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
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
func (h *WordCountTestHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	s := ""
	stopwatchViewSeconds := strconv.FormatFloat(h.test.stopwatch.stopwatch.Elapsed().Seconds(), 'f', 0, 64) + "s"
	stopwatch := style(stopwatchViewSeconds, m.styles.magenta)
	paragraphView := h.test.base.renderParagraph(lineLenLimit, m.styles)
	lines := strings.Split(paragraphView, "\n")
	cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), h.test.base.cursor)

	linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

	s += positionVertically(termHeight)
	avgLineLen := averageLineLen(lines)
	indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

	s += m.indent(stopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)
	s += "\n\n\n"
	s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))

	return s
}

// ValidateTransition implements StateHandler
func (h *WordCountTestHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateWordCountTest]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}
