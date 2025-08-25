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

type WordCountTestHandler struct {
	*BaseStateHandler
	stopwatch StopWatch
	base      TestBase
	completed bool
}

func NewWordCountTestHandler(menu MainMenuHandler) *WordCountTestHandler {
	menu.wordTestWordGenerator.Count = menu.currentUser.Config.Words
	return &WordCountTestHandler{
		BaseStateHandler: NewBaseStateHandler(StateWordCountTest),
		stopwatch: StopWatch{
			stopwatch: stopwatch.NewWithInterval(time.Millisecond * 100),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter:  menu.wordTestWordGenerator.Generate("Common words"),
			inputBuffer:   make([]rune, 0),
			rawInputCount: 0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor:   0,
			mainMenu: menu,
		},
		completed: false,
	}
}

func (h *WordCountTestHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
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

		elapsedMinutes := h.stopwatch.stopwatch.Elapsed().Minutes()
		if elapsedMinutes != 0 {
			h.base.wpmEachSecond = append(h.base.wpmEachSecond, h.base.calculateNormalizedWpm(elapsedMinutes))
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if h.ValidateTransition(StateMainMenu, context) {
				return NewMainMenuHandler(context.model.session.User), nil
			}
		case "ctrl+r":
			return NewWordCountTestHandler(h.base.mainMenu), nil

		case "backspace":
			handleBackspace(&h.base)
			recordInput(msg, h)
		case "ctrl+t":
			// Delete entire word
			handleCtrlBackspace(&h.base)
		default:
			switch msg.Type {
			case tea.KeyRunes, tea.KeySpace:
				if !h.stopwatch.isRunning {
					commands = append(commands, h.stopwatch.stopwatch.Init())
					h.stopwatch.isRunning = true
				}
				handleCharacterInputFromMsg(msg, &h.base)
				recordInput(msg, h)
			}
		}
	}

	if len(h.base.wordsToEnter) == len(h.base.inputBuffer) &&
		!h.base.mistakes.mistakesAt[len(h.base.inputBuffer)-1] {
		//termenv.DefaultOutput().Reset()
		results := h.calculateResults()
		return &results, tea.Batch(commands...)
	}

	return h, tea.Batch(commands...)
}

func (h *WordCountTestHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	s := ""
	stopwatchViewSeconds := strconv.FormatFloat(h.stopwatch.stopwatch.Elapsed().Seconds(), 'f', 0, 64) + "s"
	stopwatch := style(stopwatchViewSeconds, m.styles.magenta)
	paragraphView := h.base.renderParagraph(lineLenLimit, m.styles)
	lines := strings.Split(paragraphView, "\n")
	cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), h.base.cursor)

	linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

	s += positionVertically(termHeight)
	avgLineLen := averageLineLen(lines)
	indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

	s += m.indent(stopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)
	s += "\n\n\n"
	s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))

	return s
}

func (h *WordCountTestHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateWordCountTest]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}

func (test WordCountTestHandler) calculateResults() ResultsHandler {
	elapsedMinutes := test.stopwatch.stopwatch.Elapsed().Minutes()
	wpm := test.base.calculateNormalizedWpm(elapsedMinutes)

	return ResultsHandler{
		testType:      "wordcount",
		wpm:           int(wpm),
		accuracy:      test.base.calculateAccuracy(),
		rawWpm:        int(test.base.calculateRawWpm(elapsedMinutes)),
		cpm:           test.base.calculateCpm(elapsedMinutes),
		time:          test.stopwatch.stopwatch.Elapsed(),
		test:          test.base,
		wpmEachSecond: test.base.wpmEachSecond,
		mainMenu:      test.base.mainMenu,
		resultsSelection: []string{
			"Next Test",
			"Main Menu",
			"Replay",
		},
	}
}
