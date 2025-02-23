package cmd

import (
	"fmt"
	"reflect"

	"net/mail"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	blurredStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

	focusedButton = focusedStyle.Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))
)

type TextInputState interface {
	updateInputs(msg tea.Msg) []tea.Cmd
}

type Login struct {
	inputs []textinput.Model
	cursor int
}

type Register struct {
	form              *huh.Form
	formData          RegisterFormData
	isFormInitialized bool
}

type RegisterFormData struct {
	email           string
	password        string
	confirmPassword string
}

type PreAuthentication struct {
	authMenu []State
	cursor   int
}

type ForgotPassword struct {
	emailInput textinput.Model
}

type ResetPassword struct {
	inputs []textinput.Model
	cursor int
}

func initLoginScreen() Login {
	inputs := make([]textinput.Model, 2)
	var t textinput.Model
	for i := range inputs {
		t = textinput.New()
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Email"
			t.CharLimit = 64
			t.Focus()
		case 1:
			t.Placeholder = "Password"
			t.CharLimit = 32
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		inputs[i] = t
	}
	return Login{
		inputs: inputs,
		cursor: 0,
	}
}

func initRegisterScreen() Register {
	data := RegisterFormData{}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Email").
				Value(&data.email).
				Validate(func(str string) error {
					_, err := mail.ParseAddress(str)
					if err == nil {
						return nil
					}
					return fmt.Errorf("invalid email address")
				}),
			huh.NewInput().
				Title("password").
				Value(&data.password).
				EchoMode(huh.EchoModePassword).
				Validate(func(str string) error {
					if len(str) < 8 {
						return fmt.Errorf("password must be at least 8 characters")
					}
					return nil
				}),
			huh.NewInput().
				Title("Confirm password").
				Value(&data.confirmPassword).
				EchoMode(huh.EchoModePassword).
				Validate(func(str string) error {
					if str != data.password {
						return fmt.Errorf("password must match")
					}
					return nil
				}),
		),
	)

	return Register{
		form:              form,
		formData:          data,
		isFormInitialized: false,
	}
}

func initPreAuthentication() PreAuthentication {
	return PreAuthentication{
		authMenu: []State{
			initRegisterScreen(),
			initLoginScreen(),
		},
		cursor: 0,
	}
}

func (l Login) renderLoginScreen(m model) string {
	termWidth, termHeight := m.width-2, m.height-2
	register := style("Login", m.styles.magenta)
	register = lipgloss.NewStyle().PaddingBottom(1).Render(register)

	var inputStrings []string

	for _, input := range l.inputs {
		inputStrings = append(inputStrings, input.View())
	}

	inputStrings = append(inputStrings, wrapWithCursor(l.cursor == len(l.inputs), focusedButton, m.styles.toEnter))

	joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{register}, inputStrings...)...)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	return centeredText
}

func (r Register) renderRegisterScreen(m model) string {
	termWidth, termHeight := m.width-2, m.height-2
	register := style("Register", m.styles.magenta)
	register = lipgloss.NewStyle().PaddingBottom(1).Render(register)

	joined := lipgloss.JoinVertical(lipgloss.Center, []string{register, r.form.View()}...)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	return centeredText
}

func (p PreAuthentication) renderPreAuthentication(m model) string {
	termWidth, termHeight := m.width-2, m.height-2
	termtyper := style("TermTyper", m.styles.magenta)
	termtyper = lipgloss.NewStyle().PaddingBottom(1).Render(termtyper)

	var authMenu []string
	menuItemsStyle := lipgloss.NewStyle().PaddingTop(1)

	for i, choice := range p.authMenu {
		choiceShow := style(reflect.TypeOf(choice).Name(), m.styles.toEnter)

		choiceShow = wrapWithCursor(p.cursor == i, choiceShow, m.styles.toEnter)
		choiceShow = menuItemsStyle.Render(choiceShow)
		authMenu = append(authMenu, choiceShow)
	}

	joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{termtyper}, authMenu...)...)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	return centeredText
}

func (state *Login) updateLogin(msg tea.Msg) (State, []tea.Cmd) {
	cmds := make([]tea.Cmd, len(state.inputs))
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			if s == "enter" && state.cursor == len(state.inputs) {
				return initMainMenu(), []tea.Cmd{}
			} else {
				if s == "up" || s == "shift+tab" {
					state.cursor--
				} else {
					state.cursor++
				}

				if state.cursor > len(state.inputs)+1 {
					state.cursor = 0
				} else if state.cursor < 0 {
					state.cursor = len(state.inputs)
				}

				for i := 0; i <= len(state.inputs)-1; i++ {
					state.inputs[i].Blur()
					if i == state.cursor {
						cmds[i] = state.inputs[i].Focus()
						continue
					}
				}
			}
		}
	}

	return *state, cmds
}

func (state *PreAuthentication) updatePreAuthentication(msg tea.Msg) State {

	newCursor := state.cursor
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return state.authMenu[newCursor]

		case "up", "k":
			if state.cursor == 0 {
				newCursor = len(state.authMenu) - 1
			} else {
				newCursor--
			}

		case "down", "j":
			if state.cursor == len(state.authMenu)-1 {
				newCursor = 0
			} else {
				newCursor++
			}
		}
	}
	state.cursor = newCursor
	return *state
}

func (state *Login) updateInputs(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, len(state.inputs))

	for i := range state.inputs {
		state.inputs[i], cmds[i] = state.inputs[i].Update(msg)
	}

	return cmds
}
