package cmd

import (
	"fmt"
	"math"
	"regexp"
	"sort"
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
	m.session.mu.Lock()
	defer m.session.mu.Unlock()

	termWidth := m.width - 2
	reactiveLimit := (termWidth * 6) / 10
	lineLenLimit = int(math.Min(float64(maxLineLen), math.Max(float64(minLineLen), float64(reactiveLimit))))

	return m.stateMachine.Render()
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

			var mistakeSlice []rune
			mistakeSlice = base.wordsToEnter[mistakeAt : mistakeAt+1]
			if string(mistakeSlice) == " " {
				mistakeSlice = base.inputBuffer[mistakeAt : mistakeAt+1]
			}

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
