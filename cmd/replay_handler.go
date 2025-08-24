package cmd

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ReplayHandler struct {
	*BaseStateHandler
	test              TestBase
	results           ResultsHandler
	stopwatch         StopWatch
	isReplayInProcess bool
}

func NewReplayHandler(results ResultsHandler) *ReplayHandler {
	return &ReplayHandler{
		BaseStateHandler:  NewBaseStateHandler(StateReplay),
		test:              results.test,
		results:           results,
		isReplayInProcess: false,
		stopwatch: StopWatch{
			stopwatch: stopwatch.NewWithInterval(time.Millisecond),
			isRunning: false,
		},
	}
}

func (h *ReplayHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	var commands []tea.Cmd
	switch msg := msg.(type) {
	case stopwatch.StartStopMsg:
		stopwatchUpdate, cmdUpdate := h.stopwatch.stopwatch.Update(msg)
		h.stopwatch.stopwatch = stopwatchUpdate
		commands = append(commands, cmdUpdate)

	case stopwatch.TickMsg:
		stopwatchUpdate, cmdUpdate := h.stopwatch.stopwatch.Update(msg)
		h.stopwatch.stopwatch = stopwatchUpdate
		commands = append(commands, cmdUpdate)

	case tea.KeyMsg:
		if !h.isReplayInProcess {
			h.isReplayInProcess = true
		}
	}

	if h.isReplayInProcess {
		if !h.stopwatch.isRunning && len(h.test.testRecord) > 0 {
			commands = append(commands, h.stopwatch.stopwatch.Init())
			h.stopwatch.isRunning = true
		}

		if len(h.test.testRecord) > 0 {
			currentKeyPress := h.test.testRecord[0]
			if currentKeyPress.timestamp <= h.stopwatch.stopwatch.Elapsed().Milliseconds() {
				switch currentKeyPress.key {
				case '\b':
					handleBackspace(&h.test)
				default:
					handleCharacterInputFromRune(currentKeyPress.key, &h.test)
				}
				h.test.testRecord = h.test.testRecord[1:]
			}
		}

		if len(h.test.wordsToEnter) == len(h.test.inputBuffer) &&
			!h.test.mistakes.mistakesAt[len(h.test.inputBuffer)-1] {

			commands = append(commands, h.stopwatch.stopwatch.Stop())
			h.stopwatch.isRunning = false
		}

		// case tea.KeyMsg:
		// 	switch msg.String() {
		// 	case "esc":
		// 		if h.ValidateTransition(StateResults, context) {
		// 			//return NewResultsHandler(initResults(h.replay.results)), nil
		// 		}

		// 	}
	}
	return h, tea.Batch(commands...)
}

func (h *ReplayHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	s := ""

	stopwatchViewSeconds := strconv.FormatFloat(h.stopwatch.stopwatch.Elapsed().Seconds(), 'f', 0, 64) + "s"
	stopwatch := style(stopwatchViewSeconds, m.styles.magenta)
	paragraphView := h.test.renderParagraph(lineLenLimit, m.styles)
	lines := strings.Split(paragraphView, "\n")
	cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), h.test.cursor)

	linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

	s += positionVertically(termHeight)
	avgLineLen := averageLineLen(lines)
	indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

	s += m.indent(stopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)
	s += "\n\n\n"

	if h.isReplayInProcess {
		s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("Replay in progress..", m.styles.toEnter))
	} else {
		s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("Press any key to start Replay", m.styles.toEnter))
	}

	return s
}

func (h *ReplayHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateReplay]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}
