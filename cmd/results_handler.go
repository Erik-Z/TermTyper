package cmd

import (
	"fmt"
	"math"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResultsHandler struct {
	*BaseStateHandler
	testType string
	wpm      int
	accuracy float64
	//deltaWpm         float64
	rawWpm int
	cpm    int
	time   time.Duration
	//wordList         string
	test             TestBase
	wpmEachSecond    []float64
	mainMenu         MainMenuHandler
	resultsSelection []string
	cursor           int
}

func NewResultsHandler() *ResultsHandler {
	return &ResultsHandler{
		BaseStateHandler: NewBaseStateHandler(StateResults),
	}
}

func (h *ResultsHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	//TODO: Fix this
	newCursor := h.cursor
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {

		case "esc":
			return NewMainMenuHandler(context.model.session.User), nil

		case "enter":
			if h.resultsSelection[newCursor] == "Next Test" {
				if h.testType == "timer" {
					return NewTimerTestHandler(h.mainMenu), nil
				} else if h.testType == "wordcount" {
					return NewWordCountTestHandler(h.mainMenu), nil
				}
			} else if h.resultsSelection[newCursor] == "Main Menu" {
				return NewMainMenuHandler(context.model.session.User), nil
			} else if h.resultsSelection[newCursor] == "Replay" {
				return NewReplayHandler(*h), nil
			}

		case "left", "h":
			if h.cursor == 0 {
				newCursor = len(h.resultsSelection) - 1
			} else {
				newCursor--
			}

		case "right", "l":
			if h.cursor == len(h.resultsSelection)-1 {
				newCursor = 0
			} else {
				newCursor++
			}
		}
	}
	h.cursor = newCursor
	return h, nil
}

func (h *ResultsHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2

	title := style("Test Results", m.styles.magenta)
	title = lipgloss.NewStyle().PaddingBottom(1).Render(title)

	var content []string
	content = append(content, fmt.Sprintf("WPM: %d", h.wpm))
	content = append(content, fmt.Sprintf("Accuracy: %.1f%%", h.accuracy))

	// if h.results.testType == "timer" {
	// 	content = append(content, fmt.Sprintf("Time: %s", formatDuration(h.results.duration)))
	// 	content = append(content, fmt.Sprintf("Words: %d", h.results.wordsTyped))
	// } else if h.results.testType == "wordcount" {
	// 	content = append(content, fmt.Sprintf("Words: %d/%d", h.results.wordsTyped, h.results.targetWords))
	// 	content = append(content, fmt.Sprintf("Time: %s", formatDuration(h.results.duration)))
	// }

	var menuItems []string
	menuItemsStyle := lipgloss.NewStyle().Padding(0, 2, 0, 2)
	for i, choice := range h.resultsSelection {
		choiceShow := style(choice, m.styles.toEnter)
		choiceShow = wrapWithCursor(h.cursor == i, choiceShow, m.styles.correct)
		choiceShow = menuItemsStyle.Render(choiceShow)
		menuItems = append(menuItems, choiceShow)
	}

	resultsMenu := lipgloss.JoinHorizontal(lipgloss.Center, menuItems...)

	fullParagraph := lipgloss.JoinVertical(
		lipgloss.Center, resultsStyle.Padding(0).Render(title),
		resultsStyle.Padding(0).Render(content...),
		resultsStyle.Render(resultsMenu),
	)
	s := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, fullParagraph)

	return s
}

func (h *ResultsHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateResults]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}

func formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%.2fs", d.Seconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

func (base TestBase) calculateRawWpm(elapsedMinutes float64) float64 {
	return base.calculateWpm(len(strings.Split(string(base.inputBuffer), " ")), elapsedMinutes)
}

func (base TestBase) calculateWpm(wordCount int, elapsedMinutes float64) float64 {
	if elapsedMinutes == 0 {
		return 0
	} else {
		grossWpm := float64(wordCount) / elapsedMinutes
		netWpm := grossWpm - float64(len(base.mistakes.mistakesAt))/elapsedMinutes

		return math.Max(0, netWpm)
	}
}

func (base TestBase) calculateNormalizedWpm(elapsedMinutes float64) float64 {
	return base.calculateWpm(len(base.inputBuffer)/5, elapsedMinutes)
}

func (base TestBase) calculateCpm(elapsedMinutes float64) int {
	return int(float64(base.rawInputCount) / elapsedMinutes)
}

func (base TestBase) calculateAccuracy() float64 {
	mistakesRate := float64(base.mistakes.rawMistakesCnt*100) / float64(base.rawInputCount)
	accuracy := 100 - mistakesRate
	return accuracy
}
