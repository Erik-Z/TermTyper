package cmd

import (
	"fmt"
	"math"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

var lineLenLimit int
var minLineLen int = 5
var maxLineLen int = 40

func (m model) View() string {
	var result strings.Builder

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
		var s string
		var timer string

		timer = style(state.timer.timer.View(), m.styles.magenta)

		paragraph := state.base.renderParagraph(lineLenLimit, m.styles)
		lines := strings.Split(paragraph, "\n")
		cursorLine := findCursorLine(lines, state.base.cursor)

		linesAroundCursor := strings.Join(getLinesAroundCursor(lines, cursorLine), "\n")

		s += positionVerticaly(termHeight)
		avgLineLen := averageLineLen(lines)
		indentBy := uint(math.Max(0, float64(termWidth/2-avgLineLen/2)))

		s += m.indent(timer, indentBy) + "\n\n" + m.indent(linesAroundCursor, indentBy)

		if !state.timer.isRunning {
			s += "\n\n\n"
			s += lipgloss.PlaceHorizontal(termWidth, lipgloss.Center, style("ctrl+r to restart, ctrl+q to menu", m.styles.toEnter))
		}

		return s + "\n" + string(state.base.inputBuffer)
	}

	return result.String()
}

func (base *TestBase) renderParagraph(lineLimit int, styles Styles) string {
	paragraph := base.renderInput(styles)
	paragraph += base.renderCursor(styles)
	paragraph += base.renderWordsToEnter(styles)

	wrappedParagraph := wrapParagraph(paragraph, lineLimit)
	return wrappedParagraph
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
	cursorLetter := base.wordsToEnter[len(base.inputBuffer) : len(base.inputBuffer)+1]

	return style(string(cursorLetter), styles.cursor)
}

func (base *TestBase) renderWordsToEnter(styles Styles) string {
	wordsToEnter := base.wordsToEnter[len(base.inputBuffer)+1:] // without cursor

	return style(string(wordsToEnter), styles.toEnter)
}

func positionVerticaly(termHeight int) string {
	var acc strings.Builder

	for i := 0; i < termHeight/2-3; i++ {
		acc.WriteRune('\n')
	}

	return acc.String()
}

func (m model) indent(block string, indentBy uint) string {
	indentation := strings.Repeat(" ", int(indentBy))
	lines := strings.Split(block, "\n")

	for i, line := range lines {
		lines[i] = indentation + line
	}

	return strings.Join(lines, "\n")
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
