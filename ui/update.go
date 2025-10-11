package ui

import (
	"time"

	"polyterm/types"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.maxDisplay = (msg.Height - 12)
		if m.maxDisplay < 10 {
			m.maxDisplay = 10
		}
		if m.maxDisplay > 25 {
			m.maxDisplay = 25
		}
		if !m.ready {
			m.ready = true
		}
		return m, nil

	case types.FetchResult:
		m.loading = false
		m.err = msg.Err
		if msg.Err == nil {
			m.markets = msg.Markets
			m.stats = msg.Stats
			m.lastUpdate = time.Now()
			m.applyFiltersAndSort()
		}
		if !m.ready {
			m.ready = true
		}
		return m, nil

	case tickMsg:
		if m.autoRefresh && !m.loading {
			return m, tea.Batch(
				fetchMarketsCmd(500),
				tickCmd(),
			)
		}
		return m, tickCmd()

	case tea.KeyMsg:
		if m.searchMode {
			switch msg.String() {
			case "esc", "enter":
				m.searchMode = false
				m.applyFiltersAndSort()
				return m, nil
			case "backspace":
				if len(m.searchQuery) > 0 {
					m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
					m.applyFiltersAndSort()
				}
				return m, nil
			case "ctrl+u":
				m.searchQuery = ""
				m.applyFiltersAndSort()
				return m, nil
			default:
				if len(msg.String()) == 1 {
					m.searchQuery += msg.String()
					m.applyFiltersAndSort()
				}
				return m, nil
			}
		}
		
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		
		case "esc":
			if m.currentView == viewDetail {
				m.currentView = viewList
				m.selectedMarket = -1
				return m, nil
			}
			return m, tea.Quit
		
		case "r":
			if m.currentView == viewList {
				m.loading = true
				m.err = nil
				return m, tea.Batch(m.spinner.Tick, fetchMarketsCmd(500))
			}
			return m, nil
		
		case "a":
			m.autoRefresh = !m.autoRefresh
			return m, nil
		
		case "/":
			if m.currentView == viewList && m.currentPage == pageMarkets {
				m.searchMode = true
				return m, nil
			}
			return m, nil
		
		case "f":
			if m.currentView == viewList && m.currentPage == pageMarkets {
				m.filterBy = (m.filterBy + 1) % 5
				m.cursor = 0
				m.scroll = 0
				m.applyFiltersAndSort()
				return m, nil
			}
			return m, nil
		
		case "s":
			if m.currentView == viewList && m.currentPage == pageMarkets {
				m.sortBy = (m.sortBy + 1) % 3
				m.cursor = 0
				m.scroll = 0
				m.applyFiltersAndSort()
				return m, nil
			}
			return m, nil
		
		case "c":
			if m.currentView == viewList && m.currentPage == pageMarkets {
				m.searchQuery = ""
				m.filterBy = filterAll
				m.sortBy = sortVolume
				m.cursor = 0
				m.scroll = 0
				m.applyFiltersAndSort()
				return m, nil
			}
			return m, nil
		
		case "tab":
			if m.currentView == viewList {
				m.currentPage = (m.currentPage + 1) % 2
				return m, nil
			}
			return m, nil
		
		case "1":
			if m.currentView == viewList {
				m.currentPage = pageMarkets
				return m, nil
			}
			return m, nil
		
		case "2":
			if m.currentView == viewList {
				m.currentPage = pageStats
				return m, nil
			}
			return m, nil
		
		case "enter":
			if m.currentView == viewList && m.currentPage == pageMarkets && len(m.filteredMarkets) > 0 {
				m.selectedMarket = m.cursor
				if m.selectedMarket < len(m.filteredMarkets) {
					m.currentView = viewDetail
				}
				return m, nil
			}
			return m, nil
		
		case "up", "k":
			if m.currentView == viewList && m.currentPage == pageMarkets {
				if m.cursor > 0 {
					m.cursor--
					
					middle := m.maxDisplay / 2
					if m.cursor < m.scroll+middle && m.scroll > 0 {
						m.scroll--
					}
				}
			}
			return m, nil
		
		case "down", "j":
			if m.currentView == viewList && m.currentPage == pageMarkets {
				maxLen := len(m.filteredMarkets)
				if maxLen == 0 {
					maxLen = len(m.markets)
				}
				if m.cursor < maxLen-1 {
					m.cursor++
					
					middle := m.maxDisplay / 2
					maxScroll := maxLen - m.maxDisplay
					if maxScroll < 0 {
						maxScroll = 0
					}
					
					if m.cursor >= m.scroll+middle && m.scroll < maxScroll {
						m.scroll++
					}
				}
			}
			return m, nil
		
		case "home", "g":
			if m.currentView == viewList && m.currentPage == pageMarkets {
				m.cursor = 0
				m.scroll = 0
			}
			return m, nil
		
		case "end", "G":
			if m.currentView == viewList && m.currentPage == pageMarkets {
				maxLen := len(m.filteredMarkets)
				if maxLen == 0 {
					maxLen = len(m.markets)
				}
				m.cursor = maxLen - 1
				maxScroll := maxLen - m.maxDisplay
				if maxScroll < 0 {
					maxScroll = 0
				}
				m.scroll = maxScroll
			}
			return m, nil
		
		case "pageup":
			if m.currentView == viewList && m.currentPage == pageMarkets {
				m.cursor -= m.maxDisplay
				if m.cursor < 0 {
					m.cursor = 0
				}
				m.scroll -= m.maxDisplay
				if m.scroll < 0 {
					m.scroll = 0
				}
			}
			return m, nil
		
		case "pagedown":
			if m.currentView == viewList && m.currentPage == pageMarkets {
				maxLen := len(m.filteredMarkets)
				if maxLen == 0 {
					maxLen = len(m.markets)
				}
				maxCursor := maxLen - 1
				m.cursor += m.maxDisplay
				if m.cursor > maxCursor {
					m.cursor = maxCursor
				}
				
				maxScroll := maxLen - m.maxDisplay
				if maxScroll < 0 {
					maxScroll = 0
				}
				m.scroll += m.maxDisplay
				if m.scroll > maxScroll {
					m.scroll = maxScroll
				}
			}
			return m, nil
		}
	}

	if m.loading {
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}
