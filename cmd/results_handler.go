package cmd

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ResultsHandler struct {
	*BaseStateHandler
	testType         string
	wpm              int
	accuracy         float64
	deltaWpm         float64
	rawWpm           int
	cpm              int
	time             time.Duration
	wordList         string
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
				//return wordCountTestResults.showReplay()
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
