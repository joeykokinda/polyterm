package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"polyterm/api"
	"polyterm/types"

	"github.com/charmbracelet/lipgloss"
)

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing...\n"
	}

	if m.loading && len(m.markets) == 0 {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			BrandStyle.Render("POLYTERM"),
			"",
			LoadingStyle.Render("Loading Polymarket data... ")+m.spinner.View(),
			"",
			HelpStyle.Render("Press q to quit"),
		)
	}

	if m.err != nil {
		return lipgloss.JoinVertical(
			lipgloss.Left,
			"",
			BrandStyle.Render("POLYTERM"),
			"",
			ErrorStyle.Render("Error: "+m.err.Error()),
			"",
			HelpStyle.Render("Press r to retry | q to quit"),
		)
	}

	if m.currentView == viewDetail && m.selectedMarket >= 0 && m.selectedMarket < len(m.markets) {
		return m.renderMarketDetail()
	}

	switch m.currentPage {
	case pageMarkets:
		return m.renderMarketsPage()
	case pageStats:
		return m.renderStatsPage()
	default:
		return m.renderMarketsPage()
	}
}

func (m Model) renderMarketsPage() string {
	header := m.renderHeader()
	tabs := m.renderTabs()
	filterBar := m.renderFilterBar()
	stats := m.renderStats()
	table := m.renderTable()
	help := m.renderHelp()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		header,
		tabs,
		filterBar,
		stats,
		"",
		table,
		help,
	)
}

func (m Model) renderStatsPage() string {
	header := m.renderHeader()
	tabs := m.renderTabs()
	statsDetail := m.renderAdvancedStats()
	help := m.renderHelp()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		header,
		tabs,
		"",
		statsDetail,
		help,
	)
}

func (m Model) renderHeader() string {
	title := BrandStyle.Render("POLYTERM")
	subtitle := HeaderStyle.Render("Polymarket Analytics Platform")

	timeStr := ""
	if !m.lastUpdate.IsZero() {
		timeStr = MutedStyle.Render(fmt.Sprintf("Updated: %s", m.lastUpdate.Format("15:04:05")))
	}

	refreshStatus := ""
	if m.autoRefresh {
		refreshStatus = MutedStyle.Render("Auto: ON")
	} else {
		refreshStatus = MutedStyle.Render("Auto: OFF")
	}

	headerLine := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		" ",
		subtitle,
		"  ",
		timeStr,
		"  ",
		refreshStatus,
	)

	return headerLine
}

func (m Model) renderTabs() string {
	activeTab := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#4F46E5")).
		Padding(0, 2)

	inactiveTab := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#9CA3AF")).
		Background(lipgloss.Color("#1F2937")).
		Padding(0, 2)

	tab1 := inactiveTab.Render("[1] Markets")
	tab2 := inactiveTab.Render("[2] Analytics")

	switch m.currentPage {
	case pageMarkets:
		tab1 = activeTab.Render("[1] Markets")
	case pageStats:
		tab2 = activeTab.Render("[2] Analytics")
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, tab1, " ", tab2)
}

func (m Model) renderFilterBar() string {
	sortName := ""
	switch m.sortBy {
	case sortVolume:
		sortName = "Volume"
	case sortChange:
		sortName = "Change"
	case sortLiquidity:
		sortName = "Liquidity"
	}

	filterName := ""
	switch m.filterBy {
	case filterAll:
		filterName = "All"
	case filterCrypto:
		filterName = "Crypto"
	case filterPolitics:
		filterName = "Politics"
	case filterSports:
		filterName = "Sports"
	case filterEntertainment:
		filterName = "Entertainment"
	}

	sortStyle := lipgloss.NewStyle().Foreground(polyPurple).Bold(true)
	filterStyle := lipgloss.NewStyle().Foreground(polyBlue).Bold(true)
	searchStyle := lipgloss.NewStyle().Foreground(polyPink).Bold(true)

	displayLen := len(m.filteredMarkets)
	if displayLen == 0 && len(m.markets) > 0 {
		displayLen = len(m.markets)
	}

	var parts []string
	parts = append(parts, MutedStyle.Render("Sort: ")+sortStyle.Render(sortName))
	parts = append(parts, MutedStyle.Render("Filter: ")+filterStyle.Render(filterName))

	if m.searchMode {
		parts = append(parts, searchStyle.Render("Search: ")+m.searchQuery+"_")
	} else if m.searchQuery != "" {
		parts = append(parts, MutedStyle.Render("Search: ")+searchStyle.Render(m.searchQuery))
	} else {
		parts = append(parts, MutedStyle.Render("Search: ")+MutedStyle.Render("-"))
	}

	parts = append(parts, MutedStyle.Render(fmt.Sprintf("Results: %d", displayLen)))

	return lipgloss.JoinHorizontal(lipgloss.Top, parts[0], "  ", parts[1], "  ", parts[2], "  ", parts[3])
}

