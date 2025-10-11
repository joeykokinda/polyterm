package ui

import "github.com/charmbracelet/lipgloss"

var (
	polyBlue   = lipgloss.Color("#6366F1")
	polyPurple = lipgloss.Color("#8B5CF6")
	polyPink   = lipgloss.Color("#EC4899")
	polyDark   = lipgloss.Color("#1E1B4B")
	polyLight  = lipgloss.Color("#C7D2FE")
	
	greenYes   = lipgloss.Color("#10B981")
	redNo      = lipgloss.Color("#EF4444")
	yellowWarn = lipgloss.Color("#F59E0B")
	grayMuted  = lipgloss.Color("#6B7280")
	
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(polyBlue).
		Background(polyDark).
		Padding(0, 2).
		MarginBottom(1)
	
	BrandStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(polyPink).
		Background(polyDark).
		Padding(0, 2)
	
	HeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(polyLight).
		Background(polyDark).
		Padding(0, 1)
	
	StatsBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(polyBlue).
		Padding(0, 2).
		MarginRight(2)
	
	StatsLabelStyle = lipgloss.NewStyle().
		Foreground(grayMuted).
		Bold(false)
	
	StatsValueStyle = lipgloss.NewStyle().
		Foreground(polyBlue).
		Bold(true)
	
	TableHeaderStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(polyPurple).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(polyBlue).
		BorderBottom(true).
		Padding(0, 1)
	
	TableCellStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#E5E7EB")).
		Padding(0, 1)
	
	YesOddsStyle = lipgloss.NewStyle().
		Foreground(greenYes).
		Bold(true).
		Padding(0, 1)
	
	NoOddsStyle = lipgloss.NewStyle().
		Foreground(redNo).
		Bold(true).
		Padding(0, 1)
	
	VolumeStyle = lipgloss.NewStyle().
		Foreground(polyPink).
		Bold(true).
		Padding(0, 1)
	
	ErrorStyle = lipgloss.NewStyle().
		Foreground(redNo).
		Bold(true)
	
	MutedStyle = lipgloss.NewStyle().
		Foreground(grayMuted)
	
	HelpStyle = lipgloss.NewStyle().
		Foreground(grayMuted).
		Italic(true).
		MarginTop(1)
	
	LoadingStyle = lipgloss.NewStyle().
		Foreground(polyPurple).
		Bold(true)
)
