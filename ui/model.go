package ui

import (
	"context"

	"polyterm/api"
	"polyterm/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
)

type Model struct {
	spinner spinner.Model
	loading bool
	event   *types.Event
	err     error
	width   int
	height  int
}

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return Model{
		spinner: s,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, fetchEventCmd())
}

func fetchEventCmd() tea.Cmd {
	return func() tea.Msg {
		event, err := api.FetchOneEvent(context.Background())
		return types.FetchResult{Event: event, Err: err}
	}
}
