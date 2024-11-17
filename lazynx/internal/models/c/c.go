package c

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModelC struct {
	status bool
}

func Spawn() ModelC {
	return ModelC{status: false}
}

func (m ModelC) Init() tea.Cmd {
	return nil
}

func (m ModelC) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "t", " ", "enter":
			m.status = !m.status
		}

	}

	return m, nil
}

func (m ModelC) View() string {
	status := "Off"
	if m.status {
		status = "On"
	}

	color := '1'
	if m.status {
		color = '2'
	}

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			"Model C: ",
			lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(status),
		),
		lipgloss.NewStyle().Foreground(lipgloss.Color('7')).Render("---"),
		lipgloss.NewStyle().Foreground(lipgloss.Color('7')).Render("Toggle on/off: t | space | enter"),
	)
}
