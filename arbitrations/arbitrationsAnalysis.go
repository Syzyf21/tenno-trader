package arbitrations

import (
	"fmt"
	"log"
	"time"

	"github.com/Syzyf21/tenno-trader/internal"
)

type ProgressFunc func(done, total int, currentItem string)

type StatusFunc func(text string)

func BuildRows(onStatus StatusFunc, onProgress ProgressFunc, isMax bool) ([]internal.ArbitrationRow, internal.AnalysisWindow, error) {
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

		avgPrice, avgVolume, count := internal.AverageInWindow(entries, timeWindow, isMax)
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
	notify("Fetching complete. Inserting data to database")

	sql := `DELETE FROM arbitration_data`
	_, err = internal.DBConn.Query(sql)
	if err != nil {
		return nil, timeWindow, fmt.Errorf("Error while clearing arbitration data database: %v", err)
	}

	for _, row := range rows {
		sql := `INSERT INTO arbitration_data (item_name, vitus_cost, avg_platinum, avg_volume, plat_vitus, has_market_data) 
		VALUES (?, ?, ?, ?, ?, ?)`

		stmt, err := internal.DBConn.Prepare(sql)
		if err != nil {
			log.Fatalf("Błąd przygotowania zapytania: %q", err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(row.Name, row.Vitus, row.AvgPlatinum, row.AvgVolume, row.PlatPerVitus, row.NoMarketData)
		if err != nil {
			log.Fatalf("Błąd podczas wstawiania danych: %q", err)
		}
	}

	return rows, timeWindow, nil
}
