package main

import (
	"fmt"

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

/*
cdd8f2, 9cb0d7, 6d88bc, 4161a1, 153a85
#cdd8f2, #9cb0d7, #6d88bc, #4161a1, #153a85

cdd8f2,9cb0d7,6d88bc,4161a1,153a85
#cdd8f2,#9cb0d7,#6d88bc,#4161a1,#153a85

["cdd8f2","9cb0d7","6d88bc","4161a1","153a85"]
["#cdd8f2","#9cb0d7","#6d88bc","#4161a1","#153a85"]
*/
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

/* MODEL */
/* TYPES */
type model struct {
	inputs  []textinput.Model
	spinner spinner.Model

	creating bool
	done     bool

	focused int
	err     error
}
type createProjMsg bool

/* TYPES */

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tea.EnterAltScreen)
}

func createProj() tea.Msg {
	err := NewCreateProject(nameVal, languageVal)
	if err != nil {
		panic(err)
	}

	return createProjMsg(true)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				m.creating = true

				nameVal = m.inputs[0].Value()
				languageVal = m.inputs[1].Value()

				return m, createProj
			}
			m.nextInput()
		case tea.KeyCtrlC /*, tea.KeyEsc */ :
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case createProjMsg:
		return m, tea.Quit

	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.creating {
		return fmt.Sprintf(`
%s

 %s Creating Project
`,
			titleStyle.Render(GetLogo()),
			lipgloss.NewStyle().PaddingTop(2).PaddingLeft(10).Render(m.spinner.View()))
	}
	return fmt.Sprintf(`
%s

 %s
 %s

 %s  
 %s
 
 %s 
`,
		titleStyle.Render(GetLogo()),
		inputStyle.Width(30).Render("Project Name"),
		lipgloss.NewStyle().PaddingLeft(10).Render(m.inputs[name].View()),
		inputStyle.Width(8).Render("Language"),
		lipgloss.NewStyle().PaddingLeft(10).Render(m.inputs[language].View()),
		continueStyle.Render("Continue ->"),
	) + "\n"
}

func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

/* MODEL */

func NewModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[0] = textinput.New()
	inputs[0].Placeholder = "Project Name"
	inputs[0].Focus()
	inputs[0].CharLimit = 20
	inputs[0].Width = 30
	inputs[0].Prompt = ""

	inputs[1] = textinput.New()
	inputs[1].Placeholder = "Language"
	inputs[1].Focus()
	inputs[1].CharLimit = 20
	inputs[1].Width = 30
	inputs[1].Prompt = ""

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

func initUI() error {
	p := tea.NewProgram(NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		return err
	}
	return nil
}
