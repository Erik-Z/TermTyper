package cmd

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/indent"
	"github.com/muesli/reflow/wordwrap"
)

var lineLenLimit int
var minLineLen int = 5
var maxLineLen int = 40
var resultsStyle = lipgloss.NewStyle().
	Align(lipgloss.Center).
	PaddingTop(1).
	PaddingBottom(1).
	PaddingLeft(5).
	PaddingRight(5)

func (m model) View() string {
	var s string

	termWidth, termHeight := m.width-2, m.height-2
	reactiveLimit := (termWidth * 6) / 10
	lineLenLimit = int(math.Min(float64(maxLineLen), math.Max(float64(minLineLen), float64(reactiveLimit))))

	switch state := m.state.(type) {
	case MainMenu:
		termtyper := style("TermTyper", m.styles.magenta)
		termtyper = lipgloss.NewStyle().PaddingBottom(1).Render(termtyper)
		var menuItems []string
		menuItemsStyle := lipgloss.NewStyle().PaddingTop(1)

		for i, choice := range state.MainMenuSelection {
			choiceShow := style(choice, m.styles.toEnter)

			choiceShow = wrapWithCursor(state.cursor == i, choiceShow, m.styles.toEnter)
			choiceShow = menuItemsStyle.Render(choiceShow)
			menuItems = append(menuItems, choiceShow)
		}

		joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{termtyper}, menuItems...)...)
		s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
		centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#FF00FF"))

		return borderStyle.Render(centeredText)

	case TimerTest:
		var timer string

		timer = style(state.timer.timer.View(), m.styles.magenta)

		paragraph := state.base.renderParagraph(lineLenLimit, m.styles)
		lines := strings.Split(paragraph, "\n")
		cursorLine := findCursorLine(lines, state.base.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		s += positionVertically(termHeight)
		avgLineLen := averageLineLen(lines)
		indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

		s += m.indent(timer, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

		if !state.timer.isRunning {
			s += "\n\n\n"
			s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
		}

		return s + "\n" + string(state.base.inputBuffer)

	case ZenMode:
		stopwatch := style(state.stopwatch.stopwatch.View(), m.styles.magenta)
		paragraphView := state.base.renderParagraphZenMode(lineLenLimit, m.styles)
		lines := strings.Split(paragraphView, "\n")
		cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), state.base.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		s += positionVertically(termHeight)
		//avgLineLen := averageLineLen(lines)
		//indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

		//s += m.indent(stopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)
		s += stopwatch + "\n\n" + linesAroundCursor
		if !state.stopwatch.isRunning {
			s += "\n\n\n"
			s += style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter)
			//s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
		}
		centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center+lipgloss.Position(termHeight/2), s)
		borderStyle := lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#FF00FF"))

		return borderStyle.Render(centeredText)

	case WordCountTest:
		stopwatch := style(state.stopwatch.stopwatch.View(), m.styles.magenta)
		paragraphView := state.base.renderParagraph(lineLenLimit, m.styles)
		lines := strings.Split(paragraphView, "\n")
		cursorLine := findCursorLine(strings.Split(paragraphView, "\n"), state.base.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		s += positionVertically(termHeight)
		avgLineLen := averageLineLen(lines)
		indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

		s += m.indent(stopwatch, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)
		s += "\n\n\n"
		s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))

	case WordCountTestResults:
		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.magenta)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.magenta)
		givenTime := "time: " + style(state.results.time.String(), m.styles.magenta)
		wordCount := "count: " + style(strconv.Itoa(state.wordCount), m.styles.magenta)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.magenta)

		statsLine1 := fmt.Sprintf("%s %s %s", accuracy, rawWpmShow, givenTime)
		statsLine2 := wordCount

		var menuItems []string
		menuItemsStyle := lipgloss.NewStyle().Padding(0, 2, 0, 2)
		for i, choice := range state.results.resultsSelection {
			choiceShow := style(choice, m.styles.toEnter)

			choiceShow = wrapWithCursor(state.results.cursor == i, choiceShow, m.styles.toEnter)
			choiceShow = menuItemsStyle.Render(choiceShow)
			menuItems = append(menuItems, choiceShow)
		}

		resultsMenu := lipgloss.JoinHorizontal(lipgloss.Center, menuItems...)

		fullParagraph := lipgloss.JoinVertical(
			lipgloss.Center, resultsStyle.Padding(1).Render(wpm),
			resultsStyle.Padding(0).Render(statsLine1),
			resultsStyle.Render(statsLine2), resultsStyle.Render(resultsMenu),
		)
		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, fullParagraph)

	case TimerTestResult:
		rawWpmShow := "raw: " + style(strconv.Itoa(state.results.rawWpm), m.styles.magenta)
		wpm := "wpm: " + style(strconv.Itoa(state.results.wpm), m.styles.magenta)
		givenTime := "time: " + style(state.results.time.String(), m.styles.magenta)
		accuracy := "accuracy: " + style(fmt.Sprintf("%.1f", state.results.accuracy), m.styles.magenta)

		statsLine1 := fmt.Sprintf("%s %s %s", accuracy, rawWpmShow, givenTime)

		fullParagraph := lipgloss.JoinVertical(lipgloss.Center, resultsStyle.Padding(1).Render(wpm), resultsStyle.Padding(0).Render(statsLine1))
		s = lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, fullParagraph)
	}

	return s
}

