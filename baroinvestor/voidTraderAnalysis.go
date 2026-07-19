package baroinvestor

import (
	"fmt"
	"time"

	"github.com/Syzyf21/tenno-trader/internal"
)

type ProgressFunc func(done, total int, currentItem string)

type StatusFunc func(text string)

func BuildRows(onStatus StatusFunc, onProgress ProgressFunc) ([]internal.Row, internal.AnalysisWindow, error) {
	notify := func(s string) {
		if onStatus != nil {
			onStatus(s)
		}
	}

	notify("Fetching Baro Ki'Teer inventory...")
	inv, err := fetchBaroInventory()
	if err != nil {
		return nil, internal.AnalysisWindow{}, fmt.Errorf("Error while fetching Baro inventory: %w", err)
	}

	var refDate time.Time
	refDate, err = inv.activationDate()
	if err != nil {
		return nil, internal.AnalysisWindow{}, fmt.Errorf("Error while determining baro activation date: %w", err)
	}

	timeWindow := internal.AnalysisWindow{
		Start: refDate.AddDate(0, 0, -10),
		End:   refDate.AddDate(0, 0, -1),
	}

	notify("Fetching warframe.market item catalogue...")
	index, err := internal.FetchMarketItemIndex()
	if err != nil {
		return nil, internal.AnalysisWindow{}, fmt.Errorf("Error while fetching warframe.market item catalogue: %w", err)
	}

	total := len(inv.Inventory)
	rows := make([]internal.Row, 0, total)

	for i, baroItem := range inv.Inventory {
		if onProgress != nil {
			onProgress(i, total, baroItem.Item)
		}

		urlName, ok := index[internal.NormalizeItemName(baroItem.Item)]
		if !ok {
			continue
		}

		entries, err := internal.FetchItemStatistics(urlName)
		time.Sleep(internal.MarketRequestPause)

		if err != nil {
			continue
		}

		row := internal.Row{
			Name:    baroItem.Item,
			Ducats:  baroItem.Ducats,
			Credits: baroItem.Credits,
		}

		avgPrice, avgVolume, count := internal.AverageInWindow(entries, timeWindow)
		row.AvgPlatinum = avgPrice
		row.AvgVolume = avgVolume
		row.DataPoints = count
		if row.Ducats > 0 && count > 0 {
			row.PlatPerDucat = avgPrice / float64(row.Ducats)
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
