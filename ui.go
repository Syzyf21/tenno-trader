package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Syzyf21/tenno-trader/arbitrations"
	"github.com/Syzyf21/tenno-trader/baroinvestor"
)

func buildSidebar(onVoidTraderSelect func(), onArbitrationsSelect func()) fyne.CanvasObject {
	baroInvestor := widget.NewButtonWithIcon("Baro Investor", baroinvestor.DucatIconResource(), onVoidTraderSelect)
	baroInvestor.Alignment = widget.ButtonAlignLeading
	baroInvestor.Importance = widget.LowImportance

	arbitrations := widget.NewButtonWithIcon("Arbitrations", arbitrations.VitusIconResource(), onArbitrationsSelect)
	arbitrations.Alignment = widget.ButtonAlignLeading
	arbitrations.Importance = widget.LowImportance

	title := widget.NewLabelWithStyle("TENNO TRADER", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})

	box := container.NewVBox(
		title,
		widget.NewSeparator(),
		baroInvestor,
		arbitrations,
		layout.NewSpacer(),
	)
	padded := container.NewPadded(box)
	return padded
}
