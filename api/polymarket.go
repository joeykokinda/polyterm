package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"time"

	"polyterm/types"
)

const (
	BaseURL = "https://gamma-api.polymarket.com"
	Timeout = 10 * time.Second
)

func FetchMarkets(ctx context.Context, limit int) ([]types.Market, types.GlobalStats, error) {
	ctx, cancel := context.WithTimeout(ctx, Timeout)
	defer cancel()

	url := fmt.Sprintf("%s/markets?closed=false&limit=%d", BaseURL, limit)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, types.GlobalStats{}, err
	}
	req.Header.Set("User-Agent", "polyterm/1.0.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, types.GlobalStats{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, types.GlobalStats{}, err
	}

	var rawData interface{}
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, types.GlobalStats{}, fmt.Errorf("invalid JSON: %w", err)
	}

	var markets []types.Market
	
	if err := json.Unmarshal(body, &markets); err == nil {
		if len(markets) == 0 {
			return nil, types.GlobalStats{}, fmt.Errorf("no markets returned")
		}
	} else {
		var marketsResp types.MarketsResponse
		if err := json.Unmarshal(body, &marketsResp); err != nil {
			return nil, types.GlobalStats{}, fmt.Errorf("failed to parse as array or object - first 500 chars of response: %s", truncateString(string(body), 500))
		}

		switch {
		case len(marketsResp.Data) > 0:
			markets = marketsResp.Data
		case len(marketsResp.Markets) > 0:
			markets = marketsResp.Markets
		case len(marketsResp.Items) > 0:
			markets = marketsResp.Items
		default:
			return nil, types.GlobalStats{}, fmt.Errorf("no markets in response")
		}
	}

	if len(markets) == 0 {
		return nil, types.GlobalStats{}, fmt.Errorf("no markets returned")
	}

	activeMarkets := make([]types.Market, 0)
	for _, m := range markets {
		vol := m.GetVolume()
		if vol > 100 || m.Volume24hr > 10 {
			activeMarkets = append(activeMarkets, m)
		}
	}

	sort.Slice(activeMarkets, func(i, j int) bool {
		volI := activeMarkets[i].GetVolume()
		volJ := activeMarkets[j].GetVolume()
		if volI == 0 && volJ == 0 {
			return activeMarkets[i].Volume24hr > activeMarkets[j].Volume24hr
		}
		return volI > volJ
	})
	
	if len(activeMarkets) > limit {
		activeMarkets = activeMarkets[:limit]
	}

	stats := calculateStats(activeMarkets)
	return activeMarkets, stats, nil
}

func calculateStats(markets []types.Market) types.GlobalStats {
	stats := types.GlobalStats{}
	
	activeCount := 0
	totalVolume := 0.0
	volume24h := 0.0
	totalLiquidity := 0.0
	
	var topGainer *types.Market
	topGainerChange := 0.0
	var topVolume *types.Market
	maxVolume24h := 0.0
	
	for i, m := range markets {
		if m.Active {
			activeCount++
		}
		totalVolume += m.GetVolume()
		volume24h += m.Volume24hr
		totalLiquidity += m.GetLiquidity()
		
		if m.OneDayPriceChange > topGainerChange {
			topGainerChange = m.OneDayPriceChange
			topGainer = &markets[i]
		}
		
		if m.Volume24hr > maxVolume24h {
			maxVolume24h = m.Volume24hr
			topVolume = &markets[i]
		}
	}
	
	stats.ActiveMarkets = activeCount
	stats.TotalMarkets = len(markets)
	stats.TotalVolume = totalVolume
	stats.Volume24h = volume24h
	stats.TopGainer = topGainer
	stats.TopGainerChange = topGainerChange
	stats.TopVolume = topVolume
	
	if len(markets) > 0 {
		stats.AvgLiquidity = totalLiquidity / float64(len(markets))
	}
	
	return stats
}

func ParseOdds(m *types.Market) (yesOdds, noOdds float64) {
	prices := m.GetOutcomePrices()
	if len(prices) >= 2 {
		if yes, err := strconv.ParseFloat(prices[0], 64); err == nil {
			yesOdds = yes * 100
		}
		if no, err := strconv.ParseFloat(prices[1], 64); err == nil {
			noOdds = no * 100
		}
	}
	return
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
