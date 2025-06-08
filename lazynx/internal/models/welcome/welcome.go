package welcome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type WelcomeModel struct {
	workspacePath string
	width         int
	height        int
}

func New(workspacePath string) WelcomeModel {
	if workspacePath == "" {
		workspacePath = "./"
	}

	return WelcomeModel{
		workspacePath: workspacePath,
	}
}

func (m WelcomeModel) Init() tea.Cmd {
	return nil
}

func (m WelcomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m WelcomeModel) View() string {
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
		AlignHorizontal(lipgloss.Center)

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
		lipgloss.Center,
		logoStyle.Render(logo),
		"",
		subtitle,
		"",
		"",
		workspaceInfo,
		"",
		instructions,
	)

	// Calculate max width (80 chars or terminal width, whichever is smaller)
	maxWidth := 80
	if m.width > 0 && m.width < maxWidth {
		maxWidth = m.width - 4 // Leave some margin
	}

	// Create a container with proper centering
	container := lipgloss.NewStyle().
		Width(maxWidth).
		Height(m.height).
		Align(lipgloss.Center).
		AlignVertical(lipgloss.Center).
		Render(content)

	return container
}
