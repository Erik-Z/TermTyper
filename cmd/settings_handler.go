package cmd

import (
	"fmt"
	"termtyper/database"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

//TODO: add "unsaved settings. Do you want to save before returning to main menu"

type SettingsHandler struct {
	*BaseStateHandler
	settingsCursor    int
	settingSelections []TestSetting
	userConfig        database.UserConfig
}

type TestSetting interface {
	render(styles Styles) string
	MoveLeft()
	MoveRight()
	SaveSettings(context *StateContext)
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

type PunctuationSettings struct {
	enabled    bool
	savedValue bool
}

type ThemeSettings struct {
	themeIndex int
	savedIndex int
}

func NewSettingsHandler(user *database.ApplicationUser) *SettingsHandler {
	wordCountSelection := []int{15, 30, 45, 60}
	timerSelection := []int{15, 30, 60, 120}

	timerSettings := TimerSettings{
		timerSelection:  timerSelection,
		timerCursor:     findTimerIndex(user.Config, timerSelection),
		selectionCursor: findTimerIndex(user.Config, timerSelection),
	}

	wordsSettings := WordsSettings{
		wordsSelection:  wordCountSelection,
		wordsCursor:     findWordsIndex(user.Config, wordCountSelection),
		selectionCursor: findWordsIndex(user.Config, wordCountSelection),
	}

	punctuationSettings := PunctuationSettings{
		enabled:    user.Config.Punctuation,
		savedValue: user.Config.Punctuation,
	}

	themeSettings := ThemeSettings{
		themeIndex: GetThemeIndex(user.Config.Theme),
		savedIndex: GetThemeIndex(user.Config.Theme),
	}

	return &SettingsHandler{
		BaseStateHandler:  NewBaseStateHandler(StateSettings),
		settingsCursor:    0,
		settingSelections: []TestSetting{&timerSettings, &wordsSettings, &punctuationSettings, &themeSettings},
		userConfig:        *user.Config,
	}
}

func (h *SettingsHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	newCursor := h.settingsCursor
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc", "ctrl+q":
			if h.ValidateTransition(StateMainMenu, context) {
				return NewMainMenuHandler(context.model.session.User, context.model), nil
			}

		case "enter":
			h.settingSelections[h.settingsCursor].SaveSettings(context)

		case "up", "k":
			if h.settingsCursor == 0 {
				newCursor = len(h.settingSelections) - 1
			} else {
				newCursor--
			}

		case "down", "j":
			if h.settingsCursor == len(h.settingSelections)-1 {
				newCursor = 0
			} else {
				newCursor++
			}

		case "left", "h":
			h.settingSelections[h.settingsCursor].MoveLeft()

		case "right", "l":
			h.settingSelections[h.settingsCursor].MoveRight()

		}

	}
	h.settingsCursor = newCursor
	return h, nil
}

func (h *SettingsHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	settings := style("Settings", m.styles.themeFunc)
	settings = lipgloss.NewStyle().PaddingBottom(1).Render(settings)

	menuItemsStyle := lipgloss.NewStyle().PaddingTop(1)
	var settingSelection []string
	for i, choice := range h.settingSelections {
		choiceShow := choice.render(m.styles)

		choiceShow = wrapWithCursor(h.settingsCursor == i, choiceShow, m.styles.toEnter)
		choiceShow = menuItemsStyle.Render(choiceShow)
		settingSelection = append(settingSelection, choiceShow)
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{settings}, settingSelection...)...)
	renderString := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, renderString)

	return centeredText
}

func (t *TimerSettings) render(styles Styles) string {
	var renderColor StringStyle
	formattedDuration := formatSettingsDuration(t.timerSelection[t.selectionCursor])
	if t.selectionCursor == t.timerCursor {
		renderColor = styles.themeFunc
	} else {
		renderColor = styles.toEnter
	}
	selectionsStr := "[" + style(formattedDuration, renderColor) + "]"
	return fmt.Sprintf("%s %s", "Timer", selectionsStr)
}

