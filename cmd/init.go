package cmd

import (
	"runtime"
	"termtyper/words"
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/muesli/termenv"
)

func (m model) Init() tea.Cmd {
	return nil
}

func initMainMenu() MainMenu {
	return MainMenu{
		MainMenuSelection: []string{
			"Timer",
			"Word Count",
			"Zen",
			"Config",
		},
		cursor:                 0,
		timerTestWordGenerator: words.NewGenerator(),
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
	}
}

func initModel(termProfile termenv.Profile, foregroundColor termenv.Color, width, height int) model {
	return model{
		state:  initMainMenu(),
		width:  width,
		height: height,
		styles: Styles{
			correct: func(str string) termenv.Style {
				return termenv.String(str).Foreground(foregroundColor)
			},
			toEnter: func(str string) termenv.Style {
				return termenv.String(str).Foreground(foregroundColor).Faint()
			},
			mistake: func(str string) termenv.Style {
				return termenv.String(str).Foreground(termProfile.Color("1")).Underline()
			},
			cursor: func(str string) termenv.Style {
				return termenv.String(str).Reverse().Bold()
			},
			magenta: func(str string) termenv.Style {
				return termenv.String(str).Foreground(termProfile.Color("#FF00FF"))
			},
		},
	}
}

func OsInit() {
	// enable colors for one guy who uses windows
	if runtime.GOOS == "windows" {
		mode, err := termenv.EnableWindowsANSIConsole()
		if err != nil {
			panic(err)
		}
		defer termenv.RestoreWindowsConsole(mode)
	}
}
