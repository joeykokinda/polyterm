package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
)

type gammaEvent struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type eventsResp struct {
	// Gamma commonly returns arrays; in some builds it returns objects with "data".
	// We'll support both safely.
	Data  []gammaEvent `json:"data"`
	Items []gammaEvent `json:"items"`
}

type fetched struct {
	ev *gammaEvent
	err error
}

type model struct {
	spin     spinner.Model
	loading  bool
	event    *gammaEvent
	err      error
	width    int
	height   int
}

//init
func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	return model{
		spin:    s,
		loading: true,
	}
}

func fetchOneEventCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()

		req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
			"https://gamma-api.polymarket.com/events?closed=false&limit=1", nil)
		req.Header.Set("User-Agent", "polyterm/0.0.1")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fetched{err: err}
		}
		defer resp.Body.Close()
		b, _ := io.ReadAll(resp.Body)

		var er eventsResp
		_ = json.Unmarshal(b, &er)

		var list []gammaEvent
		switch {
		case len(er.Data) > 0:
			list = er.Data
		case len(er.Items) > 0:
			list = er.Items
		default:
			// try raw array fallback
			if err2 := json.Unmarshal(b, &list); err2 != nil {
				return fetched{err: fmt.Errorf("unexpected response shape")}
			}
		}

		if len(list) == 0 {
			return fetched{err: fmt.Errorf("no events returned")}
		}
		ev := list[0]
		return fetched{ev: &ev}
	}
}

// ---- tea.Model
func (m model) Init() tea.Cmd {
	return tea.Batch(m.spin.Tick, fetchOneEventCmd())
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case fetched:
		m.loading = false
		m.err = msg.err
		m.event = msg.ev
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.loading = true
			m.err = nil
			m.event = nil
			return m, tea.Batch(m.spin.Tick, fetchOneEventCmd())
		}
	}
	if m.loading {
		var cmd tea.Cmd
		m.spin, cmd = m.spin.Update(msg)
		return m, cmd
	}
	return m, nil
}

var (
	titleStyle = lipgloss.NewStyle().Bold(true).MarginBottom(1)
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	muted      = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
)

func (m model) View() string {
	if m.loading {
		return titleStyle.Render("Polyterm • Fetching one Polymarket event… ") + m.spin.View() +
			"\n\nPress q to quit."
	}
	if m.err != nil {
		return titleStyle.Render("Polyterm") + "\n" +
			errStyle.Render("Error: "+m.err.Error()) + "\n\n" +
			muted.Render("Press r to retry, q to quit.")
	}
	body := fmt.Sprintf(
		"%s\n  slug: %s\n\n%s",
		titleStyle.Render("Polymarket Event"),
		m.event.Slug,
		muted.Render("Press r to refetch, q to quit."),
	)
	return body
}

func main() {
	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("error:", err)
	}
}
