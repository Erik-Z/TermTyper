package cmd

import (
	"math"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/stopwatch"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type ReplayHandler struct {
	*BaseStateHandler
	test              TestBase
	results           ResultsHandler
	stopwatch         StopWatch
	isReplayInProcess bool
	replayDone        bool
	replaySelection   []string
	replayCursor      int
}

func NewReplayHandler(results ResultsHandler) *ReplayHandler {
	results.test.inputBuffer = make([]rune, 0)
	results.test.cursor = 0
	return &ReplayHandler{
		BaseStateHandler:  NewBaseStateHandler(StateReplay),
		test:              results.test,
		results:           results,
		isReplayInProcess: false,
		replayDone:        false,
		replaySelection:   []string{"New Test", "Main Menu", "Replay"},
		replayCursor:      0,
		stopwatch: StopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
	}
}

func (h *ReplayHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	var commands []tea.Cmd

	if h.replayDone {
		switch msg := msg.(type) {
		case tea.KeyPressMsg:
			switch msg.String() {
			case "enter":
				switch h.replaySelection[h.replayCursor] {
				case "Replay":
					return NewReplayHandler(h.results), nil
				case "New Test":
					return NewWordCountTestHandler(h.results.mainMenu), nil
				case "Main Menu":
					return NewMainMenuHandler(context.model.session.User, context.model), nil
				}
			case "left", "h":
				if h.replayCursor == 0 {
					h.replayCursor = len(h.replaySelection) - 1
				} else {
					h.replayCursor--
				}
			case "right", "l":
				if h.replayCursor == len(h.replaySelection)-1 {
					h.replayCursor = 0
				} else {
					h.replayCursor++
				}
			case "esc":
				return NewMainMenuHandler(context.model.session.User, context.model), nil
			}
		}
		return h, nil
	}

	switch msg := msg.(type) {
	case stopwatch.StartStopMsg:
		stopwatchUpdate, cmdUpdate := h.stopwatch.stopwatch.Update(msg)
		h.stopwatch.stopwatch = stopwatchUpdate
		commands = append(commands, cmdUpdate)

	case stopwatch.TickMsg:
		stopwatchUpdate, cmdUpdate := h.stopwatch.stopwatch.Update(msg)
		h.stopwatch.stopwatch = stopwatchUpdate
		commands = append(commands, cmdUpdate)

	case tea.KeyPressMsg:
		if !h.isReplayInProcess {
			h.isReplayInProcess = true
		}
	}

	if h.isReplayInProcess {
		if !h.stopwatch.isRunning && len(h.test.testRecord) > 0 {
			h.stopwatch.startTime = time.Now()
			commands = append(commands, h.stopwatch.stopwatch.Init())
			h.stopwatch.isRunning = true
		}

		if len(h.test.testRecord) > 0 {
			currentKeyPress := h.test.testRecord[0]
			if currentKeyPress.timestamp <= h.stopwatch.Elapsed().Milliseconds() {
				switch currentKeyPress.key {
				case '\b':
					handleBackspace(&h.test)
				default:
					handleCharacterInputFromRune(currentKeyPress.key, &h.test)
				}
				h.test.testRecord = h.test.testRecord[1:]
			}
		}

		if len(h.test.testRecord) == 0 {
			h.isReplayInProcess = false
			h.replayDone = true

			commands = append(commands, h.stopwatch.stopwatch.Stop())
			h.stopwatch.isRunning = false
		}
	}
	return h, tea.Batch(commands...)
}

func (h *ReplayHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	s := ""

	stopwatchViewSeconds := strconv.FormatFloat(h.stopwatch.Elapsed().Seconds(), 'f', 0, 64) + "s"
	stopwatch := style(stopwatchViewSeconds, m.styles.themeFunc)
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
	} else if h.replayDone {
		// TODO: replay menu is inconsistent with other menus, needs to be fixed
		var menuItems []string
		for i, choice := range h.replaySelection {
			choiceShow := style(choice, m.styles.toEnter)
			if i == h.replayCursor {
				choiceShow = style(choice, m.styles.themeFunc)
			}
			menuItems = append(menuItems, choiceShow)
		}
		s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, strings.Join(menuItems, " | "))
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
