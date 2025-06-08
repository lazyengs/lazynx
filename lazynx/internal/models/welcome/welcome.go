package welcome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeModel struct {
	form          *huh.Form
	workspacePath string
}

type WorkspacePathMsg struct {
	Path string
}

func New() WelcomeModel {
	workspacePath := "./"

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nx Workspace Path").
				Description("Enter the path to your Nx workspace").
				Placeholder("./").
				Value(&workspacePath).
				CharLimit(256),
		),
	).WithShowHelp(false)

	return WelcomeModel{
		form:          form,
		workspacePath: workspacePath,
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
	}

	if m.form.State == huh.StateCompleted {
		workspacePath := m.workspacePath
		if workspacePath == "" {
			workspacePath = "./"
		}
		return m, func() tea.Msg {
			return WorkspacePathMsg{Path: workspacePath}
		}
	}

	return m, cmd
}

func (m WelcomeModel) View() string {
	title := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF6B6B")).
		Bold(true).
		Render("Welcome to LazyNX!")

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Render("A TUI for Nx workspace management")

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		lipgloss.NewStyle().Align(lipgloss.Center).Render(title),
		lipgloss.NewStyle().Align(lipgloss.Center).Render(subtitle),
		"",
		"",
		m.form.View(),
	)
}
