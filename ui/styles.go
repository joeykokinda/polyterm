package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	MutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)
