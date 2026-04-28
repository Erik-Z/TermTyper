package cmd

import (
	"termtyper/database"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	"github.com/muesli/termenv"
)

type model struct {
	stateMachine    *StateMachine
	styles          Styles
	width, height   int
	context         database.Context
	session         *Session
	termProfile     termenv.Profile
	foregroundColor termenv.Color
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

type Timer struct {
	timer      timer.Model
	duration   time.Duration
	timedout   bool
	isRunning  bool
	startTime  time.Time
	elapsed    time.Duration
}

type StopWatch struct {
	stopwatch  stopwatch.Model
	isRunning bool
	startTime time.Time
}

type StringStyle func(string) termenv.Style

type Styles struct {
	correct  StringStyle
	mistake  StringStyle
	cursor   StringStyle
	toEnter  StringStyle
	themeFunc StringStyle
}

func (t *Timer) Elapsed() time.Duration {
	if t.timedout {
		return t.duration
	}
	if !t.isRunning {
		return 0
	}
	return time.Since(t.startTime)
}

func (sw *StopWatch) Elapsed() time.Duration {
	if !sw.isRunning {
		return 0
	}
	return time.Since(sw.startTime)
}
