package types

import (
	"encoding/json"
	"strconv"
)

type Market struct {
	ID                string  `json:"id"`
	Question          string  `json:"question"`
	Description       string  `json:"description"`
	VolumeStr         string  `json:"volume"`
	Volume24hr        float64 `json:"volume24hr"`
	LiquidityStr      string  `json:"liquidity"`
	VolumeNum         float64 `json:"volumeNum"`
	LiquidityNum      float64 `json:"liquidityNum"`
	EndDate           string  `json:"endDate"`
	Active            bool    `json:"active"`
	Closed            bool    `json:"closed"`
	MarketSlug        string  `json:"marketSlug"`
	OutcomesStr       string  `json:"outcomes"`
	OutcomePricesStr  string  `json:"outcomePrices"`
	CloseTime         string  `json:"closeTime"`
	Category          string  `json:"category"`
	LastTradePrice    float64 `json:"lastTradePrice"`
	OneDayPriceChange float64 `json:"oneDayPriceChange"`
}

func (m *Market) GetVolume() float64 {
	if m.VolumeNum > 0 {
		return m.VolumeNum
	}
	if m.VolumeStr != "" {
		if v, err := strconv.ParseFloat(m.VolumeStr, 64); err == nil {
			return v
		}
	}
	return 0
}

func (m *Market) GetLiquidity() float64 {
	if m.LiquidityNum > 0 {
		return m.LiquidityNum
	}
	if m.LiquidityStr != "" {
		if l, err := strconv.ParseFloat(m.LiquidityStr, 64); err == nil {
			return l
		}
	}
	return 0
}

func (m *Market) GetOutcomePrices() []string {
	var prices []string
	if m.OutcomePricesStr != "" {
		json.Unmarshal([]byte(m.OutcomePricesStr), &prices)
	}
	return prices
}

func (m *Market) GetOutcomes() []string {
	var outcomes []string
	if m.OutcomesStr != "" {
		json.Unmarshal([]byte(m.OutcomesStr), &outcomes)
	}
	return outcomes
}

type MarketsResponse struct {
	Data    []Market `json:"data"`
	Markets []Market `json:"markets"`
	Items   []Market `json:"items"`
}

type GlobalStats struct {
	TotalVolume     float64
	Volume24h       float64
	ActiveMarkets   int
	TotalMarkets    int
	AvgLiquidity    float64
	TopGainer       *Market
	TopGainerChange float64
	TopVolume       *Market
}

type FetchResult struct {
	Markets     []Market
	Stats       GlobalStats
	Err         error
}
