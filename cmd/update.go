package cmd

import (
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var commands []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width == 0 && msg.Height == 0 {
			return m, nil
		} else {
			m.width = msg.Width
			m.height = msg.Height
			return m, nil
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

	}
	switch state := m.state.(type) {
	case MainMenu:
		m.state = state.handleInput(msg)

	case TimerTest:
		switch msg := msg.(type) {
		case timer.TickMsg:
			timerUpdate, cmdUpdate := state.timer.timer.Update(msg)
			state.timer.timer = timerUpdate
			commands = append(commands, cmdUpdate)

			elapsedMinute := state.timer.duration.Minutes() - state.timer.timer.Timeout.Minutes()
			if elapsedMinute != 0 {
				state.base.wpmEachSecond = append(state.base.wpmEachSecond, state.base.calculateNormalizedWpm(elapsedMinute))
			}

			m.state = state

			if state.timer.timer.Timedout() {
				state.timer.timedout = true

				// var results = state.calculateResults()
				// m.state = TimerTestResult{
				// 	wpmEachSecond: state.base.wpmEachSecond,
				// 	results:       results,
				// }

				m.state = initMainMenu()
			}

		case tea.KeyMsg:
			switch msg.String() {

			case "backspace":
				m.state = state
			case "ctrl+w":
				m.state = state.base.mainMenu
				return m, nil
			case "ctrl+r":
				// Restart the test
				m.state = initTimerTest(state.base.mainMenu)
				return m, nil

			default:
				switch msg.Type {
				case tea.KeyRunes:
					if !state.timer.isRunning {
						commands = append(commands, state.timer.timer.Init())
						state.timer.isRunning = true
					}

					handleCharacterInput(msg, &state.base)
					m.state = state
				}
			}
		}

	}

	return m, tea.Batch(commands...)
}

func (menu MainMenu) handleInput(msg tea.Msg) State {
	newCursor := menu.cursor
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return initTimerTest(menu)
		case "up", "k":
			if menu.cursor == 0 {
				newCursor = len(menu.MainMenuSelection) - 1
			} else {
				newCursor--
			}

		case "down", "j":
			if menu.cursor == len(menu.MainMenuSelection)-1 {
				newCursor = 0
			} else {
				newCursor++
			}

		}

	}
	menu.cursor = newCursor
	return menu
}

func handleBackspace(base *TestBase) {

}

func handleCharacterInput(msg tea.KeyMsg, base *TestBase) {
	inputLetter := msg.Runes[len(msg.Runes)-1]
	currInputBufferLen := len(base.inputBuffer)
	correctNextLetter := base.wordsToEnter[currInputBufferLen]

	base.inputBuffer = append(base.inputBuffer, inputLetter)
	base.rawInputCount += 1

	if inputLetter != correctNextLetter {
		base.mistakes.mistakesAt[currInputBufferLen] = true
		base.mistakes.rawMistakesCnt = base.mistakes.rawMistakesCnt + 1
	}

	newCursorPosition := len(base.inputBuffer)
	base.cursor = newCursorPosition
}