func (m Model) renderStats() string {
	compactStyle := lipgloss.NewStyle().
		Foreground(grayMuted).
		Bold(false)

	valueStyle := lipgloss.NewStyle().
		Foreground(polyBlue).
		Bold(true)

	parts := []string{
		compactStyle.Render("24h: ") + valueStyle.Render(formatCurrency(m.stats.Volume24h)),
		compactStyle.Render("Vol: ") + valueStyle.Render(formatCurrency(m.stats.TotalVolume)),
		compactStyle.Render("Markets: ") + valueStyle.Render(fmt.Sprintf("%d", m.stats.ActiveMarkets)),
		compactStyle.Render("Liq: ") + valueStyle.Render(formatCurrency(m.stats.AvgLiquidity)),
	}

	if m.stats.TopVolume != nil {
		parts = append(parts, compactStyle.Render("Hot: ")+lipgloss.NewStyle().Foreground(polyPink).Render(truncate(m.stats.TopVolume.Question, 40)))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, parts[0], "  ", parts[1], "  ", parts[2], "  ", parts[3], "  ", parts[4])
}

func (m Model) renderStatBox(label, value string) string {
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		StatsLabelStyle.Render(label),
		StatsValueStyle.Render(value),
	)
	return StatsBoxStyle.Render(content)
}

func (m Model) renderTable() string {
	displayMarkets := m.filteredMarkets
	if len(displayMarkets) == 0 {
		displayMarkets = m.markets
	}

	if len(displayMarkets) == 0 {
		return MutedStyle.Render("No markets available")
	}

	colWidths := []int{4, 55, 12, 12, 15, 10}

	headers := []string{"#", "Market", "Yes %", "No %", "Total Vol", "24h Vol"}
	headerRow := m.renderTableRow(headers, colWidths, TableHeaderStyle)

	var rows []string
	rows = append(rows, headerRow)

	start := m.scroll
	end := m.scroll + m.maxDisplay
	if end > len(displayMarkets) {
		end = len(displayMarkets)
	}

	for i := start; i < end; i++ {
		market := displayMarkets[i]

		yesOdds, noOdds := api.ParseOdds(&market)

		cells := []string{
			fmt.Sprintf("%d", i+1),
			truncate(market.Question, 53),
			fmt.Sprintf("%.1f%%", yesOdds),
			fmt.Sprintf("%.1f%%", noOdds),
			formatCurrency(market.GetVolume()),
			formatCurrency(market.Volume24hr),
		}

		rowStyle := TableCellStyle
		if i == m.cursor {
			rowStyle = rowStyle.Background(lipgloss.Color("#6366F1")).Bold(true)
		} else if i%2 == 0 {
			rowStyle = rowStyle.Background(lipgloss.Color("#1F2937"))
		}

		row := m.renderTableRow(cells, colWidths, rowStyle)
		rows = append(rows, row)
	}

	scrollInfo := ""
	if len(displayMarkets) > m.maxDisplay {
		scrollInfo = MutedStyle.Render(fmt.Sprintf(
			"Showing %d-%d of %d markets",
			start+1, end, len(displayMarkets),
		))
		rows = append(rows, scrollInfo)
	}

	return strings.Join(rows, "\n")
}

