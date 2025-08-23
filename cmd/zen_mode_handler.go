package cmd

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ZenModeHandler handles the zen mode state
type ZenModeHandler struct {
	*BaseStateHandler
	zen ZenMode
}

func NewZenModeHandler(zen ZenMode) *ZenModeHandler {
	return &ZenModeHandler{
		BaseStateHandler: NewBaseStateHandler(StateZenMode),
		zen:              zen,
	}
}

func (h *ZenModeHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if h.ValidateTransition(StateMainMenu, context) {
				return NewMainMenuHandler(context.model.session.User), nil
			}
			// 	case "enter":
			// 		if !h.zen.started {
			// 			h.zen.started = true
			// 			h.zen.stopwatch.Start()
			// 			return h, h.zen.stopwatch.Tick()
			// 		}
			// 	case "tab":
			// 		if h.zen.started {
			// 			h.zen.currentWordIndex++
			// 			if h.zen.currentWordIndex >= len(h.zen.words) {
			// 				h.zen.finished = true
			// 				h.zen.stopwatch.Stop()
			// 				if h.ValidateTransition(StateResults, context) {
			// 					return NewResultsHandler(initResults(h.zen)), nil
			// 				}
			// 			}
			// 		}
			// 	default:
			// 		if h.zen.started && !h.zen.finished {
			// 			// Handle typing input
			// 			h.zen.currentInput += msg.String()
			// 			if len(h.zen.currentInput) > len(h.zen.words[h.zen.currentWordIndex]) {
			// 				h.zen.currentInput = h.zen.currentInput[:len(h.zen.words[h.zen.currentWordIndex])]
			// 			}
			// 		}
			// 	}
			// case stopwatch.TickMsg:
			// 	if h.zen.started && !h.zen.finished {
			// 		return h, h.zen.stopwatch.Tick()
		}
	}
	return h, nil
}

// Render implements StateHandler
func (h *ZenModeHandler) Render(m *model) string {
	s := ""
	termWidth, termHeight := m.width-2, m.height-2

	stopwatch := style(h.zen.stopwatch.stopwatch.View(), m.styles.magenta)
	paragraphView := h.zen.base.renderParagraphZenMode(lineLenLimit, m.styles)
	lines := strings.Split(paragraphView, "\n")

	cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), h.zen.base.cursor)
	linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")
	s += positionVertically(termHeight)

	s += stopwatch + "\n\n" + linesAroundCursor
	if !h.zen.stopwatch.isRunning {
		s += "\n\n\n"
		s += style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter)
		//s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
	}
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center+lipgloss.Position(termHeight/2), s)

	return centeredText
}

func (h *ZenModeHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateZenMode]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}
