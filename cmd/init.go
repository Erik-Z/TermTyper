package cmd

import (
	"fmt"
	"runtime"
	"termtyper/database"

	"github.com/charmbracelet/bubbles/stopwatch"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func (m model) Init() tea.Cmd {
	return nil
}

// func initTimerTest(menu MainMenu) TimerTest {
// 	return TimerTest{
// 		timer: Timer{
// 			timer:     timer.NewWithInterval(30*time.Second, time.Second),
// 			duration:  30 * time.Second,
// 			isRunning: false,
// 			timedout:  false,
// 		},
// 		base: TestBase{
// 			wordsToEnter:  menu.timerTestWordGenerator.Generate("Common words"),
// 			inputBuffer:   make([]rune, 0),
// 			rawInputCount: 0,
// 			mistakes: mistakes{
// 				mistakesAt:     make(map[int]bool, 0),
// 				rawMistakesCnt: 0,
// 			},
// 			cursor:      0,
// 			mainMenuOld: menu,
// 		},
// 		completed: false,
// 	}
// }

func initZenMode(menu MainMenuHandler) ZenMode {
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
			cursor: 0,
			//mainMenuOld: menu,
		},
	}
}

func initModel(termProfile termenv.Profile, foregroundColor termenv.Color, width, height int, sess *Session) model {
	databaseContext := database.InitDB()
	m := model{
		//state:   initPreAuthentication(&databaseContext),
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
