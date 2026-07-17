package main

import (
	"fmt"
	"time"
)

// ProgressFunc is called after each item is processed so the UI can show
// live progress (e.g. "12 / 46 items checked").
type ProgressFunc func(done, total int, currentItem string)

// StatusFunc reports coarse-grained status text to the UI.
type StatusFunc func(text string)

// buildRows fetches Baro's current inventory, matches each item against the
// warframe.market catalogue, and computes the 10-day price/volume averages
// for the window ending on Baro's arrival date.
//
// overrideActivation, if non-zero, replaces the arrival date reported by the
// worldstate API — useful for reproducing a specific historical stock such
// as the one on 2026-07-10.
func buildRows(overrideActivation time.Time, onStatus StatusFunc, onProgress ProgressFunc) ([]Row, AnalysisWindow, error) {
	notify := func(s string) {
		if onStatus != nil {
			onStatus(s)
		}
	}

	notify("Fetching Baro Ki'Teer inventory...")
	inv, err := fetchBaroInventory()
	if err != nil {
		return nil, AnalysisWindow{}, fmt.Errorf("baro inventory: %w", err)
	}

	var refDate time.Time
	if !overrideActivation.IsZero() {
		refDate = overrideActivation
	} else {
		refDate, err = inv.activationDate()
		if err != nil {
			return nil, AnalysisWindow{}, fmt.Errorf("determine stock date: %w", err)
		}
	}

	window := AnalysisWindow{
		Start: refDate.AddDate(0, 0, -9),
		End:   refDate,
	}

	notify("Fetching warframe.market item catalogue...")
	index, err := fetchMarketItemIndex()
	if err != nil {
		return nil, AnalysisWindow{}, fmt.Errorf("market item index: %w", err)
	}

	total := len(inv.Inventory)
	rows := make([]Row, 0, total)

	for i, baroItem := range inv.Inventory {
		if onProgress != nil {
			onProgress(i, total, baroItem.Item)
		}

		row := Row{
			Name:    baroItem.Item,
			Ducats:  baroItem.Ducats,
			Credits: baroItem.Credits,
		}

		urlName, ok := index[normalizeItemName(baroItem.Item)]
		if !ok {
			row.NoMarketData = true
			rows = append(rows, row)
			continue
		}

		entries, err := fetchItemStatistics(urlName)
		// Always pause to respect warframe.market's request rate guidance,
		// regardless of success/failure of the call above.
		time.Sleep(marketRequestPause)

		if err != nil {
			row.NoMarketData = true
			rows = append(rows, row)
			continue
		}

		avgPrice, avgVolume, count := averageInWindow(entries, window)
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

	return rows, window, nil
}
