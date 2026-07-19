package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/Syzyf21/tenno-trader/arbitrations"
	"github.com/Syzyf21/tenno-trader/baroinvestor"
)

func buildSidebar(onSelect func()) fyne.CanvasObject {
	baroInvestor := widget.NewButtonWithIcon("Baro Investor", baroinvestor.DucatIconResource(), onSelect)
	baroInvestor.Alignment = widget.ButtonAlignLeading
	baroInvestor.Importance = widget.LowImportance

	arbitrations := widget.NewButtonWithIcon("Arbitrations", arbitrations.VitusIconResource(), onSelect)

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
