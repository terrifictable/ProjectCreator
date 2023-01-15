package ui

import (
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg error
)

/* Globals */
const (
	name = iota
	language
)

const (
	lightBlue = lipgloss.Color("#5cffff")
	darkBlue  = lipgloss.Color("#252f8f")
	darkGray  = lipgloss.Color("#1e1e1e")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(lightBlue).PaddingLeft(10)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray).PaddingLeft(12)
	titleStyle    = lipgloss.NewStyle().Foreground(darkBlue).
			PaddingLeft(2).
			PaddingRight(2).
			BorderStyle(lipgloss.RoundedBorder())
)

var languageVal string
var nameVal string

/* Globals */

func NewModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Project Name"
	inputs[0].Focus()
	inputs[0].CharLimit = 20
	inputs[0].Width = 30
	inputs[0].Prompt = "┃ "

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Language"
	inputs[1].Focus()
	inputs[1].CharLimit = 20
	inputs[1].Width = 30
	inputs[1].Prompt = "┃ "

	spinner := spinner.New()
	spinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("63"))

	return model{
		inputs:  inputs,
		spinner: spinner,

		focused: 0,
		err:     nil,

		done:     false,
		creating: false,
	}
}

func InitUI() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
