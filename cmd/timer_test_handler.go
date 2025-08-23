package cmd

import (
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TimerTestHandler struct {
	*BaseStateHandler
	base      TestBase
	completed bool
	timer     Timer
}

func NewTimerTestHandler(menu MainMenuHandler) *TimerTestHandler {
	testDuration := time.Duration(menu.currentUser.Config.Time) * time.Second
	return &TimerTestHandler{
		BaseStateHandler: NewBaseStateHandler(StateTimerTest),
		timer: Timer{
			timer:     timer.NewWithInterval(testDuration, time.Second),
			duration:  testDuration,
			isRunning: false,
			timedout:  false,
		},
		base: TestBase{
			wordsToEnter:  menu.timerTestWordGenerator.Generate("Common words"),
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

func (h *TimerTestHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	var commands []tea.Cmd
	switch msg := msg.(type) {
	case timer.TickMsg:
		timerUpdate, cmdUpdate := h.timer.timer.Update(msg)
		h.timer.timer = timerUpdate
		commands = append(commands, cmdUpdate)

		elapsedMinute := h.timer.duration.Minutes() - h.timer.timer.Timeout.Minutes()
		if elapsedMinute != 0 {
			h.base.wpmEachSecond = append(h.base.wpmEachSecond, h.base.calculateNormalizedWpm(elapsedMinute))
		}

		if h.timer.timer.Timedout() {
			h.timer.timedout = true

			results := h.calculateResults()
			return &results, tea.Batch(commands...)
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			if h.ValidateTransition(StateMainMenu, context) {
				return NewMainMenuHandler(context.model.session.User), nil
			}

		case "backspace":
			handleBackspace(&h.base)
		case "ctrl+t":
			// Delete entire word
			handleCtrlBackspace(&h.base)
		case "ctrl+w":
			return NewMainMenuHandler(context.model.session.User), nil
		case "ctrl+r":
			return NewTimerTestHandler(h.base.mainMenu), nil
		default:
			switch msg.Type {
			case tea.KeyRunes, tea.KeySpace:
				if !h.timer.isRunning {
					commands = append(commands, h.timer.timer.Init())
					h.timer.isRunning = true
				}

				handleCharacterInputFromMsg(msg, &h.base)
			}
		}
	}

	return h, tea.Batch(commands...)
}

func (h *TimerTestHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2

	timer := style(h.timer.timer.View(), m.styles.magenta)
	s := ""

	paragraph := h.base.renderParagraph(lineLenLimit, m.styles)
	lines := strings.Split(paragraph, "\n")
	cursorLine := findCursorLine(lines, h.base.cursor)

	linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

	s += positionVertically(termHeight)
	avgLineLen := averageLineLen(lines)
	indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

	s += m.indent(timer, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

	if !h.timer.isRunning {
		s += "\n\n\n"
		s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
	}

	return s + "\n" //+ string(h.base.inputBuffer)
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

func (test TimerTestHandler) calculateResults() ResultsHandler {
	elapsedMinutes := test.timer.duration.Minutes()
	wpm := test.base.calculateNormalizedWpm(elapsedMinutes)

	return ResultsHandler{
		testType:      "timer",
		wpm:           int(wpm),
		accuracy:      test.base.calculateAccuracy(),
		rawWpm:        int(test.base.calculateRawWpm(elapsedMinutes)),
		cpm:           test.base.calculateCpm(elapsedMinutes),
		time:          test.timer.duration,
		wpmEachSecond: test.base.wpmEachSecond,
		mainMenu:      test.base.mainMenu,
		resultsSelection: []string{
			"Next Test",
			"Main Menu",
			"Replay",
		},
	}
}
