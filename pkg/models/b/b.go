package b

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ModelB struct {
	count float32
}

func Spawn() ModelB {
	return ModelB{count: 1}
}

func (m ModelB) Init() tea.Cmd {
	return nil
}

func (m ModelB) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "left":
			m.count = m.count / 2
		case "right":
			m.count = m.count * 2
		}
	}

	return m, nil
}

func (m ModelB) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			"Model B: ",
			lipgloss.NewStyle().Foreground(lipgloss.Color('3')).Render(fmt.Sprintf("$%.3f", m.count)),
		),
		lipgloss.NewStyle().Foreground(lipgloss.Color('7')).Render("---"),
		lipgloss.NewStyle().Foreground(lipgloss.Color('7')).Render("Duplicate: → | Half: ←"),
	)
}
