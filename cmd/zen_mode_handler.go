package cmd

import (
	"strings"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ZenModeHandler struct {
	*BaseStateHandler
	base      TestBase
	stopwatch StopWatch
}

func NewZenModeHandler(menu MainMenuHandler) *ZenModeHandler {
	return &ZenModeHandler{
		BaseStateHandler: NewBaseStateHandler(StateZenMode),
		stopwatch: StopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter:  make([]rune, 0),
			inputBuffer:   make([]rune, 0),
			rawInputCount: 0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor:   0,
			mainMenu: menu,
		},
	}
}

func (h *ZenModeHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
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
		case "ctrl+w":
			return NewMainMenuHandler(context.model.session.User), nil
		case "ctrl+r":
			return NewZenModeHandler(h.base.mainMenu), nil
		case "ctrl+backspace":
			handleCtrlBackspace(&h.base)
		case "backspace":
			handleBackspace(&h.base)
		default:
			switch msg.Type {
			case tea.KeyRunes:
				if !h.stopwatch.isRunning {
					commands = append(commands, h.stopwatch.stopwatch.Init())
					h.stopwatch.isRunning = true
				}

				handleCharacterInputZenMode(msg, &h.base)
			}
		}
	}
	return h, tea.Batch(commands...)
}

func (h *ZenModeHandler) Render(m *model) string {
	s := ""
	termWidth, termHeight := m.width-2, m.height-2

	stopwatch := style(h.stopwatch.stopwatch.View(), m.styles.magenta)
	paragraphView := h.base.renderParagraphZenMode(lineLenLimit, m.styles)
	lines := strings.Split(paragraphView, "\n")

	cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), h.base.cursor)
	linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")
	s += positionVertically(termHeight)

	s += stopwatch + "\n\n" + linesAroundCursor
	if !h.stopwatch.isRunning {
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
