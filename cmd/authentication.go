package cmd

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/crypto/bcrypt"
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
	inputs []textinput.Model
	cursor int
}

type ForgotPassword struct {
}

type ResetPassword struct {
}

func initRegisterScreen() Register {

	inputs := make([]textinput.Model, 3)
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
		case 2:
			t.Placeholder = "Confirm Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		inputs[i] = t
	}

	return Register{
		inputs: inputs,
	}
}

func (r Register) renderRegisterScreen(m model) string {
	termWidth, termHeight := m.width-2, m.height-2
	register := style("Register", m.styles.magenta)
	register = lipgloss.NewStyle().PaddingBottom(1).Render(register)

	var inputStrings []string

	for _, input := range r.inputs {
		inputStrings = append(inputStrings, input.View())
	}

	inputStrings = append(inputStrings, wrapWithCursor(r.cursor == len(r.inputs), focusedButton, m.styles.toEnter))

	joined := lipgloss.JoinVertical(lipgloss.Left, append([]string{register}, inputStrings...)...)
	s := lipgloss.NewStyle().Align(lipgloss.Left).Render(joined)
	centeredText := lipgloss.Place(termWidth, termHeight, lipgloss.Center, lipgloss.Center, s)

	return centeredText
}

func (state *Register) updateInputs(msg tea.Msg) []tea.Cmd {
	cmds := make([]tea.Cmd, len(state.inputs))

	for i := range state.inputs {
		state.inputs[i], cmds[i] = state.inputs[i].Update(msg)
	}

	return cmds
}

func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func hashPassword(password, salt string) (string, error) {
	saltedPassword := password + salt

	bytes, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash, salt string) bool {
	saltedPassword := password + salt

	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(saltedPassword))
	return err == nil
}
