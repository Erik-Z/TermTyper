package cmd

import (
	"fmt"
	"termtyper/database"
)

type Settings struct {
	settingsCursor    int
	settingSelections []TestSetting
}

type TestSetting interface {
	render(styles Styles) string
}

type TimerSettings struct {
	timerCursor     int
	timerSelection  []int
	selectionCursor int
}

type WordsSettings struct {
	wordsCursor     int
	wordsSelection  []int
	selectionCursor int
}

func (t TimerSettings) render(style Styles) string {
	selections := []string{formatSettingsDuration(t.timerSelection[t.timerCursor])}
	selectionsStr := showSelections(selections, t.selectionCursor, style)
	return fmt.Sprintf("%s %s", "Timer", selectionsStr)
}

func (w WordsSettings) render(style Styles) string {
	selections := []string{fmt.Sprintf("%d", w.wordsSelection[w.wordsCursor])}
	selectionsStr := showSelections(selections, w.selectionCursor, style)
	return fmt.Sprintf("%s %s", "Words", selectionsStr)
}

func findTimerIndex(config *database.UserConfig, timerSelection []int) int {
	for i, num := range timerSelection {
		if num == config.Time {
			return i
		}
	}
	return 0
}

func showSelections(selections []string, cursor int, styles Styles) string {
	var selectionsStr string
	for i, option := range selections {
		if i+1 == cursor {
			selectionsStr += "[" + style(option, styles.magenta) + "]"
		} else {
			selectionsStr += "[" + style(option, styles.toEnter) + "]"
		}
		selectionsStr += " "
	}
	return selectionsStr
}

func formatSettingsDuration(seconds int) string {
	minutes := seconds / 60
	remainingSeconds := seconds % 60
	return fmt.Sprintf("%dm%ds", minutes, remainingSeconds)
}
