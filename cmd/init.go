package cmd

import (
	"fmt"
	"runtime"
	"termtyper/database"
	"termtyper/words"
	"time"

	"github.com/charmbracelet/bubbles/stopwatch"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func (m model) Init() tea.Cmd {
	return nil
}

func initMainMenu(user *database.ApplicationUser) MainMenu {
	return MainMenu{
		MainMenuSelection: []string{
			"Timer",
			"Word Count",
			"Zen",
			"Config",
		},
		currentUser:            user,
		cursor:                 0,
		timerTestWordGenerator: words.NewGenerator(),
		wordTestWordGenerator:  words.NewGenerator(),
	}
}

func initTimerTest(menu MainMenu) TimerTest {
	return TimerTest{
		timer: Timer{
			timer:     timer.NewWithInterval(30*time.Second, time.Second),
			duration:  30 * time.Second,
			isRunning: false,
			timedout:  false,
		},
		base: TestBase{
			wordsToEnter:  menu.timerTestWordGenerator.Generate("Common words"),
			inputBuffer:   make([]rune, 0),
			rawInputCount: 0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor:   0,
			mainMenu: menu,
		},
		completed: false,
	}
}

func initWordCountTest(menu MainMenu) WordCountTest {
	menu.wordTestWordGenerator.Count = 30
	return WordCountTest{
		stopwatch: StopWatch{
			stopwatch: stopwatch.NewWithInterval(time.Millisecond * 100),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter:  menu.wordTestWordGenerator.Generate("Common words"),
			inputBuffer:   make([]rune, 0),
			rawInputCount: 0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor:   0,
			mainMenu: menu,
		},
		completed: false,
	}
}

func initZenMode(menu MainMenu) ZenMode {
	return ZenMode{
		stopwatch: StopWatch{
			stopwatch: stopwatch.New(),
			isRunning: false,
		},
		base: TestBase{
			wordsToEnter:  make([]rune, 0),
			inputBuffer:   make([]rune, 0),
			rawInputCount: 0,
			mistakes: mistakes{
				mistakesAt:     make(map[int]bool, 0),
				rawMistakesCnt: 0,
			},
			cursor:   0,
			mainMenu: menu,
		},
	}
}

func initModel(termProfile termenv.Profile, foregroundColor termenv.Color, width, height int, sess *Session) model {
	databaseContext := database.InitDB()
	m := model{
		state:   initPreAuthentication(&databaseContext),
		width:   width,
		height:  height,
		context: databaseContext,
		session: sess,
		styles: Styles{
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(foregroundColor)
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(foregroundColor).Faint()
			},
			mistake: func(str string) termenv.Style {
				return termenv.String(str).Foreground(termProfile.Color("1"))
			},
			cursor: func(str string) termenv.Style {
				return termenv.String(str).Reverse().Bold()
			},
			magenta: func(str string) termenv.Style {
				return termenv.String(str).Foreground(termProfile.Color("#FF00FF"))
			},
		},
	}

	// Initialize state machine with pre-authentication state
	m.stateMachine = NewStateMachine(&m)
	m.stateMachine.SetCurrentState(StatePreAuth)
	m.stateMachine.handlers[StatePreAuth] = NewPreAuthHandler(&databaseContext)

	return m
}

func OsInit() {
	if runtime.GOOS == "windows" {
		mode, err := termenv.EnableWindowsANSIConsole()
		if err != nil {
			fmt.Println(err)
		}
		defer termenv.RestoreWindowsConsole(mode)
	}
}
