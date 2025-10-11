package ui

import "fmt"

func (m Model) View() string {
	if m.loading {
		return TitleStyle.Render("Polyterm • Fetching one Polymarket event… ") +
			m.spinner.View() +
			"\n\nPress q to quit."
	}

	if m.err != nil {
		return TitleStyle.Render("Polyterm") + "\n" +
			ErrorStyle.Render("Error: "+m.err.Error()) + "\n\n" +
			MutedStyle.Render("Press r to retry, q to quit.")
	}

	body := fmt.Sprintf(
		"%s\n  slug: %s\n\n%s",
		TitleStyle.Render("Polymarket Event"),
		m.event.Slug,
		MutedStyle.Render("Press r to refetch, q to quit."),
	)
	return body
}
