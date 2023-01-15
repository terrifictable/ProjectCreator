package ui

import (
	"fmt"
	"main/util"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* Types */
type model struct {
	inputs  []textinput.Model
	spinner spinner.Model

	creating bool
	done     bool

	focused int
	err     error
}
type createProjMsg bool

/* Types */

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tea.EnterAltScreen)
}

func createProj() tea.Msg {
	err := util.NewCreateProject(nameVal, languageVal)
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
		case tea.KeyCtrlC:
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
			titleStyle.Render(util.GetLogo()),
			lipgloss.NewStyle().PaddingTop(2).PaddingLeft(10).Render(m.spinner.View()))
	}

	if m.focused == len(m.inputs)-1 && m.inputs[0].Value() != "" && m.inputs[1].Value() != "" {
		continueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).PaddingLeft(12)
	} else {
		continueStyle = lipgloss.NewStyle().Foreground(darkGray).PaddingLeft(12)
	}
	return fmt.Sprintf(`
%s

 %s
 %s

 %s  
 %s
 
 %s 
`,
		titleStyle.Render(util.GetLogo()),
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

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
