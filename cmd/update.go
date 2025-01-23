package cmd

import (
	"github.com/charmbracelet/bubbles/stopwatch"
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

				var results = state.calculateResults()
				m.state = TimerTestResult{
					wpmEachSecond: state.base.wpmEachSecond,
					results:       results,
				}
			}

		case tea.KeyMsg:
			switch msg.String() {
			case "backspace":
				handleBackspace(&state.base)
				m.state = state
			case "ctrl+t":
				// Delete entire word
				handleCtrlBackspace(&state.base)
				m.state = state
			case "ctrl+w":
				m.state = state.base.mainMenu
				return m, nil
			case "ctrl+r":
				m.state = initTimerTest(state.base.mainMenu)
				return m, nil
			default:
				switch msg.Type {
				case tea.KeyRunes:
					if !state.timer.isRunning {
						commands = append(commands, state.timer.timer.Init())
						state.timer.isRunning = true
					}

					handleCharacterInputFromMsg(msg, &state.base)
					m.state = state
				}
			}
		}

	case ZenMode:
		switch msg := msg.(type) {
		case stopwatch.StartStopMsg:
			stopwatchUpdate, cmdUpdate := state.stopwatch.stopwatch.Update(msg)
			state.stopwatch.stopwatch = stopwatchUpdate
			commands = append(commands, cmdUpdate)

			m.state = state
		case stopwatch.TickMsg:
			stopwatchUpdate, cmdUpdate := state.stopwatch.stopwatch.Update(msg)
			state.stopwatch.stopwatch = stopwatchUpdate
			commands = append(commands, cmdUpdate)

			elapsedMinutes := state.stopwatch.stopwatch.Elapsed().Minutes()
			if elapsedMinutes != 0 {
				state.base.wpmEachSecond = append(state.base.wpmEachSecond, state.base.calculateNormalizedWpm(elapsedMinutes))
			}

			m.state = state

		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "tab":
			case "ctrl+q":
				m.state = state.base.mainMenu
				return m, nil
			case "ctrl+r":
				m.state = initZenMode(state.base.mainMenu)
				return m, nil
			case "ctrl+backspace":
				handleCtrlBackspace(&state.base)
				m.state = state
			case "backspace":
				handleBackspace(&state.base)
				m.state = state
			default:
				switch msg.Type {
				case tea.KeyRunes:
					if !state.stopwatch.isRunning {
						commands = append(commands, state.stopwatch.stopwatch.Init())
						state.stopwatch.isRunning = true
					}

					handleCharacterInputZenMode(msg, &state.base)
					m.state = state
				}
			}

		}

	case WordCountTest:
		switch msg := msg.(type) {
		case stopwatch.StartStopMsg:

			stopwatchUpdate, cmdUpdate := state.stopwatch.stopwatch.Update(msg)
			state.stopwatch.stopwatch = stopwatchUpdate
			commands = append(commands, cmdUpdate)

			m.state = state

		case stopwatch.TickMsg:

			stopwatchUpdate, cmdUpdate := state.stopwatch.stopwatch.Update(msg)
			state.stopwatch.stopwatch = stopwatchUpdate
			commands = append(commands, cmdUpdate)

			elapsedMinutes := state.stopwatch.stopwatch.Elapsed().Minutes()
			if elapsedMinutes != 0 {
				state.base.wpmEachSecond = append(state.base.wpmEachSecond, state.base.calculateNormalizedWpm(elapsedMinutes))
			}

			m.state = state

		case tea.KeyMsg:
			switch msg.String() {
			case "enter", "tab":
			case "ctrl+q":
				m.state = state.base.mainMenu
				return m, nil
			case "ctrl+r":
				m.state = initWordCountTest(state.base.mainMenu)
				return m, nil
			case "ctrl+backspace":
				handleCtrlBackspace(&state.base)
				m.state = state
			case "backspace":
				handleBackspace(&state.base)
				recordInput(msg, &state)
				m.state = state
			default:
				switch msg.Type {
				case tea.KeyRunes, tea.KeySpace:
					if !state.stopwatch.isRunning {
						commands = append(commands, state.stopwatch.stopwatch.Init())
						state.stopwatch.isRunning = true
					}

					handleCharacterInputFromMsg(msg, &state.base)
					recordInput(msg, &state)

					m.state = state
				}
			}
		}
		if len(state.base.wordsToEnter) == len(state.base.inputBuffer) &&
			!state.base.mistakes.mistakesAt[len(state.base.inputBuffer)-1] {
			//termenv.DefaultOutput().Reset()
			results := state.calculateResults()

			results.mainMenu = state.base.mainMenu
			m.state = WordCountTestResults{
				wpmEachSecond: state.base.wpmEachSecond,
				wordCount:     state.base.mainMenu.wordTestWordGenerator.Count,
				results:       results,
			}
		}

	case WordCountTestResults:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+q":
				m.state = state.results.mainMenu
				return m, nil
			case "ctrl+r":
				m.state = initWordCountTest(state.results.mainMenu)
				return m, nil
			}
		}
		m.state = state.handleInput(msg)

	case Replay:
		switch msg := msg.(type) {
		case stopwatch.StartStopMsg:
			stopwatchUpdate, cmdUpdate := state.stopwatch.stopwatch.Update(msg)
			state.stopwatch.stopwatch = stopwatchUpdate
			commands = append(commands, cmdUpdate)

			m.state = state
		case stopwatch.TickMsg:
			stopwatchUpdate, cmdUpdate := state.stopwatch.stopwatch.Update(msg)
			state.stopwatch.stopwatch = stopwatchUpdate
			commands = append(commands, cmdUpdate)

			m.state = state
		}

		if !state.stopwatch.isRunning && len(state.test.testRecord) > 0 {
			commands = append(commands, state.stopwatch.stopwatch.Init())
			state.stopwatch.isRunning = true
		}

		if len(state.test.testRecord) > 0 {
			currentKeyPress := state.test.testRecord[0]
			if currentKeyPress.timestamp <= state.stopwatch.stopwatch.Elapsed().Milliseconds() {
				switch currentKeyPress.key {
				case '\b':
					handleBackspace(&state.test)
				default:
					handleCharacterInputFromRune(currentKeyPress.key, &state.test)
				}
				state.test.testRecord = state.test.testRecord[1:]
			}
		}

		if len(state.test.wordsToEnter) == len(state.test.inputBuffer) &&
			!state.test.mistakes.mistakesAt[len(state.test.inputBuffer)-1] {

			commands = append(commands, state.stopwatch.stopwatch.Stop())
			state.stopwatch.isRunning = false
		}

		m.state = state
	}

	return m, tea.Batch(commands...)
}