func (base *TestBase) renderParagraph(lineLimit int, styles Styles) string {
	paragraph := base.renderInput(styles)
	paragraph += base.renderCursor(styles)
	paragraph += base.renderWordsToEnter(styles)

	wrappedParagraph := wrapParagraph(paragraph, lineLimit)
	return wrappedParagraph
}

func (base *TestBase) renderParagraphZenMode(lineLimit int, styles Styles) string {
	paragraph := base.renderInputZenMode(styles)
	paragraph += base.renderCursorZenMode(styles)

	wrappedParagraph := wrapParagraph(paragraph, lineLimit)
	return wrappedParagraph
}

func (base *TestBase) renderCursorZenMode(styles Styles) string {
	cursorLetter := [1]rune{' '}

	return style(string(cursorLetter[:]), styles.cursor)
}

func (base *TestBase) renderInputZenMode(styles Styles) string {
	var input strings.Builder
	input.WriteString(styleAll(base.inputBuffer, styles.correct))
	return input.String()
}

func (base *TestBase) renderInput(styles Styles) string {
	mistakes := mapToKeysSlice(base.mistakes.mistakesAt)
	sort.Ints(mistakes)

	var input strings.Builder

	if len(mistakes) == 0 {
		input.WriteString(styleAll(base.inputBuffer, styles.correct))
	} else {
		previousMistake := -1

		for _, mistakeAt := range mistakes {
			sliceUntilMistake := base.inputBuffer[previousMistake+1 : mistakeAt]
			mistakeSlice := base.wordsToEnter[mistakeAt : mistakeAt+1]

			input.WriteString(styleAll(sliceUntilMistake, styles.correct))
			input.WriteString(style(string(mistakeSlice), styles.mistake))

			previousMistake = mistakeAt
		}

		inputAfterLastMistake := base.inputBuffer[previousMistake+1:]
		input.WriteString(styleAll(inputAfterLastMistake, styles.correct))
	}

	return input.String()
}

func (base *TestBase) renderCursor(styles Styles) string {
	if len(base.inputBuffer) == len(base.wordsToEnter) {
		s := [1]rune{' '}
		return style(string(s[:]), styles.cursor)
	}
	cursorLetter := base.wordsToEnter[len(base.inputBuffer) : len(base.inputBuffer)+1]

	return style(string(cursorLetter), styles.cursor)
}

func (base *TestBase) renderWordsToEnter(styles Styles) string {
	if len(base.inputBuffer) == len(base.wordsToEnter) {
		return ""
	}
	wordsToEnter := base.wordsToEnter[len(base.inputBuffer)+1:]

	return style(string(wordsToEnter), styles.toEnter)
}

func positionVertically(termHeight int) string {
	var acc strings.Builder

	for i := 0; i < termHeight/2-3; i++ {
		acc.WriteRune('\n')
	}

	return acc.String()
}

func (m model) indent(block string, indentBy uint) string {
	indentedBlock := indent.String(block, indentBy)

	return indentedBlock
}

func wrapParagraph(paragraph string, lineLimit int) string {
	paragraph = strings.ReplaceAll(paragraph, " ", "💩")

	f := wordwrap.NewWriter(lineLimit)
	f.Breakpoints = []rune{'💩'}
	f.KeepNewlines = false
	f.Write([]byte(paragraph))
	f.Close()

	paragraph = strings.ReplaceAll(f.String(), "💩", " ")

	return paragraph
}

func wrapWithCursor(shouldWrap bool, line string, stringStyle StringStyle) string {
	cursor := " "
	cursorClose := " "
	if shouldWrap {
		cursor = style(">", stringStyle)
		cursorClose = style("<", stringStyle)
	}
	return fmt.Sprintf("%s %s %s", cursor, line, cursorClose)
}

func findCursorLine(lines []string, cursorAt int) int {
	lenAcc := 0
	cursorLine := 0

	for _, line := range lines {
		lineLen := len(dropAnsiCodes(line))

		lenAcc += lineLen

		if cursorAt <= lenAcc-1 {
			return cursorLine
		} else {
			cursorLine += 1
		}
	}

	return cursorLine
}

func getLinesAroundCursor(lines []string, cursorLine int) []string {
	cursor := cursorLine

	if cursorLine == 0 {
		cursor += 3
	} else {
		cursor += 2
	}

	low := int(math.Max(0, float64(cursorLine-1)))
	high := int(math.Min(float64(len(lines)), float64(cursor)))

	return lines[low:high]
}

func style(str string, style StringStyle) string {
	return style(str).String()
}

func styleAll(runes []rune, style StringStyle) string {
	var acc strings.Builder

	for idx, char := range runes {
		_ = idx
		acc.WriteString(style(string(char)).String())
	}

	return acc.String()
}

func dropAnsiCodes(colored string) string {
	m := regexp.MustCompile("\x1b\\[[0-9;]*m")

	return m.ReplaceAllString(colored, "")
}
