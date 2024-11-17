package a

import (
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModelA struct {
	count int
}

func Spawn() ModelA {
	return ModelA{count: 0}
}

func (m ModelA) Spawn() tea.Model {
	return ModelA{count: 0}
}

func (m ModelA) Init() tea.Cmd {
	return nil
}

func (m ModelA) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.count = m.count + 1
		case "down":
			m.count = max(0, m.count-1)
		}
	}

	return m, nil
}

func (m ModelA) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			"Model A: ",
			lipgloss.NewStyle().Foreground(lipgloss.Color('5')).Render(strconv.Itoa(m.count)),
		),

		lipgloss.NewStyle().Foreground(lipgloss.Color('7')).Render("---"),
		lipgloss.NewStyle().Foreground(lipgloss.Color('7')).Render("Increment: ↑ | Decrement: ↓"),
	)
}