func (m Model) renderTableRow(cells []string, widths []int, style lipgloss.Style) string {
	var formatted []string
	for i, cell := range cells {
		width := widths[i]
		if len(cell) > width {
			cell = cell[:width]
		}

		cellStyle := style
		if i == 2 {
			cellStyle = YesOddsStyle
		} else if i == 3 {
			cellStyle = NoOddsStyle
		} else if i == 4 {
			cellStyle = VolumeStyle
		}

		padded := fmt.Sprintf("%-*s", width, cell)
		formatted = append(formatted, cellStyle.Render(padded))
	}
	return strings.Join(formatted, " ")
}

func (m Model) renderHelp() string {
	var helps []string
	if m.currentView == viewList {
		if m.currentPage == pageMarkets {
			if m.searchMode {
				helps = []string{
					"type to search",
					"enter/esc: exit search",
					"backspace: delete",
					"ctrl+u: clear",
				}
			} else {
				helps = []string{
					"↑/↓ j/k: nav",
					"enter: details",
					"/: search",
					"f: filter",
					"s: sort",
					"c: clear",
					"q: quit",
				}
			}
		} else {
			helps = []string{
				"1/2 or tab: switch page",
				"r: refresh",
				"a: auto-refresh",
				"q: quit",
			}
		}
	} else {
		helps = []string{
			"esc: back",
			"q: quit",
		}
	}
	return HelpStyle.Render(strings.Join(helps, " | "))
}

func (m Model) renderAdvancedStats() string {
	if len(m.markets) == 0 {
		return MutedStyle.Render("No data available")
	}

	var sections []string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#1F2937")).
		Padding(0, 2)

	labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#9CA3AF"))
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	sections = append(sections, titleStyle.Render("PLATFORM OVERVIEW"))

	platformStats := lipgloss.JoinVertical(
		lipgloss.Left,
		"",
		labelStyle.Render("Total Markets:     ")+valueStyle.Render(fmt.Sprintf("%d", m.stats.TotalMarkets)),
		labelStyle.Render("Active Markets:    ")+valueStyle.Render(fmt.Sprintf("%d", m.stats.ActiveMarkets)),
		labelStyle.Render("24h Volume:        ")+valueStyle.Render(formatCurrency(m.stats.Volume24h)),
		labelStyle.Render("Total Volume:      ")+valueStyle.Render(formatCurrency(m.stats.TotalVolume)),
		labelStyle.Render("Avg Liquidity:     ")+valueStyle.Render(formatCurrency(m.stats.AvgLiquidity)),
	)
	sections = append(sections, platformStats)

	topMarkets := getTopMarketsByVolume(m.markets, 10)
	if len(topMarkets) > 0 {
		sections = append(sections, "", "", titleStyle.Render("TOP 10 MARKETS BY TOTAL VOLUME"))
		sections = append(sections, "")

		for i, market := range topMarkets {
			yesOdds, _ := api.ParseOdds(&market)
			line := fmt.Sprintf("%2d. %-50s %8s  YES: %5.1f%%",
				i+1,
				truncate(market.Question, 48),
				formatCurrency(market.GetVolume()),
				yesOdds)
			sections = append(sections, lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).Render(line))
		}
	}

	topVolume24h := getTopMarketsByVolume24h(m.markets, 10)
	if len(topVolume24h) > 0 {
		sections = append(sections, "", "", titleStyle.Render("TOP 10 MARKETS BY 24H VOLUME"))
		sections = append(sections, "")

		for i, market := range topVolume24h {
			yesOdds, _ := api.ParseOdds(&market)
			line := fmt.Sprintf("%2d. %-50s %8s  YES: %5.1f%%",
				i+1,
				truncate(market.Question, 48),
				formatCurrency(market.Volume24hr),
				yesOdds)
			sections = append(sections, lipgloss.NewStyle().Foreground(lipgloss.Color("#D1D5DB")).Render(line))
		}
	}

	movers := getTopMovers(m.markets, 10)
	if len(movers) > 0 {
		sections = append(sections, "", "", titleStyle.Render("BIGGEST 24H PRICE MOVERS"))
		sections = append(sections, "")

		for i, market := range movers {
			changeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#10B981")).Bold(true)
			if market.OneDayPriceChange < 0 {
				changeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EF4444")).Bold(true)
			}

			line := fmt.Sprintf("%2d. %-50s %s",
				i+1,
				truncate(market.Question, 48),
				changeStyle.Render(fmt.Sprintf("%+6.2f%%", market.OneDayPriceChange*100)))
			sections = append(sections, line)
		}
	}

	momentum := getTopMomentum(m.markets, 10)
	if len(momentum) > 0 {
		sections = append(sections, "", "", titleStyle.Render("HIGHEST MOMENTUM (1H+1D+1W COMBINED)"))
		sections = append(sections, "")

		for i, market := range momentum {
			score := market.GetMomentumScore()
			line := fmt.Sprintf("%2d. %-50s Score: %5.1f",
				i+1,
				truncate(market.Question, 48),
				score*100)
			sections = append(sections, lipgloss.NewStyle().Foreground(lipgloss.Color("#A78BFA")).Render(line))
		}
	}

	engaged := getMostEngaged(m.markets, 10)
	if len(engaged) > 0 {
		sections = append(sections, "", "", titleStyle.Render("MOST ENGAGED (VOLUME + COMMENTS)"))
		sections = append(sections, "")

		for i, market := range engaged {
			line := fmt.Sprintf("%2d. %-50s Vol: %6s  Comments: %d",
				i+1,
				truncate(market.Question, 48),
				formatCurrency(market.Volume24hr),
				market.CommentCount)
			sections = append(sections, lipgloss.NewStyle().Foreground(lipgloss.Color("#FCD34D")).Render(line))
		}
	}

	tightSpreads := getTightestSpreads(m.markets, 10)
	if len(tightSpreads) > 0 {
		sections = append(sections, "", "", titleStyle.Render("TIGHTEST SPREADS (BEST LIQUIDITY)"))
		sections = append(sections, "")

		for i, market := range tightSpreads {
			spread := market.GetSpread()
			if spread > 0 {
				line := fmt.Sprintf("%2d. %-50s Spread: %.4f",
					i+1,
					truncate(market.Question, 48),
					spread)
				sections = append(sections, lipgloss.NewStyle().Foreground(lipgloss.Color("#34D399")).Render(line))
			}
		}
	}

	openInterest := getHighestOpenInterest(m.markets, 10)
	if len(openInterest) > 0 {
		sections = append(sections, "", "", titleStyle.Render("HIGHEST OPEN INTEREST (MONEY AT STAKE)"))
		sections = append(sections, "")

		for i, market := range openInterest {
			if market.OpenInterest > 0 {
				line := fmt.Sprintf("%2d. %-50s %8s",
					i+1,
					truncate(market.Question, 48),
					formatCurrency(market.OpenInterest))
				sections = append(sections, lipgloss.NewStyle().Foreground(lipgloss.Color("#F472B6")).Render(line))
			}
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func getTopMarketsByVolume(markets []types.Market, limit int) []types.Market {
	sorted := make([]types.Market, len(markets))
	copy(sorted, markets)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].GetVolume() > sorted[j].GetVolume()
	})

	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted
}

