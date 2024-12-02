package cmd

import (
	"termtyper/words"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	"github.com/muesli/termenv"
)

type model struct {
	state         State
	styles        Styles
	width, height int
}

type State interface{}

type MainMenu struct {
	MainMenuSelection      []string
	cursor                 int
	timerTestWordGenerator words.WordGenerator
	wordTestWordGenerator  words.WordGenerator
}

type TestBase struct {
	wordsToEnter  []rune
	inputBuffer   []rune
	wpmEachSecond []float64
	rawInputCount int
	mistakes      mistakes
	cursor        int
	testRecord    []KeyPress
	mainMenu      MainMenu
}

type KeyPress struct {
	key       string
	timestamp time.Time
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

type TimerTestResult struct {
	wpmEachSecond []float64
	results       Results
}

type Timer struct {
	timer     timer.Model
	duration  time.Duration
	timedout  bool
	isRunning bool
}

type Results struct {
	wpm           int
	accuracy      float64
	deltaWpm      float64
	rawWpm        int
	cpm           int
	time          time.Duration
	wordList      string
	wpmEachSecond []float64
}

type StringStyle func(string) termenv.Style

type Styles struct {
	correct StringStyle
	mistake StringStyle
	cursor  StringStyle
	toEnter StringStyle
	magenta StringStyle
}
