package internal

import "time"

// BaroInventoryItem is one entry sold by Baro Ki'Teer, as reported by the
// warframestat.us worldstate API.
type BaroInventoryItem struct {
	Item       string `json:"item"`
	Ducats     int    `json:"ducats"`
	Credits    int    `json:"credits"`
	UniqueName string `json:"uniqueName"`
}

// MarketItem is one entry from warframe.market's /v2/items list.
type MarketItem struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	I18n struct {
		En struct {
			Name string `json:"name"`
		} `json:"en"`
	} `json:"i18n"`
}

// StatEntry is one daily (or hourly) bucket from warframe.market statistics.
type StatEntry struct {
	Datetime string  `json:"datetime"`
	Volume   int     `json:"volume"`
	AvgPrice float64 `json:"avg_price"`
	ModRank  int     `json:"mod_rank"`
}

// VoidTraderRow is a fully computed line of the results table shown in the UI.
type VoidTraderRow struct {
	Name         string
	Ducats       int
	Credits      int
	AvgPlatinum  float64
	AvgVolume    float64
	PlatPerDucat float64
	DataPoints   int
	NoMarketData bool
}

type ArbitrationRow struct {
	Name         string
	Vitus        int
	AvgPlatinum  float64
	AvgVolume    float64
	PlatPerVitus float64
	NoMarketData bool
}

type ArbitrationItem struct {
	ID    string
	Name  string
	Vitus int
}

// AnalysisWindow is the date range (inclusive, whole days, UTC) over which
// price/volume statistics are averaged.
type AnalysisWindow struct {
	Start time.Time
	End   time.Time
}
