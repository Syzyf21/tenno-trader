package baroinvestor

import (
	"fmt"
	"time"

	"github.com/Syzyf21/tenno-trader/internal"
)

type ProgressFunc func(done, total int, currentItem string)

type StatusFunc func(text string)

func BuildRows(onStatus StatusFunc, onProgress ProgressFunc, isLive bool) ([]internal.VoidTraderRow, internal.AnalysisWindow, error) {
	notify := func(s string) {
		if onStatus != nil {
			onStatus(s)
		}
	}

	var inv *VoidTraderResponse
	var invDB []internal.BaroInventoryItemDB
	var err error
	if isLive {
		notify("Fetching Baro Ki'Teer inventory...")
		inv, err = fetchBaroInventory()
		if err != nil {
			return nil, internal.AnalysisWindow{}, fmt.Errorf("Error while fetching Baro inventory: %w", err)
		}
	} else {
		notify("Fetching Baro Ki'Teer inventory from database...")
		invDB, err = GetVoidTraderItems(internal.DBConn)
		if err != nil {
			return nil, internal.AnalysisWindow{}, fmt.Errorf("Error while fetching Baro inventory from database: %w", err)
		}
	}

	var timeWindow internal.AnalysisWindow
	if isLive {
		var refDate time.Time
		refDate, err = inv.activationDate()
		if err != nil {
			return nil, internal.AnalysisWindow{}, fmt.Errorf("Error while determining baro activation date: %w", err)
		}

		timeWindow = internal.AnalysisWindow{
			Start: refDate.AddDate(0, 0, -14),
			End:   refDate.AddDate(0, 0, -1),
		}
	} else {
		timeWindow = internal.AnalysisWindow{
			Start: time.Now().AddDate(0, 0, -14),
			End:   time.Now(),
		}
	}

	notify("Fetching warframe.market item catalogue...")
	index, err := internal.FetchMarketItemIndex()
	if err != nil {
		return nil, internal.AnalysisWindow{}, fmt.Errorf("Error while fetching warframe.market item catalogue: %w", err)
	}

	if isLive {
		total := len(inv.Inventory)
		rows := make([]internal.VoidTraderRow, 0, total)

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

			row := internal.VoidTraderRow{
				Name:   baroItem.Item,
				Ducats: baroItem.Ducats,
			}

			avgPrice, avgVolume, count := internal.AverageInWindow(entries, timeWindow, false)
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
	} else {
		total := len(invDB)
		rows := make([]internal.VoidTraderRow, 0, total)

		for i, baroItem := range invDB {
			if onProgress != nil {
				onProgress(i, total, baroItem.Name)
			}

			urlName, ok := index[internal.NormalizeItemName(baroItem.Name)]
			if !ok {
				continue
			}

			entries, err := internal.FetchItemStatistics(urlName)
			time.Sleep(internal.MarketRequestPause)

			if err != nil {
				continue
			}

			row := internal.VoidTraderRow{
				Name:   baroItem.Name,
				Ducats: baroItem.Ducats,
			}

			avgPrice, avgVolume, count := internal.AverageInWindow(entries, timeWindow, false)
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
}
