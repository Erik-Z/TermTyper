package main

import (
	"log"
	"termtyper/cmd"
)

// var (
// 	baseStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
// 	correctStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
// 	wrongStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
// 	promptStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
// )

// type model struct {
// 	prompt     string
// 	input      string
// 	startTime  time.Time
// 	endTime    time.Time
// 	typingDone bool
// }

// func initialModel() model {
// 	return model{
// 		prompt: "The quick brown fox jumps over the lazy dog.",
// 	}
// }

// func (m model) Init() tea.Cmd {
// 	return nil
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		if m.typingDone {
// 			switch msg.String() {
// 			case "q":
// 				os.Exit(0)
// 			}
// 			return m, nil
// 		}

// 		switch msg.String() {
// 		case "enter":
// 			if m.input == m.prompt {
// 				m.endTime = time.Now()
// 				m.typingDone = true
// 			}
// 		case "backspace":
// 			if len(m.input) > 0 {
// 				m.input = m.input[:len(m.input)-1]
// 			}
// 		default:
// 			m.input += msg.String()
// 		}

// 		if m.startTime.IsZero() {
// 			m.startTime = time.Now()
// 		}
// 	}

// 	return m, nil
// }

// func (m model) View() string {
// 	var result strings.Builder

// 	result.WriteString(promptStyle.Render("Typing Test") + "\n\n")

// 	if m.typingDone {
// 		duration := m.endTime.Sub(m.startTime).Seconds()
// 		chars := len(m.prompt)
// 		wpm := int(float64((chars / 5)) / (duration / 60))
// 		result.WriteString(fmt.Sprintf("You completed the test! \n\nTime: %.2f seconds\nWPM: %d\n\nPress q to quit.",
// 			duration, wpm,
// 		))

// 		return result.String()
// 	}

// 	for i, char := range m.prompt {
// 		if i < len(m.input) {
// 			if m.input[i] == byte(char) {
// 				result.WriteString(correctStyle.Render(string(char)))
// 			} else {
// 				if byte(char) == ' ' {
// 					result.WriteString(wrongStyle.Render(string(m.input[i])))
// 				} else {
// 					result.WriteString(wrongStyle.Render(string(char)))
// 				}
// 			}
// 		} else {
// 			result.WriteString(baseStyle.Render(string(char)))
// 		}
// 	}

// 	result.WriteString("\n\n" + lipgloss.NewStyle().Italic(true).Render("Start typing above!") + "\n")
// 	result.WriteString(baseStyle.Render("Press Enter when done, Esc to quit."))

// 	return result.String()
// }

func main() {
	// p := tea.NewProgram(initialModel())
	// if _, err := p.Run(); err != nil {
	// 	fmt.Printf("Error starting app: %v\n", err)
	// 	os.Exit(1)
	// }

	cmd.OsInit()
	if err := cmd.RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
