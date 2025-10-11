package ui

import (
	"context"
	"sort"
	"strings"
	"time"

	"polyterm/api"
	"polyterm/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/spinner"
)

type tickMsg time.Time

type viewMode int

const (
	viewList viewMode = iota
	viewDetail
	viewStats
)

type pageMode int

const (
	pageMarkets pageMode = iota
	pageStats
)

type sortMode int

const (
	sortVolume sortMode = iota
	sortChange
	sortLiquidity
)

type filterMode int

const (
	filterAll filterMode = iota
	filterCrypto
	filterPolitics
	filterSports
	filterEntertainment
)

type Model struct {
	spinner         spinner.Model
	loading         bool
	markets         []types.Market
	filteredMarkets []types.Market
	stats           types.GlobalStats
	err             error
	width           int
	height          int
	scroll          int
	cursor          int
	maxDisplay      int
	lastUpdate      time.Time
	autoRefresh     bool
	currentView     viewMode
	currentPage     pageMode
	selectedMarket  int
	ready           bool
	searchMode      bool
	searchQuery     string
	sortBy          sortMode
	filterBy        filterMode
}

func NewModel() Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = LoadingStyle
	return Model{
		spinner:         s,
		loading:         true,
		scroll:          0,
		cursor:          0,
		maxDisplay:      20,
		autoRefresh:     true,
		currentView:     viewList,
		currentPage:     pageMarkets,
		selectedMarket:  -1,
		ready:           false,
		searchMode:      false,
		searchQuery:     "",
		sortBy:          sortVolume,
		filterBy:        filterAll,
		filteredMarkets: []types.Market{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchMarketsCmd(500),
		tickCmd(),
	)
}

func fetchMarketsCmd(limit int) tea.Cmd {
	return func() tea.Msg {
		markets, stats, err := api.FetchMarkets(context.Background(), limit)
		return types.FetchResult{Markets: markets, Stats: stats, Err: err}
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*30, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m *Model) applyFiltersAndSort() {
	filtered := make([]types.Market, 0)
	
	for _, market := range m.markets {
		vol := market.GetVolume()
		if vol < 100 && market.Volume24hr < 10 {
			continue
		}
		
		if m.filterBy != filterAll {
			category := strings.ToLower(market.Category)
			question := strings.ToLower(market.Question)
			
			match := false
			switch m.filterBy {
			case filterCrypto:
				match = strings.Contains(category, "crypto") || 
					strings.Contains(question, "bitcoin") ||
					strings.Contains(question, "ethereum") ||
					strings.Contains(question, "crypto")
			case filterPolitics:
				match = strings.Contains(category, "politics") ||
					strings.Contains(question, "election") ||
					strings.Contains(question, "president") ||
					strings.Contains(question, "congress") ||
					strings.Contains(question, "senate") ||
					strings.Contains(question, "mayor") ||
					strings.Contains(question, "governor")
			case filterSports:
				match = strings.Contains(category, "sports") ||
					strings.Contains(question, "nba") ||
					strings.Contains(question, "nfl") ||
					strings.Contains(question, "fifa") ||
					strings.Contains(question, "champion") ||
					strings.Contains(question, "world series") ||
					strings.Contains(question, "playoff")
			case filterEntertainment:
				match = strings.Contains(category, "entertainment") ||
					strings.Contains(question, "movie") ||
					strings.Contains(question, "oscar") ||
					strings.Contains(question, "box office")
			}
			
			if !match {
				continue
			}
		}
		
		if m.searchQuery != "" {
			query := strings.ToLower(m.searchQuery)
			question := strings.ToLower(market.Question)
			description := strings.ToLower(market.Description)
			
			if !strings.Contains(question, query) && !strings.Contains(description, query) {
				continue
			}
		}
		
		filtered = append(filtered, market)
	}
	
	switch m.sortBy {
	case sortVolume:
		sort.Slice(filtered, func(i, j int) bool {
			volI := filtered[i].GetVolume()
			volJ := filtered[j].GetVolume()
			if volI == 0 && volJ == 0 {
				return filtered[i].Volume24hr > filtered[j].Volume24hr
			}
			return volI > volJ
		})
	case sortChange:
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].OneDayPriceChange > filtered[j].OneDayPriceChange
		})
	case sortLiquidity:
		sort.Slice(filtered, func(i, j int) bool {
			return filtered[i].GetLiquidity() > filtered[j].GetLiquidity()
		})
	}
	
	m.filteredMarkets = filtered
	
	if m.cursor >= len(filtered) {
		m.cursor = len(filtered) - 1
		if m.cursor < 0 {
			m.cursor = 0
		}
	}
	if m.scroll >= len(filtered) {
		m.scroll = len(filtered) - m.maxDisplay
		if m.scroll < 0 {
			m.scroll = 0
		}
	}
}
