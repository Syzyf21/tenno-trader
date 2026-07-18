package main

import "time"

// BaroInventoryItem is one entry sold by Baro Ki'Teer, as reported by the
// warframestat.us worldstate API.
type BaroInventoryItem struct {
	Item       string `json:"item"`
	Ducats     int    `json:"ducats"`
	Credits    int    `json:"credits"`
	UniqueName string `json:"uniqueName"`
}

type VoidTraderResponse struct {
	Activation string              `json:"activation"`
	Expiry     string              `json:"expiry"`
	Character  string              `json:"character"`
	Location   string              `json:"location"`
	Inventory  []BaroInventoryItem `json:"inventory"`
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

// Row is a fully computed line of the results table shown in the UI.
type Row struct {
	Name         string
	Ducats       int
	Credits      int
	AvgPlatinum  float64
	AvgVolume    float64
	PlatPerDucat float64
	DataPoints   int
	NoMarketData bool
}

// AnalysisWindow is the date range (inclusive, whole days, UTC) over which
// price/volume statistics are averaged.
type AnalysisWindow struct {
	Start time.Time
	End   time.Time
}
