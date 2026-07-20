package arbitrations

import (
	"fmt"
	"time"

	"github.com/Syzyf21/tenno-trader/internal"
)

type ProgressFunc func(done, total int, currentItem string)

type StatusFunc func(text string)

func BuildRows(onStatus StatusFunc, onProgress ProgressFunc) ([]internal.ArbitrationRow, internal.AnalysisWindow, error) {
	notify := func(s string) {
		if onStatus != nil {
			onStatus(s)
		}
	}

	notify("Fetching arbitration items from database...")
	items, err := GetArbitrationItems(internal.DBConn)
	if err != nil {
		return nil, internal.AnalysisWindow{}, fmt.Errorf("Error while fetching arbitration items: %w", err)
	}

	timeWindow := internal.AnalysisWindow{
		Start: time.Now().AddDate(0, 0, -14),
		End:   time.Now(),
	}

	notify("Fetching warframe.market item catalogue...")
	index, err := internal.FetchMarketItemIndex()
	if err != nil {
		return nil, internal.AnalysisWindow{}, fmt.Errorf("Error while fetching warframe.market item catalogue: %w", err)
	}

	total := len(items)
	rows := make([]internal.ArbitrationRow, 0, total)

	for i, arbitrationItem := range items {
		if onProgress != nil {
			onProgress(i, total, arbitrationItem.Name)
		}

		urlName, ok := index[internal.NormalizeItemName(arbitrationItem.Name)]
		if !ok {
			continue
		}

		entries, err := internal.FetchItemStatistics(urlName)
		time.Sleep(internal.MarketRequestPause)

		if err != nil {
			continue
		}

		row := internal.ArbitrationRow{
			Name:  arbitrationItem.Name,
			Vitus: arbitrationItem.Vitus,
		}

		avgPrice, avgVolume, count := internal.AverageInWindow(entries, timeWindow)
		row.AvgPlatinum = avgPrice
		row.AvgVolume = avgVolume
		if row.Vitus > 0 && count > 0 {
			row.PlatPerVitus = avgPrice / float64(row.Vitus)
		}
		if count == 0 {
			row.NoMarketData = true
		}
		rows = append(rows, row)
	}

	if onProgress != nil {
		onProgress(total, total, "")
	}
	notify("Done.")

	return rows, timeWindow, nil
}
