package components

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpComponent struct {
	width    int
	height   int
	keys     []key.Binding
	maxWidth int
}

func NewHelpComponent() *HelpComponent {
	return &HelpComponent{
		maxWidth: 90,
	}
}

func (h *HelpComponent) Init() tea.Cmd {
	return nil
}

func (h *HelpComponent) SetBindings(keys []key.Binding) {
	h.keys = removeDuplicateBindings(keys)
}

func (h *HelpComponent) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h.width = h.maxWidth
		h.height = msg.Height
		if msg.Width < h.maxWidth {
			h.width = msg.Width - 4 // Leave margin for borders
		}
	}
	return h, nil
}

func removeDuplicateBindings(bindings []key.Binding) []key.Binding {
	seen := make(map[string]struct{})
	result := make([]key.Binding, 0, len(bindings))

	// Process bindings in reverse order
	for i := len(bindings) - 1; i >= 0; i-- {
		b := bindings[i]
		k := strings.Join(b.Keys(), " ")
		if _, ok := seen[k]; ok {
			// duplicate, skip
			continue
		}
		seen[k] = struct{}{}
		// Add to the beginning of result to maintain original order
		result = append([]key.Binding{b}, result...)
	}

	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (h *HelpComponent) renderContent() string {
	baseStyle := lipgloss.NewStyle()

	helpKeyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#4ECDC4")).
		PaddingRight(1)

	helpDescStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFFFF"))

	// Compile list of bindings to render
	bindings := h.keys

	// Enumerate through each group of bindings, populating a series of
	// pairs of columns, one for keys, one for descriptions
	var (
		pairs []string
		width int
		rows  = 10 // Reduced rows per column for better fit
	)

	for i := 0; i < len(bindings); i += rows {
		var (
			keys  []string
			descs []string
		)
		for j := i; j < min(i+rows, len(bindings)); j++ {
			keyStr := strings.Join(bindings[j].Keys(), "/")
			keys = append(keys, helpKeyStyle.Render(keyStr))
			descs = append(descs, helpDescStyle.Render(bindings[j].Help().Desc))
		}

		// Render pair of columns; beyond the first pair, render a three space
		// left margin, in order to visually separate the pairs.
		var cols []string
		if len(pairs) > 0 {
			cols = []string{baseStyle.Render("   ")}
		}

		// Normalize column widths
		maxDescWidth := 0
		for _, desc := range descs {
			if maxDescWidth < lipgloss.Width(desc) {
				maxDescWidth = lipgloss.Width(desc)
			}
		}
		for i := range descs {
			remainingWidth := maxDescWidth - lipgloss.Width(descs[i])
			if remainingWidth > 0 {
				descs[i] = descs[i] + baseStyle.Render(strings.Repeat(" ", remainingWidth))
			}
		}

		maxKeyWidth := 0
		for _, key := range keys {
			if maxKeyWidth < lipgloss.Width(key) {
				maxKeyWidth = lipgloss.Width(key)
			}
		}
		for i := range keys {
			remainingWidth := maxKeyWidth - lipgloss.Width(keys[i])
			if remainingWidth > 0 {
				keys[i] = keys[i] + baseStyle.Render(strings.Repeat(" ", remainingWidth))
			}
		}

		cols = append(cols,
			strings.Join(keys, "\n"),
			strings.Join(descs, "\n"),
		)

		pair := baseStyle.Render(lipgloss.JoinHorizontal(lipgloss.Top, cols...))
		// check whether it exceeds the maximum width avail (the width of the
		// terminal, subtracting 2 for the borders).
		width += lipgloss.Width(pair)
		if width > h.width-2 {
			break
		}
		pairs = append(pairs, pair)
	}

	// Handle multiple columns with proper alignment
	if len(pairs) > 1 {
		prefix := pairs[:len(pairs)-1]
		lastPair := pairs[len(pairs)-1]
		if len(prefix) > 0 {
			prefix = append(prefix, lipgloss.Place(
				lipgloss.Width(lastPair),   // width
				lipgloss.Height(prefix[0]), // height
				lipgloss.Left,              // x
				lipgloss.Top,               // y
				lastPair,                   // content
				lipgloss.WithWhitespaceBackground(lipgloss.Color("#1a1a1a")),
			))
		}
		content := baseStyle.Width(h.width).Render(
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				prefix...,
			),
		)
		return content
	}

	// Join pairs of columns
	content := baseStyle.Width(h.width).Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			pairs...,
		),
	)
	return content
}

func (h *HelpComponent) View() string {
	baseStyle := lipgloss.NewStyle()

	content := h.renderContent()
	header := baseStyle.
		Bold(true).
		Width(lipgloss.Width(content)).
		Foreground(lipgloss.Color("#4ECDC4")).
		AlignHorizontal(lipgloss.Center).
		Render("LazyNX - Keyboard Shortcuts")

	return baseStyle.
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4ECDC4")).
		Background(lipgloss.Color("#1a1a1a")).
		Width(h.width).
		Render(
			lipgloss.JoinVertical(lipgloss.Center,
				header,
				"",
				content,
			),
		)
}