func getTopMarketsByVolume24h(markets []types.Market, limit int) []types.Market {
	sorted := make([]types.Market, len(markets))
	copy(sorted, markets)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Volume24hr > sorted[j].Volume24hr
	})

	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted
}

func getTopMovers(markets []types.Market, limit int) []types.Market {
	sorted := make([]types.Market, len(markets))
	copy(sorted, markets)

	sort.Slice(sorted, func(i, j int) bool {
		absI := sorted[i].OneDayPriceChange
		absJ := sorted[j].OneDayPriceChange
		if absI < 0 {
			absI = -absI
		}
		if absJ < 0 {
			absJ = -absJ
		}
		return absI > absJ
	})

	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted
}

func getTopMomentum(markets []types.Market, limit int) []types.Market {
	sorted := make([]types.Market, len(markets))
	copy(sorted, markets)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].GetMomentumScore() > sorted[j].GetMomentumScore()
	})

	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted
}

func getMostEngaged(markets []types.Market, limit int) []types.Market {
	sorted := make([]types.Market, len(markets))
	copy(sorted, markets)

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].GetEngagementScore() > sorted[j].GetEngagementScore()
	})

	if len(sorted) > limit {
		sorted = sorted[:limit]
	}
	return sorted
}

func getTightestSpreads(markets []types.Market, limit int) []types.Market {
	filtered := make([]types.Market, 0)
	for _, m := range markets {
		if m.GetSpread() > 0 {
			filtered = append(filtered, m)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].GetSpread() < filtered[j].GetSpread()
	})

	if len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered
}