func (menu MainMenu) handleInput(msg tea.Msg) State {
	newCursor := menu.cursor
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if menu.MainMenuSelection[newCursor] == "Timer" {
				return initTimerTest(menu)
			} else if menu.MainMenuSelection[newCursor] == "Zen" {
				return initZenMode(menu)
			} else if menu.MainMenuSelection[newCursor] == "Word Count" {
				return initWordCountTest(menu)
			}

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

func (wordCountTestResults WordCountTestResults) handleInput(msg tea.Msg) State {
	newCursor := wordCountTestResults.results.cursor
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if wordCountTestResults.results.resultsSelection[newCursor] == "Next Test" {
				return initWordCountTest(wordCountTestResults.results.mainMenu)
			} else if wordCountTestResults.results.resultsSelection[newCursor] == "Main Menu" {
				return initMainMenu()
			} else if wordCountTestResults.results.resultsSelection[newCursor] == "Replay" {
				return wordCountTestResults.showReplay()
			}
		case "left", "h":
			if wordCountTestResults.results.cursor == 0 {
				newCursor = len(wordCountTestResults.results.resultsSelection) - 1
			} else {
				newCursor--
			}

		case "right", "l":
			if wordCountTestResults.results.cursor == len(wordCountTestResults.results.resultsSelection)-1 {
				newCursor = 0
			} else {
				newCursor++
			}
		}

	}
	wordCountTestResults.results.cursor = newCursor
	return wordCountTestResults
}

func handleBackspace(base *TestBase) {
	if len(base.mistakes.mistakesAt) == 0 && len(base.wordsToEnter) > 0 {
		return
	}

	base.inputBuffer = deleteLastChar(base.inputBuffer)
	inputLength := len(base.inputBuffer)
	_, ok := base.mistakes.mistakesAt[inputLength]

	if ok {
		delete(base.mistakes.mistakesAt, inputLength)
	}

	base.cursor = inputLength
}

func handleCtrlBackspace(base *TestBase) {
	//TODO: Fix this
	//TODO if multiple punctuation is in a row we delete the punctuations
	if len(base.mistakes.mistakesAt) == 0 && len(base.wordsToEnter) > 0 {
		return
	}

	punctuation := [6]rune{' ', ',', '.', '!', '?', ';'}

	charToDelete := 0 // for some reason an ascii 0 is added to the input buffer when you press ctrl+_
	for i := len(base.inputBuffer) - 1; i >= 0; i-- {
		if !containsChar(punctuation[:], base.inputBuffer[i]) {
			charToDelete += 1
			delete(base.mistakes.mistakesAt, i)
		} else if charToDelete == 1 && containsChar(punctuation[:], base.inputBuffer[i]) {
			// if most recent character is a punctuation, delete that also
			charToDelete += 1
			delete(base.mistakes.mistakesAt, i)
		} else {
			break
		}
	}

	base.inputBuffer = base.inputBuffer[0 : len(base.inputBuffer)-(charToDelete)]
	base.cursor = base.cursor - charToDelete
}

func handleCharacterInputFromMsg(msg tea.KeyMsg, base *TestBase) {
	if len(base.inputBuffer) == len(base.wordsToEnter) {
		return
	}
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

func handleCharacterInputFromRune(char rune, base *TestBase) {
	if len(base.inputBuffer) == len(base.wordsToEnter) {
		return
	}
	currInputBufferLen := len(base.inputBuffer)
	correctNextLetter := base.wordsToEnter[currInputBufferLen]

	base.inputBuffer = append(base.inputBuffer, char)
	base.rawInputCount += 1

	if char != correctNextLetter {
		base.mistakes.mistakesAt[currInputBufferLen] = true
		base.mistakes.rawMistakesCnt = base.mistakes.rawMistakesCnt + 1
	}

	newCursorPosition := len(base.inputBuffer)
	base.cursor = newCursorPosition
}

func handleCharacterInputZenMode(msg tea.KeyMsg, base *TestBase) {
	inputLetter := msg.Runes[len(msg.Runes)-1]
	base.inputBuffer = append(base.inputBuffer, inputLetter)
	base.rawInputCount += 1

	newCursorPosition := len(base.inputBuffer)
	base.cursor = newCursorPosition
}

func recordInput(msg tea.KeyMsg, state *WordCountTest) {
	var keyPress KeyPress
	if msg.String() == "backspace" {
		keyPress = KeyPress{
			key:       '\b',
			timestamp: state.stopwatch.stopwatch.Elapsed().Milliseconds(),
		}
	} else {
		keyPress = KeyPress{
			key:       msg.Runes[len(msg.Runes)-1],
			timestamp: state.stopwatch.stopwatch.Elapsed().Milliseconds(),
		}
	}

	state.base.testRecord = append(state.base.testRecord, keyPress)
}
