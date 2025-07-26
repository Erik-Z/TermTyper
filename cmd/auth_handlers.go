package cmd

import (
	"fmt"
	"net/mail"
	"termtyper/database"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type LoginHandler struct {
	*BaseStateHandler
	form              *huh.Form
	formData          *LoginFormData
	isFormInitialized bool
	errorMessage      string
}

func NewLoginHandler(errorMessage string) *LoginHandler {
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

	return &LoginHandler{
		BaseStateHandler:  NewBaseStateHandler(StateLogin),
		form:              form,
		formData:          data,
		isFormInitialized: false,
		errorMessage:      errorMessage,
	}
}

func (h *LoginHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	var commands []tea.Cmd
	if !h.isFormInitialized {
		initCmd := h.form.Init()
		commands = append(commands, initCmd)
		h.isFormInitialized = true
	}

	updatedForm, formCmd := h.form.Update(msg)
	if f, ok := updatedForm.(*huh.Form); ok {
		h.form = f
		commands = append(commands, formCmd)
	}

	if h.form.State == huh.StateCompleted {
		authUser, err := database.AuthenticateUser(context.model.context.UserRepository, h.formData.email, h.formData.password)
		if err == nil && authUser != nil {
			context.model.session.User = authUser
			context.model.session.Authenticated = true
			newState := NewMainMenuHandler(initMainMenu(authUser))
			return newState, tea.Batch(commands...)
		} else {
			newState := NewLoginHandler("❌" + err.Error())
			newState.formData.email = h.formData.email
			return newState, tea.Batch(commands...)
		}
	}
	return h, tea.Batch(commands...)
}

func (h *LoginHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	login := style("Login"+" "+h.errorMessage, m.styles.magenta)
	login = lipgloss.NewStyle().PaddingBottom(1).Render(login)

	joined := lipgloss.JoinVertical(lipgloss.Left, []string{login, h.form.View()}...)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	return centeredText
}

func (h *LoginHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateLogin]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}

type RegisterHandler struct {
	*BaseStateHandler
	context           *database.Context
	form              *huh.Form
	formData          *RegisterFormData
	isFormInitialized bool
	errorMessage      string
}

func NewRegisterHandler(context *database.Context,
	errorMessage string,
	email string,
	password string,
	confirmPassword string) *RegisterHandler {
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

	return &RegisterHandler{
		BaseStateHandler:  NewBaseStateHandler(StateRegister),
		context:           context,
		form:              form,
		formData:          data,
		isFormInitialized: false,
		errorMessage:      errorMessage,
	}
}

// TODO: form only renders after mouse movement. Fix this. Should be similar to replay.
func (h *RegisterHandler) HandleInput(msg tea.Msg, context *StateContext) (StateHandler, tea.Cmd) {
	var commands []tea.Cmd
	if !h.isFormInitialized {
		initCmd := h.form.Init()
		commands = append(commands, initCmd)
		h.isFormInitialized = true
	}

	updatedForm, formCmd := h.form.Update(msg)
	if f, ok := updatedForm.(*huh.Form); ok {
		h.form = f
		commands = append(commands, formCmd)
	}

	if h.form.State == huh.StateCompleted {
		if h.formData.password != h.formData.confirmPassword {
			newHandler := NewRegisterHandler(
				h.context,
				"❌ Passwords must match",
				h.formData.email,
				h.formData.password,
				h.formData.confirmPassword,
			)

			return newHandler, tea.Batch(commands...)
		}

		newUser, err := database.CreateUser(context.model.context.UserRepository, h.formData.email, h.formData.password)
		if err != nil {
			newHandler := NewRegisterHandler(
				h.context,
				"❌ "+err.Error(),
				h.formData.email,
				h.formData.password,
				h.formData.confirmPassword,
			)

			return newHandler, tea.Batch(commands...)
		}
		context.model.session.User = newUser
		context.model.session.Authenticated = true
		MainMenuHandler := NewMainMenuHandler(initMainMenu(newUser))
		return MainMenuHandler, tea.Batch(commands...)
	}
	return h, tea.Batch(commands...)
}

func (h *RegisterHandler) Render(m *model) string {
	termWidth, termHeight := m.width-2, m.height-2
	register := style("Register"+" "+h.errorMessage, m.styles.magenta)
	register = lipgloss.NewStyle().PaddingBottom(1).Render(register)

	joined := lipgloss.JoinVertical(lipgloss.Center, []string{register, h.form.View()}...)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	return centeredText
}

func (h *RegisterHandler) ValidateTransition(to StateType, context *StateContext) bool {
	validTransitions := context.transitionMap[StateRegister]
	for _, validState := range validTransitions {
		if validState == to {
			return true
		}
	}
	return false
}
