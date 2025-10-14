package types

import (
	"encoding/json"
	"strconv"
)

type Market struct {
	ID                  string  `json:"id"`
	Question            string  `json:"question"`
	Description         string  `json:"description"`
	VolumeStr           string  `json:"volume"`
	Volume24hr          float64 `json:"volume24hr"`
	LiquidityStr        string  `json:"liquidity"`
	VolumeNum           float64 `json:"volumeNum"`
	LiquidityNum        float64 `json:"liquidityNum"`
	EndDate             string  `json:"endDate"`
	Active              bool    `json:"active"`
	Closed              bool    `json:"closed"`
	MarketSlug          string  `json:"marketSlug"`
	OutcomesStr         string  `json:"outcomes"`
	OutcomePricesStr    string  `json:"outcomePrices"`
	CloseTime           string  `json:"closeTime"`
	Category            string  `json:"category"`
	LastTradePrice      float64 `json:"lastTradePrice"`
	OneDayPriceChange   float64 `json:"oneDayPriceChange"`
	OneHourPriceChange  float64 `json:"oneHourPriceChange"`
	OneWeekPriceChange  float64 `json:"oneWeekPriceChange"`
	OneMonthPriceChange float64 `json:"oneMonthPriceChange"`
	Volume1wk           float64 `json:"volume1wk"`
	Volume1mo           float64 `json:"volume1mo"`
	VolumeAmm           float64 `json:"volumeAmm"`
	VolumeClob          float64 `json:"volumeClob"`
	LiquidityAmm        float64 `json:"liquidityAmm"`
	LiquidityClob       float64 `json:"liquidityClob"`
	BestAsk             float64 `json:"bestAsk"`
	BestBid             float64 `json:"bestBid"`
	Spread              float64 `json:"spread"`
	CommentCount        int     `json:"commentCount"`
	OpenInterest        float64 `json:"openInterest"`
	Featured            bool    `json:"featured"`
	Competitive         float64 `json:"competitive"`
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

func (m *Market) GetSpread() float64 {
	if m.Spread > 0 {
		return m.Spread
	}
	if m.BestAsk > 0 && m.BestBid > 0 {
		return m.BestAsk - m.BestBid
	}
	return 0
}

func (m *Market) GetMomentumScore() float64 {
	score := 0.0
	if m.OneHourPriceChange != 0 {
		if m.OneHourPriceChange < 0 {
			score += -m.OneHourPriceChange * 3
		} else {
			score += m.OneHourPriceChange * 3
		}
	}
	if m.OneDayPriceChange != 0 {
		if m.OneDayPriceChange < 0 {
			score += -m.OneDayPriceChange * 2
		} else {
			score += m.OneDayPriceChange * 2
		}
	}
	if m.OneWeekPriceChange != 0 {
		if m.OneWeekPriceChange < 0 {
			score += -m.OneWeekPriceChange
		} else {
			score += m.OneWeekPriceChange
		}
	}
	return score
}

func (m *Market) GetEngagementScore() float64 {
	volScore := m.Volume24hr / 1000.0
	commentScore := float64(m.CommentCount) * 10.0
	return volScore + commentScore
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
	Markets []Market
	Stats   GlobalStats
	Err     error
}
