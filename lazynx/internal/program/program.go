package program

import (
	"context"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lazyengs/lazynx/internal/models/welcome"
	"github.com/lazyengs/lazynx/internal/utils"
	"github.com/lazyengs/lazynx/pkg/nxlsclient"
	"go.uber.org/zap"
)

type ProgramModel struct {
	welcomeModel  tea.Model
	spinner       spinner.Model
	viewport      tea.WindowSizeMsg
	client        *nxlsclient.Client
	logger        *zap.SugaredLogger
	initialized   bool
	initializing  bool
	errorMsg      string
	workspacePath string
}

type initResultMsg struct {
	err error
}

func createProgram(client *nxlsclient.Client, logger *zap.SugaredLogger) ProgramModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return ProgramModel{
		welcomeModel: welcome.New(),
		spinner:      s,
		client:       client,
		logger:       logger,
		initialized:  false,
		initializing: false,
	}
}

func (m ProgramModel) Init() tea.Cmd {
	return m.welcomeModel.Init()
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
	case welcome.WorkspacePathMsg:
		if !m.initializing {
			m.workspacePath = msg.Path
			m.initializing = true
			return m, tea.Batch(
				m.spinner.Tick,
				m.initializeClient(),
			)
		}
	case initResultMsg:
		m.initializing = false
		if msg.err != nil {
			m.errorMsg = msg.err.Error()
		} else {
			m.initialized = true
		}
		return m, nil
	}

	if !m.initialized && !m.initializing {
		m.welcomeModel, cmd = m.welcomeModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	if m.initializing {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m ProgramModel) initializeClient() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		err := utils.InitializeNxlsclient(ctx, m.client, m.workspacePath, nil, m.logger)
		return initResultMsg{err: err}
	}
}

func (m ProgramModel) View() string {
	if m.initialized {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#00FF00")).
				Bold(true).
				Render("✓ Workspace initialized successfully!"),
			"",
			"Workspace: "+m.workspacePath,
			"",
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#888888")).
				Render("Press Ctrl+C to exit"),
		)
	}

	if m.initializing {
		return lipgloss.JoinVertical(
			lipgloss.Center,
			"",
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.spinner.View(),
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

	return m.welcomeModel.View()
}

func Create(client *nxlsclient.Client, logger *zap.SugaredLogger) *tea.Program {
	return tea.NewProgram(createProgram(client, logger), tea.WithAltScreen())
}
