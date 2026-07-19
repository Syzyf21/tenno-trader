package baroinvestor

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Syzyf21/tenno-trader/internal"
)

const ducatIconSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
  <defs>
    <linearGradient id="g" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#f6d87a"/>
      <stop offset="55%" stop-color="#d8a93a"/>
      <stop offset="100%" stop-color="#a97a1f"/>
    </linearGradient>
  </defs>
  <polygon points="32,4 58,20 58,44 32,60 6,44 6,20" fill="url(#g)" stroke="#7a5613" stroke-width="2"/>
  <polygon points="32,14 48,24 48,40 32,50 16,40 16,24" fill="none" stroke="#7a5613" stroke-width="2"/>
  <text x="32" y="38" font-family="Georgia, serif" font-size="18" font-weight="bold"
        text-anchor="middle" fill="#5c3d0e">D</text>
</svg>`

func DucatIconResource() fyne.Resource {
	return fyne.NewStaticResource("ducat_icon.svg", []byte(ducatIconSVG))
}

func BuildResultsTable(rows []internal.Row) *widget.Table {
	headers := []string{"Item", "Ducats", "Avg Platinum (10d)", "Avg Volume (10d)", "Plat / Ducat", "Data pts"}

	table := widget.NewTable(
		func() (int, int) {
			return len(rows), len(headers)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(id widget.TableCellID, obj fyne.CanvasObject) {
			label := obj.(*widget.Label)
			label.TextStyle = fyne.TextStyle{}

			r := rows[id.Row]
			switch id.Col {
			case 0:
				label.SetText(r.Name)
			case 1:
				label.SetText(fmt.Sprintf("%d", r.Ducats))
			case 2:
				if r.NoMarketData {
					label.SetText("no data")
				} else {
					label.SetText(fmt.Sprintf("%.1fp", r.AvgPlatinum))
				}
			case 3:
				if r.NoMarketData {
					label.SetText("-")
				} else {
					label.SetText(fmt.Sprintf("%.1f", r.AvgVolume))
				}
			case 4:
				if r.NoMarketData || r.Ducats == 0 {
					label.SetText("-")
				} else {
					label.SetText(fmt.Sprintf("%.3f", r.PlatPerDucat))
				}
			case 5:
				label.SetText(fmt.Sprintf("%d", r.DataPoints))
			}
		},
	)

	table.ShowHeaderRow = true
	table.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("")
	}
	table.UpdateHeader = func(id widget.TableCellID, obj fyne.CanvasObject) {
		label := obj.(*widget.Label)
		if id.Row == -1 && id.Col >= 0 && id.Col < len(headers) {
			label.TextStyle = fyne.TextStyle{Bold: true}
			label.SetText(headers[id.Col])
		} else {
			label.SetText("")
		}
	}

	table.SetColumnWidth(0, 260)
	table.SetColumnWidth(1, 80)
	table.SetColumnWidth(2, 160)
	table.SetColumnWidth(3, 160)
	table.SetColumnWidth(4, 120)
	table.SetColumnWidth(5, 90)

	return table
}
