package cmd

import (
	"fmt"
	"reflect"
	"termtyper/database"

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
	form              *huh.Form
	formData          *LoginFormData
	isFormInitialized bool
	errorMessage      string
}

type Register struct {
	context           *database.Context
	form              *huh.Form
	formData          *RegisterFormData
	isFormInitialized bool
	errorMessage      string
}

type GuestLogin struct{}

type LoginFormData struct {
	email    string
	password string
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

func initLoginScreen(errorMessage string) Login {
	data := &LoginFormData{}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Email").
				Value(&data.email),
			huh.NewInput().
				Title("password").
				Value(&data.password).
				EchoMode(huh.EchoModePassword),
		),
	)
	return Login{
		form:              form,
		formData:          data,
		isFormInitialized: false,
		errorMessage:      errorMessage,
	}
}

func initRegisterScreen(
	context *database.Context,
	errorMessage string,
	email string,
	password string,
	confirmPassword string,
) Register {
	data := &RegisterFormData{
		email:           email,
		password:        password,
		confirmPassword: confirmPassword,
	}
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Email").
				Value(&data.email).
				Validate(func(str string) error {
					_, err := mail.ParseAddress(str)
					if err != nil {
						return fmt.Errorf("invalid email address")
					} else if database.CheckEmailExists(context.UserRepository, str) {
						return fmt.Errorf("email already exists")
					}
					return nil
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
				EchoMode(huh.EchoModePassword),
		),
	)

	return Register{
		context:           context,
		form:              form,
		formData:          data,
		isFormInitialized: false,
		errorMessage:      errorMessage,
	}
}

func initPreAuthentication(context *database.Context) PreAuthentication {
	return PreAuthentication{
		authMenu: []State{
			initRegisterScreen(context, "", "", "", ""),
			initLoginScreen(""),
			GuestLogin{},
		},
		cursor: 0,
	}
}

func (l Login) renderLoginScreen(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	login := style("Login"+" "+l.errorMessage, m.styles.magenta)
	login = lipgloss.NewStyle().PaddingBottom(1).Render(login)

	joined := lipgloss.JoinVertical(lipgloss.Left, []string{login, l.form.View()}...)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	return centeredText
}

func (r Register) renderRegisterScreen(m model) string {
	termWidth, termHeight := m.width-2, m.height-2
	register := style("Register"+" "+r.errorMessage, m.styles.magenta)
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

func (state *PreAuthentication) updatePreAuthentication(msg tea.Msg) State {

	newCursor := state.cursor
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			if _, ok := state.authMenu[newCursor].(GuestLogin); ok {
				guestUser := database.ApplicationUser{
					Id:       -1,
					Username: "Guest",
					Config: &database.UserConfig{
						Time:  30,
						Words: 30,
					},
				}
				return initMainMenu(&guestUser)
			}
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
