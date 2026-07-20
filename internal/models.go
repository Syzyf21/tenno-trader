package internal

import "time"

type BaroInventoryItem struct {
	Item       string `json:"item"`
	Ducats     int    `json:"ducats"`
	Credits    int    `json:"credits"`
	UniqueName string `json:"uniqueName"`
}

type MarketItem struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`
	I18n struct {
		En struct {
			Name string `json:"name"`
		} `json:"en"`
	} `json:"i18n"`
}

type StatEntry struct {
	Datetime string  `json:"datetime"`
	Volume   int     `json:"volume"`
	AvgPrice float64 `json:"avg_price"`
	ModRank  int     `json:"mod_rank"`
}

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

type AnalysisWindow struct {
	Start time.Time
	End   time.Time
}
