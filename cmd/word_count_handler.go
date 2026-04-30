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

type WordCountTestHandler struct {
	*BaseStateHandler
	stopwatch StopWatch
	base      TestBase
	completed bool
}

func NewWordCountTestHandler(menu MainMenuHandler) *WordCountTestHandler {
	menu.wordTestWordGenerator.Count = menu.currentUser.Config.Words
	menu.wordTestWordGenerator.Punctuation = menu.currentUser.Config.Punctuation
	return &WordCountTestHandler{
		BaseStateHandler: NewBaseStateHandler(StateWordCountTest),
		stopwatch: StopWatch{
			stopwatch: stopwatch.New(),
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

		elapsedSeconds := h.stopwatch.Elapsed().Seconds()
		if int(elapsedSeconds) > len(h.base.wpmEachSecond) {
			elapsedMinutes := elapsedSeconds / 60.0
			if elapsedMinutes > 0 {
				h.base.wpmEachSecond = append(h.base.wpmEachSecond, h.base.calculateNormalizedWpm(elapsedMinutes))
			}
		}

	case tea.KeyPressMsg:
		switch msg.String() {
		case "esc":
			if h.ValidateTransition(StateMainMenu, context) {
				return NewMainMenuHandler(context.model.session.User, context.model), nil
			}
		case "ctrl+r":
			return NewWordCountTestHandler(h.base.mainMenu), nil

		case "backspace":
			handleBackspace(&h.base)
			recordInputBackspace(&h.base, h.stopwatch.Elapsed().Milliseconds())
		case "ctrl+t":
			handleCtrlBackspace(&h.base)
		default:
			if len(msg.Text) > 0 || msg.String() == "space" {
				if !h.stopwatch.isRunning {
					h.stopwatch.startTime = time.Now()
					commands = append(commands, h.stopwatch.stopwatch.Init())
					h.stopwatch.isRunning = true
				}

				handleCharacterInputFromMsg(msg, &h.base)
				recordInput(msg, &h.base, h.stopwatch.Elapsed().Milliseconds())
			}
		}
	}

	if len(h.base.wordsToEnter) == len(h.base.inputBuffer) &&
		!h.base.mistakes.mistakesAt[len(h.base.inputBuffer)-1] {
		results := h.calculateResults(context.model)
		return &results, tea.Batch(commands...)
	}

	return h, tea.Batch(commands...)
}

func (h *WordCountTestHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	s := ""
	stopwatchViewSeconds := strconv.FormatFloat(h.stopwatch.Elapsed().Seconds(), 'f', 0, 64) + "s"
	stopwatch := style(stopwatchViewSeconds, m.styles.themeFunc)
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

func (test WordCountTestHandler) calculateResults(m *model) ResultsHandler {
	elapsedMinutes := test.stopwatch.Elapsed().Minutes()
	wpm := test.base.calculateNormalizedWpm(elapsedMinutes)
	wpmChart := NewWPMChartBubble(m.width/2, m.height/2)
	wpmChart.UpdateData(test.base.wpmEachSecond)

	return ResultsHandler{
		testType:      "wordcount",
		wpm:           int(wpm),
		accuracy:      test.base.calculateAccuracy(),
		rawWpm:        int(test.base.calculateRawWpm(elapsedMinutes)),
		cpm:           test.base.calculateCpm(elapsedMinutes),
		time:          test.stopwatch.Elapsed(),
		test:          test.base,
		wpmEachSecond: test.base.wpmEachSecond,
		mainMenu:      test.base.mainMenu,
		resultsSelection: []string{
			"Next Test",
			"Main Menu",
			"Replay",
		},
		wpmChart: wpmChart,
	}
}
