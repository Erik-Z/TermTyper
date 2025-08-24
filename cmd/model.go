package cmd

import (
	"termtyper/database"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/muesli/termenv"
)

type model struct {
	state         State
	stateMachine  *StateMachine
	styles        Styles
	width, height int
	context       database.Context
	session       *Session
}

type State interface{}

type TestBase struct {
	wordsToEnter  []rune
	inputBuffer   []rune
	wpmEachSecond []float64
	rawInputCount int
	mistakes      mistakes
	cursor        int
	testRecord    []KeyPress
	mainMenu      MainMenuHandler
}

type KeyPress struct {
	key       rune
	timestamp int64
}

type mistakes struct {
	mistakesAt     map[int]bool
	rawMistakesCnt int
}

type TimerTest struct {
	base      TestBase
	completed bool
	timer     Timer
}

type Timer struct {
	timer     timer.Model
	duration  time.Duration
	timedout  bool
	isRunning bool
}

type StopWatch struct {
	stopwatch stopwatch.Model
	isRunning bool
}

type StringStyle func(string) termenv.Style

type Styles struct {
	correct StringStyle
	mistake StringStyle
	cursor  StringStyle
	toEnter StringStyle
	magenta StringStyle
}
