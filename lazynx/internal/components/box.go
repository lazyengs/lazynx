package components

import (
	"github.com/charmbracelet/lipgloss"
)

var Box = lipgloss.NewStyle().Border(lipgloss.RoundedBorder())

var (
	ActiveBox   = Box.BorderForeground(lipgloss.Color('2'))
	InactiveBox = Box.BorderForeground(lipgloss.Color('8'))
)
