package welcome

import (
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

type Model struct {
	workspacePath string
	width         int
	height        int
}

func New(workspacePath string) Model {
	if workspacePath == "" {
		workspacePath = "./"
	}

	return Model{
		workspacePath: workspacePath,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	// ASCII art for LazyNX
	logo := `
██╗      █████╗ ███████╗██╗   ██╗███╗   ██╗██╗  ██╗
██║     ██╔══██╗╚══███╔╝╚██╗ ██╔╝████╗  ██║╚██╗██╔╝
██║     ███████║  ███╔╝  ╚████╔╝ ██╔██╗ ██║ ╚███╔╝
██║     ██╔══██║ ███╔╝    ╚██╔╝  ██║╚██╗██║ ██╔██╗
███████╗██║  ██║███████╗   ██║   ██║ ╚████║██╔╝ ██╗
╚══════╝╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═╝  ╚═══╝╚═╝  ╚═╝`

	// Style the logo with gradient colors
	logoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#002F56")).
		Bold(true).
		AlignHorizontal(lipgloss.Center).MarginTop(m.height / 4)

	subtitle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#4ECDC4")).
		Italic(true).
		AlignHorizontal(lipgloss.Center).
		Render("A Terminal User Interface for Nx workspace management")

	// Show workspace path information
	workspaceInfo := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFF")).
		Align(lipgloss.Center).
		Render("Workspace: " + m.workspacePath)

	instructions := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#888888")).
		Align(lipgloss.Center).
		Render("Press Ctrl+C to exit")

	// Create the main content
	content := lipgloss.JoinVertical(
		lipgloss.Top,
		logoStyle.Render(logo),
		"",
		subtitle,
		"",
		"",
		workspaceInfo,
		"",
		instructions,
	)

	// Create a container with proper horizontal centering
	container := lipgloss.NewStyle().
		Width(m.width).
		Height(m.height).
		AlignHorizontal(lipgloss.Center).
		AlignVertical(lipgloss.Top).
		Render(content)

	return container
}