func getHighestOpenInterest(markets []types.Market, limit int) []types.Market {
	filtered := make([]types.Market, 0)
	for _, m := range markets {
		if m.OpenInterest > 0 {
			filtered = append(filtered, m)
		}
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].OpenInterest > filtered[j].OpenInterest
	})

	if len(filtered) > limit {
		filtered = filtered[:limit]
	}
	return filtered
}

func (m Model) renderMarketDetail() string {
	displayMarkets := m.filteredMarkets
	if len(displayMarkets) == 0 {
		displayMarkets = m.markets
	}

	if m.selectedMarket >= len(displayMarkets) {
		return MutedStyle.Render("Market not found")
	}

	market := displayMarkets[m.selectedMarket]
	yesOdds, noOdds := api.ParseOdds(&market)

	var sections []string

	sections = append(sections, "", "")

	header := BrandStyle.Render("POLYTERM") + " " + HeaderStyle.Render(fmt.Sprintf("Market #%d Details", m.selectedMarket+1))
	sections = append(sections, header)
	sections = append(sections, "")

	questionStyle := lipgloss.NewStyle().
		Foreground(polyLight).
		Bold(true).
		Width(m.width - 6)
	sections = append(sections, questionStyle.Render(market.Question))
	sections = append(sections, "")

	probabilityBar := renderProbabilityBar(yesOdds, noOdds, m.width-10)
	sections = append(sections, probabilityBar)
	sections = append(sections, "")

	oddsSection := renderOddsBoxes(yesOdds, noOdds, market.OneDayPriceChange)
	sections = append(sections, oddsSection)
	sections = append(sections, "")

	volumeBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(polyPink).
		Padding(1, 2).
		Width(35).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(polyPink).Bold(true).Render("VOLUME & LIQUIDITY"),
			"",
			StatsLabelStyle.Render("24h Volume:")+lipgloss.NewStyle().Render("  ")+VolumeStyle.Render(formatCurrency(market.Volume24hr)),
			StatsLabelStyle.Render("Total Volume:")+lipgloss.NewStyle().Render(" ")+StatsValueStyle.Render(formatCurrency(market.GetVolume())),
			StatsLabelStyle.Render("Liquidity:")+lipgloss.NewStyle().Render("    ")+StatsValueStyle.Render(formatCurrency(market.GetLiquidity())),
		))

	priceBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(polyBlue).
		Padding(1, 2).
		Width(35).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().Foreground(polyBlue).Bold(true).Render("PRICE DATA"),
			"",
			StatsLabelStyle.Render("Last Price:")+lipgloss.NewStyle().Render("   ")+StatsValueStyle.Render(fmt.Sprintf("$%.3f", market.LastTradePrice)),
			StatsLabelStyle.Render("24h Change:")+lipgloss.NewStyle().Render("   ")+getPriceChangeStyle(market.OneDayPriceChange).Render(fmt.Sprintf("%+.2f%%", market.OneDayPriceChange*100)),
			StatsLabelStyle.Render("Status:")+lipgloss.NewStyle().Render("        ")+getActiveStyle(market.Active).Render(getStatusText(market.Active, market.Closed)),
		))

	statsRow := lipgloss.JoinHorizontal(lipgloss.Top, volumeBox, "  ", priceBox)
	sections = append(sections, statsRow)
	sections = append(sections, "")

	if market.Description != "" {
		descBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(polyPurple).
			Padding(1, 2).
			Width(m.width - 8).
			Render(lipgloss.JoinVertical(
				lipgloss.Left,
				lipgloss.NewStyle().Foreground(polyPurple).Bold(true).Render("MARKET DETAILS"),
				"",
				lipgloss.NewStyle().Foreground(polyLight).Render(truncateDesc(market.Description, m.width-16)),
			))
		sections = append(sections, descBox)
		sections = append(sections, "")
	}

	metaBox := lipgloss.NewStyle().
		Padding(0, 2).
		Render(lipgloss.JoinVertical(
			lipgloss.Left,
			MutedStyle.Render(fmt.Sprintf("Category: %s  |  Market ID: %s  |  Closes: %s",
				getCategory(market.Category),
				market.ID,
				formatEndDate(market.EndDate))),
		))
	sections = append(sections, metaBox)

	help := m.renderHelp()
	sections = append(sections, "", help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func renderProbabilityBar(yesOdds, noOdds float64, width int) string {
	barWidth := width - 20
	if barWidth < 20 {
		barWidth = 20
	}

	yesBlocks := int((yesOdds / 100.0) * float64(barWidth))
	noBlocks := barWidth - yesBlocks

	yesBar := strings.Repeat("█", yesBlocks)
	noBar := strings.Repeat("█", noBlocks)

	yesStyle := lipgloss.NewStyle().Foreground(greenYes)
	noStyle := lipgloss.NewStyle().Foreground(redNo)

	bar := yesStyle.Render(yesBar) + noStyle.Render(noBar)

	labelStyle := lipgloss.NewStyle().Bold(true)
	labels := lipgloss.JoinHorizontal(
		lipgloss.Top,
		labelStyle.Foreground(greenYes).Render(fmt.Sprintf("YES %.1f%%", yesOdds)),
		strings.Repeat(" ", barWidth-18),
		labelStyle.Foreground(redNo).Render(fmt.Sprintf("NO %.1f%%", noOdds)),
	)

	return lipgloss.JoinVertical(lipgloss.Left, labels, bar)
}

func renderOddsBoxes(yesOdds, noOdds, change float64) string {
	yesBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(greenYes).
		Padding(2, 4).
		Width(25).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(greenYes).Bold(true).Render("YES"),
			lipgloss.NewStyle().Foreground(greenYes).Bold(true).Render(fmt.Sprintf("%.1f%%", yesOdds)),
		))

	noBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(redNo).
		Padding(2, 4).
		Width(25).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(redNo).Bold(true).Render("NO"),
			lipgloss.NewStyle().Foreground(redNo).Bold(true).Render(fmt.Sprintf("%.1f%%", noOdds)),
		))

	changeBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(polyPurple).
		Padding(2, 4).
		Width(25).
		Align(lipgloss.Center).
		Render(lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().Foreground(polyPurple).Bold(true).Render("24H CHANGE"),
			getPriceChangeStyle(change).Render(fmt.Sprintf("%+.2f%%", change*100)),
		))

	return lipgloss.JoinHorizontal(lipgloss.Top, yesBox, "  ", noBox, "  ", changeBox)
}

