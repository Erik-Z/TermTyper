package cmd

import (
	"termtyper/database"
	"termtyper/words"
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

type MainMenu struct {
	MainMenuSelection      []string
	cursor                 int
	timerTestWordGenerator words.WordGenerator
	wordTestWordGenerator  words.WordGenerator
	currentUser            *database.ApplicationUser
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

type ZenMode struct {
	base      TestBase
	stopwatch StopWatch
}

type WordCountTest struct {
	stopwatch StopWatch
	base      TestBase
	completed bool
}

type Replay struct {
	test              TestBase
	results           *WordCountTestResults
	stopwatch         StopWatch
	isReplayInProcess bool
}

type Results struct {
	testType         string
	wpm              int
	accuracy         float64
	deltaWpm         float64
	rawWpm           int
	cpm              int
	time             time.Duration
	wordList         string
	test             TestBase
	wpmEachSecond    []float64
	mainMenu         MainMenu
	resultsSelection []string
	cursor           int
}

type WordCountTestResults struct {
	wpmEachSecond []float64
	wordCount     int
	results       Results
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
