package cmd

import (
	"math"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ReplayHandler struct {
	*BaseStateHandler
	replay Replay
}

func NewReplayHandler(replay Replay) *ReplayHandler {
	return &ReplayHandler{
		BaseStateHandler: NewBaseStateHandler(StateReplay),
		replay:           replay,
	}
}

// HandleInput implements StateHandler
func (h *ReplayHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if h.ValidateTransition(StateResults, context) {
				//return NewResultsHandler(initResults(h.replay.results)), nil
			}

		}
	}
	return h, nil
}

func (h *ReplayHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	s := ""

	stopwatchViewSeconds := strconv.FormatFloat(h.replay.stopwatch.stopwatch.Elapsed().Seconds(), 'f', 0, 64) + "s"
	stopwatch := style(stopwatchViewSeconds, m.styles.magenta)
	paragraphView := h.replay.test.renderParagraph(lineLenLimit, m.styles)
	lines := strings.Split(paragraphView, "\n")
	cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), h.replay.test.cursor)

	linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

	s += positionVertically(termHeight)
	avgLineLen := averageLineLen(lines)
	indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

	s += m.indent(stopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)
	s += "\n\n\n"

	if h.replay.isReplayInProcess {
		s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("Replay in progress..", m.styles.toEnter))

	} else {
		s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("Press any key to start Replay", m.styles.toEnter))
	}

	return s
}

// ValidateTransition implements StateHandler
func (h *ReplayHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateReplay]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}
