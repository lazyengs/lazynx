package tui

import (
	"github.com/charmbracelet/bubbles/v2/key"
	"github.com/charmbracelet/bubbles/v2/spinner"
	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"

	"github.com/lazyengs/lazynx/internal/tui/components"
	"github.com/lazyengs/lazynx/internal/tui/layout"
	"github.com/lazyengs/lazynx/internal/tui/models/welcome"
	"github.com/lazyengs/lazynx/internal/tui/utils"
	"github.com/lazyengs/lazynx/pkg/nxlsclient"
	"github.com/lazyengs/lazynx/pkg/nxlsclient/commands"
	"go.uber.org/zap"
)

type activeView uint

const (
	spinnerView activeView = iota // Initial loading state
	welcomeView
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Left  key.Binding
	Right key.Binding
	Help  key.Binding
	Quit  key.Binding
}

var globalKeys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q/ctrl+c", "quit"),
	),
}

func getKeysForView(view activeView) []key.Binding {
	switch view {
	case welcomeView:
		// For welcome view, show most relevant keys
		return []key.Binding{
			globalKeys.Help,
			globalKeys.Quit,
		}
	case spinnerView:
		// For spinner view, show minimal keys
		return []key.Binding{
			globalKeys.Help,
			globalKeys.Quit,
		}
	default:
		// For unknown views, show all keys using reflection
		return utils.KeyMapToSlice(globalKeys)
	}
}

type ProgramModel struct {
	welcomeModel welcome.Model
	spinnerModel spinner.Model
	activeView   activeView

	showHelp      bool
	helpComponent *components.HelpComponent

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

	helpComp := components.NewHelpComponent()

	return ProgramModel{
		welcomeModel:  welcome.New(workspacePath),
		spinnerModel:  s,
		helpComponent: helpComp,
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
		m.helpComponent.Init(),
	)
}

func (m ProgramModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport = msg
		helpModel, helpCmd := m.helpComponent.Update(msg)
		m.helpComponent = helpModel.(*components.HelpComponent)
		cmds = append(cmds, helpCmd)
		m.welcomeModel, cmd = m.welcomeModel.Update(msg)
		cmds = append(cmds, cmd)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, globalKeys.Help):
			m.showHelp = !m.showHelp
			if m.showHelp {
				// Update help component with current view's key bindings
				m.helpComponent.SetBindings(getKeysForView(m.activeView))
			}
			return m, nil
		case key.Matches(msg, globalKeys.Quit):
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
	var baseView string

	if m.activeView == welcomeView {
		content := m.welcomeModel.View()
		helpFooter := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Render("Press ? for help")

		baseView = lipgloss.JoinVertical(
			lipgloss.Center,
			content,
			"",
			helpFooter,
		)
	} else if m.activeView == spinnerView {
		baseView = lipgloss.JoinVertical(
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
	} else if m.errorMsg != "" {
		baseView = lipgloss.JoinVertical(
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
	} else {
		baseView = ""
	}

	// Ensure base view fills the viewport
	if baseView != "" {
		baseView = lipgloss.NewStyle().
			Width(m.viewport.Width).
			Height(m.viewport.Height).
			Render(baseView)
	}

	// If help is shown, create a true overlay that preserves the background
	if m.showHelp {
		// Ensure base view fills the entire viewport
		styledBaseView := lipgloss.NewStyle().
			Width(m.viewport.Width).
			Height(m.viewport.Height).
			Render(baseView)

		// Render the help modal content
		modal := m.helpComponent.View()

		// Add instruction text below the modal
		instructions := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			AlignHorizontal(lipgloss.Center).
			Render("Press ? again to close help")

		modalContent := lipgloss.JoinVertical(
			lipgloss.Center,
			modal,
			"",
			instructions,
		)

		// Calculate center position for the modal
		modalWidth := lipgloss.Width(modalContent)
		modalHeight := lipgloss.Height(modalContent)

		// Center the modal on the screen
		x := (m.viewport.Width - modalWidth) / 2
		y := (m.viewport.Height - modalHeight) / 2

		// Use PlaceOverlay to place the modal on top of the base view
		// This will preserve the background while showing the modal on top
		overlay := layout.PlaceOverlay(x, y, modalContent, styledBaseView)

		return overlay
	}

	return baseView
}

func Create(client *nxlsclient.Client, logger *zap.SugaredLogger, workspacePath string) *tea.Program {
	return tea.NewProgram(createProgram(client, logger, workspacePath), tea.WithAltScreen())
}
