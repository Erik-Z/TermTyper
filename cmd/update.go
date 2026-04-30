package cmd

import (
	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.session.mu.Lock()
	defer m.session.mu.Unlock()
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		if msg.Width == 0 && msg.Height == 0 {
			return m, nil
		} else {
			m.stateMachine.model.width = msg.Width
			m.stateMachine.model.height = msg.Height
			return m, nil
		}
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	return m.stateMachine.HandleInput(msg)
}

func handleBackspace(base *TestBase) {
	if len(base.mistakes.mistakesAt) == 0 && len(base.wordsToEnter) > 0 {
		return
	}

	base.inputBuffer = deleteLastChar(base.inputBuffer)
	inputLength := len(base.inputBuffer)
	base.rawInputCount -= 1
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

func handleCharacterInputFromMsg(msg tea.KeyPressMsg, base *TestBase) {
	// TODO: Fix this. Doesn't work, I can still paste into input buffer.
	if len(msg.Text) > 1 {
		return
	}

	if len(base.inputBuffer) == len(base.wordsToEnter) {
		return
	}
	inputLetter := []rune(msg.Text)[len([]rune(msg.Text))-1]
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

func handleCharacterInputZenMode(msg tea.KeyPressMsg, base *TestBase) {
	inputLetter := []rune(msg.Text)[len([]rune(msg.Text))-1]
	base.inputBuffer = append(base.inputBuffer, inputLetter)
	base.rawInputCount += 1

	newCursorPosition := len(base.inputBuffer)
	base.cursor = newCursorPosition
}

func recordInput(msg tea.KeyPressMsg, base *TestBase, timestamp int64) {
	var keyPress KeyPress
	if msg.String() == "backspace" {
		keyPress = KeyPress{
			key:       '\b',
			timestamp: timestamp,
		}
	} else {
		keyPress = KeyPress{
			key:       []rune(msg.Text)[len([]rune(msg.Text))-1],
			timestamp: timestamp,
		}
	}

	base.testRecord = append(base.testRecord, keyPress)
}

func recordInputBackspace(base *TestBase, timestamp int64) {
	keyPress := KeyPress{
		key:       '\b',
		timestamp: timestamp,
	}
	base.testRecord = append(base.testRecord, keyPress)
}