func getStatusText(active, closed bool) string {
	if closed {
		return "Closed"
	}
	if active {
		return "Active"
	}
	return "Inactive"
}

func getCategory(cat string) string {
	if cat == "" {
		return "General"
	}
	return cat
}

func truncateDesc(desc string, maxWidth int) string {
	maxLen := maxWidth * 6
	if len(desc) > maxLen {
		return desc[:maxLen] + "..."
	}
	return desc
}

func getPriceChangeStyle(change float64) lipgloss.Style {
	if change > 0 {
		return lipgloss.NewStyle().Foreground(greenYes).Bold(true)
	} else if change < 0 {
		return lipgloss.NewStyle().Foreground(redNo).Bold(true)
	}
	return lipgloss.NewStyle().Foreground(grayMuted)
}

func getActiveStyle(active bool) lipgloss.Style {
	if active {
		return lipgloss.NewStyle().Foreground(greenYes).Bold(true)
	}
	return lipgloss.NewStyle().Foreground(grayMuted)
}

func formatEndDate(dateStr string) string {
	if dateStr == "" {
		return "N/A"
	}
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t.Format("Jan 02, 2006 15:04 MST")
	}
	return dateStr
}

func formatCurrency(amount float64) string {
	if amount >= 1000000 {
		return fmt.Sprintf("$%.2fM", amount/1000000)
	} else if amount >= 1000 {
		return fmt.Sprintf("$%.2fK", amount/1000)
	}
	return fmt.Sprintf("$%.2f", amount)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "N/A"
	}
	return t.Format("Jan 02 15:04")
}
