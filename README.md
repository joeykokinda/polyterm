# Polyterm

ssh polyterm@polyterm.app

A beautiful terminal-based analytics dashboard for Polymarket built with Go and Bubbletea.

## Features

- **500+ Markets** - Fetches comprehensive market data from Polymarket, filters by active volume
- **Search & Filter** - Real-time search and filter by category (Crypto, Politics, Sports, Entertainment)
- **Multiple Sort Options** - Sort by Volume, Price Change, or Liquidity
- **Multi-Page Interface** - Switch between Markets, Stats, and About pages
- **Real-time Data** - Live market data from Polymarket API with auto-refresh every 30 seconds
- **Enhanced Market Details** - Beautiful detail view with:
  - Visual probability bar chart
  - Large YES/NO odds display boxes
  - 24h price change indicator
  - Volume and liquidity metrics
  - Full market description
- **Advanced Stats** - Dedicated stats page with platform metrics, hottest markets, and biggest movers
- **Center Scrolling** - Selected market stays centered as you navigate
- **Beautiful UI** - Color-coded interface matching Polymarket's brand (blues, purples, pinks)
- **Keyboard Navigation** - Full keyboard control with vim-style bindings
- **Responsive Layout** - Adapts to your terminal size
- **Production Ready** - Fixed all rendering issues and navigation bugs

## Installation

```bash
go build -o polyterm
./polyterm
```

Or run directly:

```bash
go run main.go
```

## Usage

### Keyboard Controls

#### Markets Page
- `↑/↓` or `j/k` - Navigate markets (cursor stays centered)
- `g/G` - Jump to top/bottom
- `PgUp/PgDn` - Page up/down
- `Enter` - View detailed market information
- `/` - Enter search mode (type to search markets)
- `f` - Cycle through filters (All/Crypto/Politics/Sports/Entertainment)
- `s` - Cycle through sort options (Volume/Change/Liquidity)
- `c` - Clear all filters and search
- `1/2/3` or `Tab` - Switch between pages
- `r` - Manual refresh
- `a` - Toggle auto-refresh on/off
- `q` or `Ctrl+C` - Quit

#### Search Mode
- Type any text to search market titles in real-time
- `Backspace` - Delete last character
- `Ctrl+U` - Clear entire search
- `Enter` or `Esc` - Exit search mode

#### Stats/About Pages
- `1/2/3` or `Tab` - Switch between pages
- `r` - Manual refresh
- `a` - Toggle auto-refresh on/off
- `q` or `Ctrl+C` - Quit

#### Detail View
- `Esc` - Return to market list
- `q` or `Ctrl+C` - Quit

## Pages & Features

### Page 1: Markets
Browse and navigate trending prediction markets.

**Components:**
- **Header** - Branding, last update time, auto-refresh status
- **Tab Navigation** - Quick access to Markets, Stats, About
- **Stats Overview** - 24h Volume, Total Volume, Active Markets, Avg Liquidity, Hottest Market, Biggest Mover
- **Market Table** - Top 150 markets sorted by 24h volume
  - Market rank and question
  - Yes/No odds (color-coded green/red)
  - 24-hour volume (pink highlight)
  - Current liquidity
  - Centered cursor selection (highlighted in blue)

**Market Detail View** (Press `Enter`):
- Full market question displayed prominently
- **Visual probability bar** - Color-coded YES/NO distribution chart
- **Large odds display boxes** - YES, NO, and 24H CHANGE in dedicated boxes
- **Volume & Liquidity section** - 24h volume, total volume, liquidity
- **Price data section** - Last price, 24h change, market status
- **Full market description** - Complete details about resolution criteria
- **Market metadata** - Category, Market ID, closing date

**Search & Filter**:
- Press `/` to search for markets (e.g., "nyc mayor", "bitcoin", "election")
- Press `f` to filter by category:
  - All - Show all markets
  - Crypto - Bitcoin, Ethereum, crypto-related
  - Politics - Elections, presidents, congress
  - Sports - NBA, NFL, FIFA, championships
  - Entertainment - Movies, box office, Oscars
- Press `s` to sort by:
  - Volume - Highest 24h trading volume first (default)
  - Change - Biggest price movers first
  - Liquidity - Most liquid markets first
- Press `c` to clear all filters and reset

### Page 2: Stats
Detailed platform statistics and insights.

**Components:**
- **Platform Statistics** - Total/active markets, volumes, liquidity
- **Top Market** - Highest 24h volume market with full details
- **Biggest Mover** - Highest 24h price change with details

### Page 3: About
Information about Polyterm, features, tech stack, and open source details.

## API

Uses Polymarket's public Gamma API:
- Endpoint: `https://gamma-api.polymarket.com`
- Fetches top 150 active markets
- Client-side sorting by 24h volume for trending markets
- Auto-refreshes every 30 seconds (can be toggled off)

## Tech Stack

- [Go](https://golang.org/) - Language
- [Bubbletea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Bubbles](https://github.com/charmbracelet/bubbles) - UI components

## Color Scheme

Inspired by Polymarket's brand:
- Primary: Indigo Blue (#6366F1)
- Secondary: Purple (#8B5CF6)
- Accent: Pink (#EC4899)
- Success: Green (#10B981)
- Error: Red (#EF4444)

## License

MIT

## Contributing

This is an open-source project. Contributions welcome!
