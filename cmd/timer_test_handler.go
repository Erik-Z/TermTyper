package cmd

import (
	"math"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TimerTestHandler struct {
	*BaseStateHandler
	test TimerTest
}

func NewTimerTestHandler(test TimerTest) *TimerTestHandler {
	return &TimerTestHandler{
		BaseStateHandler: NewBaseStateHandler(StateTimerTest),
		test:             test,
	}
}

func (h *TimerTestHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	//TODO: fix this
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

func (h *TimerTestHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2

	var timer string
	s := ""

	timer = style(h.test.timer.timer.View(), m.styles.magenta)

	paragraph := h.test.base.renderParagraph(lineLenLimit, m.styles)
	lines := strings.Split(paragraph, "\n")
	cursorLine := findCursorLine(lines, h.test.base.cursor)

	linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

	s += positionVertically(termHeight)
	avgLineLen := averageLineLen(lines)
	indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

	s += m.indent(timer, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

	if !h.test.timer.isRunning {
		s += "\n\n\n"
		s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
	}

	return s + "\n" + string(h.test.base.inputBuffer)
}

func (h *TimerTestHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateTimerTest]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}
