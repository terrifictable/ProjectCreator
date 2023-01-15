package ui

import (
	"fmt"
	"main/util"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/* Types */
type keymap struct {
	tab        key.Binding
	unselect   key.Binding
	unsel_quit key.Binding
	quit       key.Binding
}
type model struct {
	inputs  []textinput.Model
	spinner spinner.Model

	creating    bool
	selected    bool
	invalidLang bool

	focused int
	err     error

	help   help.Model
	keymap keymap
}
type createProjMsg bool

/* Types */

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, tea.EnterAltScreen)
}

func createProj() tea.Msg {
	err := util.NewCreateProject(nameVal, languageVal)
	if err != nil {
		return createProjMsg(false)
	}

	return createProjMsg(true)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			for i := range m.inputs {
				if m.inputs[i].Focused() {
					m.inputs[i].Blur()
				}
				m.selected = false
			}
			return m, nil

		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 && m.inputs[0].Value() != "" && m.inputs[1].Value() != "" {
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

		if !m.selected && msg.String() == "q" {
			return m, tea.Quit
		}

		if m.focused == 1 {
			m.invalidLang = false
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
		m.selected = true

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case createProjMsg:
		if msg {
			return m, tea.Quit
		}
		m.creating = false
		m.invalidLang = true

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
		continueStyle = continueStyle.Copy().Foreground(lipgloss.Color("#ffffff"))
	} else {
		continueStyle = continueStyle.Copy().Foreground(darkGray)
	}

	helpText := "tab: select next • return: submit • esc: unselect • "
	if !m.selected {
		helpText += "q: exit"
	} else {
		helpText += "ctrl+c: exit"
	}

	errorText := ""
	if m.invalidLang {
		errorText = "Unknown Language"
	}

	return fmt.Sprintf(`
%s

%s
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
		helpStyle.Render(continueStyle.Render("Continue")),
		errorStyle.Render(errorText),
		helpStyle.Copy().PaddingLeft(2).Render(helpText),
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