func (w *WordsSettings) render(styles Styles) string {
	var renderColor StringStyle
	numberOfWords := fmt.Sprintf("%d", w.wordsSelection[w.selectionCursor])
	if w.selectionCursor == w.wordsCursor {
		renderColor = styles.themeFunc
	} else {
		renderColor = styles.toEnter
	}
	selectionsStr := "[" + style(numberOfWords, renderColor) + "]"
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

func findWordsIndex(config *database.UserConfig, wordCountSelection []int) int {
	for i, num := range wordCountSelection {
		if num == config.Words {
			return i
		}
	}
	return 0
}

func (t *TimerSettings) MoveLeft() {
	if t.selectionCursor == 0 {
		t.selectionCursor = len(t.timerSelection) - 1
	} else {
		t.selectionCursor--
	}
}

func (t *TimerSettings) MoveRight() {
	if t.selectionCursor == len(t.timerSelection)-1 {
		t.selectionCursor = 0
	} else {
		t.selectionCursor++
	}
}

func (t *TimerSettings) SaveSettings(context *StateContext) {
	if t.selectionCursor != t.timerCursor {
		t.timerCursor = t.selectionCursor
	}
	newUserConfig := context.model.session.User.Config
	newUserConfig.Time = t.timerSelection[t.timerCursor]

	database.UpdateUserConfigStandalone(
		context.model.context.UserRepository,
		context.model.session.User.Id,
		UserConfigToMap(newUserConfig))
}

func (w *WordsSettings) MoveLeft() {
	if w.selectionCursor == 0 {
		w.selectionCursor = len(w.wordsSelection) - 1
	} else {
		w.selectionCursor--
	}
}

func (w *WordsSettings) MoveRight() {
	if w.selectionCursor == len(w.wordsSelection)-1 {
		w.selectionCursor = 0
	} else {
		w.selectionCursor++
	}
}

func (w *WordsSettings) SaveSettings(context *StateContext) {
	if w.selectionCursor != w.wordsCursor {
		w.wordsCursor = w.selectionCursor
	}
	newUserConfig := context.model.session.User.Config
	newUserConfig.Words = w.wordsSelection[w.wordsCursor]

	database.UpdateUserConfigStandalone(
		context.model.context.UserRepository,
		context.model.session.User.Id,
		UserConfigToMap(newUserConfig))
}

func (p *PunctuationSettings) render(styles Styles) string {
	var renderColor StringStyle
	status := p.enabled
	if p.savedValue == p.enabled {
		renderColor = styles.themeFunc
	} else {
		renderColor = styles.toEnter
	}
	selectionsStr := "[" + style(fmt.Sprintf("%t", status), renderColor) + "]"
	return fmt.Sprintf("%s %s", "Punctuation", selectionsStr)
}

func (p *PunctuationSettings) MoveLeft() {
	p.enabled = false
}

func (p *PunctuationSettings) MoveRight() {
	p.enabled = true
}

func (p *PunctuationSettings) SaveSettings(context *StateContext) {
	p.savedValue = p.enabled
	newUserConfig := context.model.session.User.Config
	newUserConfig.Punctuation = p.enabled

	database.UpdateUserConfigStandalone(
		context.model.context.UserRepository,
		context.model.session.User.Id,
		UserConfigToMap(newUserConfig))
}

func (t *ThemeSettings) render(styles Styles) string {
	var renderColor StringStyle
	if t.themeIndex == t.savedIndex {
		renderColor = styles.themeFunc
	} else {
		renderColor = styles.toEnter
	}
	themeName := AvailableThemes[t.themeIndex].Name
	selectionsStr := "[" + style(themeName, renderColor) + "]"
	return fmt.Sprintf("%s %s", "Theme", selectionsStr)
}

func (t *ThemeSettings) MoveLeft() {
	if t.themeIndex == 0 {
		t.themeIndex = len(AvailableThemes) - 1
	} else {
		t.themeIndex--
	}
}

func (t *ThemeSettings) MoveRight() {
	if t.themeIndex == len(AvailableThemes)-1 {
		t.themeIndex = 0
	} else {
		t.themeIndex++
	}
}

func (t *ThemeSettings) SaveSettings(context *StateContext) {
	t.savedIndex = t.themeIndex
	newUserConfig := context.model.session.User.Config
	newUserConfig.Theme = AvailableThemes[t.themeIndex].Name

	database.UpdateUserConfigStandalone(
		context.model.context.UserRepository,
		context.model.session.User.Id,
		UserConfigToMap(newUserConfig))
}

func formatSettingsDuration(seconds int) string {
	minutes := seconds / 60
	remainingSeconds := seconds % 60
	return fmt.Sprintf("%dm%ds", minutes, remainingSeconds)
}
