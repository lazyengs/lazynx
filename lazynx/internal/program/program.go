package program

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lazyengs/lazynx/internal/models/welcome"
	"github.com/lazyengs/lazynx/pkg/nxlsclient"
	"github.com/lazyengs/lazynx/pkg/nxlsclient/commands"
	"go.uber.org/zap"
)

type activeView uint

const (
	spinnerView activeView = iota // Initial loading state
	welcomeView
)

type ProgramModel struct {
	welcomeModel  welcome.Model
	spinnerModel  spinner.Model
	activeView    activeView
	viewport      tea.WindowSizeMsg
	client        *nxlsclient.Client
	logger        *zap.SugaredLogger
	errorMsg      string
	workspacePath string
}

func createProgram(client *nxlsclient.Client, logger *zap.SugaredLogger, workspacePath string) ProgramModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return ProgramModel{
		welcomeModel:  welcome.New(workspacePath),
		spinnerModel:  s,
		client:        client,
		activeView:    spinnerView,
		logger:        logger,
		workspacePath: workspacePath,
	}
}

func (m ProgramModel) Init() tea.Cmd {
	// Start with spinner since initialization starts immediately in CLI
	return tea.Batch(
		m.welcomeModel.Init(),
		m.spinnerModel.Tick,
	)
}

func (m ProgramModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		default:
			// Reset error state on any key press
			if m.errorMsg != "" {
				m.errorMsg = ""
				return m, nil
			}
		}
	case *commands.InitializeRequestResult:
		// Initialization completed successfully
		m.activeView = welcomeView
		return m, nil
	}

	if m.activeView == welcomeView {
		m.welcomeModel, cmd = m.welcomeModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.activeView == spinnerView {
		m.spinnerModel, cmd = m.spinnerModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m ProgramModel) View() string {
	if m.activeView == welcomeView {
		return m.welcomeModel.View()
	}

	if m.activeView == spinnerView {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.spinnerModel.View(),
				"  Initializing workspace...",
			),
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Render("Workspace: "+m.workspacePath),
		)
	}

	if m.errorMsg != "" {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			m.welcomeModel.View(),
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF5722")).
				Render("Error: "+m.errorMsg),
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Render("Press any key to try again"),
		)
	}

	return ""
}

func Create(client *nxlsclient.Client, logger *zap.SugaredLogger, workspacePath string) *tea.Program {
	return tea.NewProgram(createProgram(client, logger, workspacePath), tea.WithAltScreen())
}
