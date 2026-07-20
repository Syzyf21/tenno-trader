package baroinvestor

import (
	"database/sql"
	"fmt"

	"github.com/Syzyf21/tenno-trader/internal"
)

func GetVoidTraderItems(db *sql.DB) ([]internal.BaroInventoryItemDB, error) {
	sql := `SELECT * FROM void_trader_stock`
	rows, err := db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching void trader items: %w", err)
	}
	defer rows.Close()

	var items []internal.BaroInventoryItemDB
	for rows.Next() {
		var item internal.BaroInventoryItemDB
		if err := rows.Scan(&item.ID, &item.Name, &item.Ducats, &item.Last_Seen); err != nil {
			return nil, fmt.Errorf("Error while scanning void trader item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}
