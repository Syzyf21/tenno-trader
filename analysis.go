package main

import (
	"fmt"
	"time"
)

type ProgressFunc func(done, total int, currentItem string)

type StatusFunc func(text string)

// buildRows fetches Baro's current inventory, matches each item against the
// warframe.market catalogue, and computes the 10-day price/volume averages
// for the window ending on Baro's arrival date.
//
// overrideActivation, if non-zero, replaces the arrival date reported by the
// worldstate API — useful for reproducing a specific historical stock such
// as the one on 2026-07-10.
func buildRows(onStatus StatusFunc, onProgress ProgressFunc) ([]Row, AnalysisWindow, error) {
	notify := func(s string) {
		if onStatus != nil {
			onStatus(s)
		}
	}

	notify("Fetching Baro Ki'Teer inventory...")
	inv, err := fetchBaroInventory()
	if err != nil {
		return nil, AnalysisWindow{}, fmt.Errorf("Error while fetching Baro inventory: %w", err)
	}

	var refDate time.Time
	refDate, err = inv.activationDate()
	if err != nil {
		return nil, AnalysisWindow{}, fmt.Errorf("Error while determining baro activation date: %w", err)
	}

	timeWindow := AnalysisWindow{
		Start: refDate.AddDate(0, 0, -10),
		End:   refDate.AddDate(0, 0, -1),
	}

	notify("Fetching warframe.market item catalogue...")
	index, err := fetchMarketItemIndex()
	if err != nil {
		return nil, AnalysisWindow{}, fmt.Errorf("Error while fetching warframe.market item catalogue: %w", err)
	}

	total := len(inv.Inventory)
	rows := make([]Row, 0, total)

	for i, baroItem := range inv.Inventory {
		if onProgress != nil {
			onProgress(i, total, baroItem.Item)
		}

		urlName, ok := index[normalizeItemName(baroItem.Item)]
		if !ok {
			continue
		}

		entries, err := fetchItemStatistics(urlName)
		time.Sleep(marketRequestPause)

		if err != nil {
			continue
		}

		row := Row{
			Name:    baroItem.Item,
			Ducats:  baroItem.Ducats,
			Credits: baroItem.Credits,
		}

		avgPrice, avgVolume, count := averageInWindow(entries, timeWindow)
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
