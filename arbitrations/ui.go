package arbitrations

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Syzyf21/tenno-trader/internal"
)

const vitusIconSVG = `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64">
  <defs>
    <!-- Metallic dark gunmetal shell gradient -->
    <linearGradient id="shell" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#464b54"/>
      <stop offset="45%" stop-color="#2c2f35"/>
      <stop offset="100%" stop-color="#121316"/>
    </linearGradient>
    
    <!-- Glowing Cyan/Teal Energy -->
    <linearGradient id="cyanGlow" x1="0%" y1="0%" x2="0%" y2="100%">
      <stop offset="0%" stop-color="#55ffff"/>
      <stop offset="100%" stop-color="#00a8ff"/>
    </linearGradient>
  </defs>

  <!-- Outer Asymmetrical Pod Shell -->
  <path d="M32,4 C48,4 56,22 54,42 C52,54 44,60 32,60 C20,60 10,52 10,38 C10,20 20,4 32,4 Z" 
        fill="url(#shell)" stroke="#1a1c1e" stroke-width="2"/>

  <!-- Inner Faceplate / Shield Layer -->
  <path d="M32,14 C42,14 46,26 44,42 C43,50 38,54 32,54 C26,54 21,50 20,42 C18,26 22,14 32,14 Z" 
        fill="#1c1e22" stroke="#32373f" stroke-width="1.5"/>

  <!-- Glowing Cyan Energy Channels -->
  <path d="M32,18 L32,50 M26,30 L38,30 M24,40 L40,40" 
        fill="none" stroke="url(#cyanGlow)" stroke-width="2" stroke-linecap="round" opacity="0.85"/>

  <!-- Hexis Style Central Metallic Overlays -->
  <rect x="29" y="24" width="6" height="14" rx="2" fill="#0d0e10" stroke="#32373f" stroke-width="1"/>
  <path d="M22,34 L28,34 M36,34 L42,34" stroke="#0d0e10" stroke-width="3" stroke-linecap="round"/>

  <!-- Core Glowing Nodes (Arbiters Signature Dots) -->
  <circle cx="32" cy="22" r="2" fill="#a6ffff" filter="drop-shadow(0px 0px 2px #00ffff)"/>
  <circle cx="32" cy="44" r="2" fill="#a6ffff" filter="drop-shadow(0px 0px 2px #00ffff)"/>
</svg>`

func VitusIconResource() fyne.Resource {
	return fyne.NewStaticResource("vitus_icon.svg", []byte(vitusIconSVG))
}

func BuildResultsTable(rows []internal.ArbitrationRow) *widget.Table {
	headers := []string{"Item", "Vitus Essence", "Avg Platinum (90d)", "Avg Volume (90d)", "Plat / Vitus"}

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
				label.SetText(fmt.Sprintf("%d", r.Vitus))
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
				if r.NoMarketData || r.Vitus == 0 {
					label.SetText("-")
				} else {
					label.SetText(fmt.Sprintf("%.3f", r.PlatPerVitus))
				}
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
