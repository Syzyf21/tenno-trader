package arbitrations

import (
	"database/sql"
	"fmt"

	"github.com/Syzyf21/tenno-trader/internal"
)

func GetArbitrationItems(db *sql.DB) ([]internal.ArbitrationItem, error) {
	sql := `SELECT * FROM arbitration_shop`
	rows, err := db.Query(sql)
	if err != nil {
		return nil, fmt.Errorf("Error while fetching arbitration items: %w", err)
	}
	defer rows.Close()

	var items []internal.ArbitrationItem
	for rows.Next() {
		var item internal.ArbitrationItem
		if err := rows.Scan(&item.ID, &item.Name, &item.Vitus); err != nil {
			return nil, fmt.Errorf("Error while scanning arbitration item: %w", err)
		}
		items = append(items, item)
	}

	return items, nil
}
