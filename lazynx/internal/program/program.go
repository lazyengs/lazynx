package program

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lazyengs/lazynx/internal/components"
	"github.com/lazyengs/lazynx/internal/models/a"
	"github.com/lazyengs/lazynx/internal/models/b"
	"github.com/lazyengs/lazynx/internal/models/c"
	"github.com/lazyengs/lazynx/internal/utils"
)

type (
	ModelId  int
	ModelMap map[ModelId]tea.Model
)

type ProgramModel struct {
	models         ModelMap
	viewport       tea.WindowSizeMsg
	focusedModelId ModelId
}

const (
	ViewA ModelId = iota
	ViewB
	ViewC
)

func createProgram() ProgramModel {
	// Model ID keeps track of the current model the user is
	// focused in as well as the key to get the actual model
	// model from the map.
	currentModelId := ViewA

	models := make(ModelMap)
	models[ViewA] = a.Spawn()
	models[ViewB] = b.Spawn()
	models[ViewC] = c.Spawn()

	return ProgramModel{
		focusedModelId: currentModelId,
		models:         models,
	}
}

func (m ProgramModel) Init() tea.Cmd {
	return nil
}

func (m ProgramModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// In order to make the program responsive, we'll need to keep track of
		// the available viewport, as every model's layout will be derived from
		// these dimensions.
		m.viewport = msg
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			m.focusedModelId = ModelId(int(m.focusedModelId + 1))
			if int(m.focusedModelId) > len(m.models)-1 {
				m.focusedModelId = 0
			}
			return m, nil
		case "shift+tab":
			m.focusedModelId = ModelId(int(m.focusedModelId - 1))
			if int(m.focusedModelId) < 0 {
				m.focusedModelId = ModelId(len(m.models) - 1)
			}
			return m, nil
		}
	}

	// If none of the above cases are met, we'll deletate the message
	// to the currently focused model and let it handle it separately
	// and perform any updates necessary.
	currentModel := m.getModel(m.focusedModelId)
	updatedModel, cmd := currentModel.Update(msg)

	// The resulting model reflects the updated state, and it must be
	// revalidated within our set of models, which is essentialy done
	// by hot-swapping the model with the updated one.
	m.updateModel(m.focusedModelId, updatedModel)

	return m, cmd
}

func (m ProgramModel) View() string {
	w1, w2 := utils.SafeHalves(m.viewport.Width)

	leftPane := lipgloss.JoinVertical(
		lipgloss.Top,
		m.getModelStyle(ViewA).Render(m.getModel(ViewA).View()),
		m.getModelStyle(ViewB).Render(m.getModel(ViewB).View()),
	)
	rightPane := m.getModelStyle(ViewC).Render(m.getModel(ViewC).View())

	bottomText := lipgloss.JoinHorizontal(
		lipgloss.Right,
		lipgloss.NewStyle().Width(w1).Foreground(lipgloss.Color('7')).Align(lipgloss.Left).Render("Next view: ⇥ | Previous view: ⇧⇥"),
		lipgloss.NewStyle().Width(w2).Foreground(lipgloss.Color('7')).Align(lipgloss.Right).Render("Lazynx v0.1.0"),
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.JoinHorizontal(
			lipgloss.Right,
			leftPane,
			rightPane,
		),
		bottomText,
	)
}

func (m ProgramModel) getModel(id ModelId) tea.Model {
	return m.models[id]
}

func (m ProgramModel) updateModel(id ModelId, model tea.Model) {
	m.models[id] = model
}

func (m ProgramModel) getModelStyle(id ModelId) lipgloss.Style {
	component := components.InactiveBox

	if id == m.focusedModelId {
		component = components.ActiveBox
	}

	// These are all model-related styles. The model decides *what* to render,
	// but the program decides *how* to render it.
	//
	// This involves:
	// - Dimensions
	// - Positioning
	// - Focused states
	//
	// Here, we're allocating the entire viewport width for placing all panes
	// but reserving a single row at the bottom, for general info & controls.
	availableWidth := m.viewport.Width
	availableHeight := m.viewport.Height - 1

	// Since we're on a TUI splitting a dimension in half isn't guaranteed to
	// result in two equal halves -or thirds-, so we'll use a helper function
	// to ensure each pane gets the right dimensions based on their position.
	//
	// This very manual approach to layouts may be revisited in the future in
	// favor of a more flexible and automated layout system.
	//
	// See: https://github.com/charmbracelet/lipgloss/issues/166
	w1, w2 := utils.SafeHalves(availableWidth)
	h1, h2 := utils.SafeHalves(availableHeight)

	h3 := availableHeight

	switch id {
	case ViewA:
		component = component.Width(w1 - 2).Height(h1 - 2)
	case ViewB:
		component = component.Width(w1 - 2).Height(h2 - 2)
	case ViewC:
		component = component.Width(w2 - 2).Height(h3 - 2)
	}

	return component
}

func Create() *tea.Program {
	return tea.NewProgram(createProgram(), tea.WithAltScreen())
}
