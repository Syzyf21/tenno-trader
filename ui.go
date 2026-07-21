package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Syzyf21/tenno-trader/arbitrations"
	"github.com/Syzyf21/tenno-trader/baroinvestor"
)

type ui struct {
	window        fyne.Window
	status        *widget.Label
	progress      *widget.ProgressBar
	header        *widget.Label
	headerButtons *fyne.Container
	resultsHolder *fyne.Container
}

func newUI(w fyne.Window) *ui {
	status := widget.NewLabel("Select \"Baro Investor\" on the left to check investing opportunities, if new Baro stock is available, it will be automatically loaded.")

	progress := widget.NewProgressBar()
	progress.Hide()

	header := widget.NewLabel("")
	header.TextStyle = fyne.TextStyle{Bold: true}

	return &ui{
		window:        w,
		status:        status,
		progress:      progress,
		header:        header,
		headerButtons: container.NewHBox(),
		resultsHolder: container.NewStack(widget.NewLabel("")),
	}
}

func (ui *ui) buildContent() fyne.CanvasObject {
	headerRow := container.NewBorder(nil, nil, ui.header, ui.headerButtons)

	return container.NewBorder(
		container.NewVBox(headerRow, ui.status, ui.progress),
		nil, nil, nil,
		ui.resultsHolder,
	)
}

func (ui *ui) buildSidebar() fyne.CanvasObject {
	baroInvestor := widget.NewButtonWithIcon("Baro Investor", baroinvestor.DucatIconResource(), ui.loadVoidTraderData)
	baroInvestor.Alignment = widget.ButtonAlignLeading
	baroInvestor.Importance = widget.LowImportance

	arbitrations := widget.NewButtonWithIcon("Arbitrations", arbitrations.VitusIconResource(), ui.loadArbitrationsData)
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

type viewButton int

const (
	viewNone viewButton = iota
	viewBaro
	viewLiveBaro
	viewArbitrations
)

func (ui *ui) setHeaderButtons(view viewButton) {
	switch view {
	case viewBaro:
		ui.headerButtons.Objects = []fyne.CanvasObject{
			widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
				ui.loadVoidTraderData()
			}),
			widget.NewButtonWithIcon("Live Baro Data", baroinvestor.DucatIconResource(), func() {
				ui.loadLiveVoidTraderData()
			}),
		}
	case viewLiveBaro:
		ui.headerButtons.Objects = []fyne.CanvasObject{
			widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
				ui.loadLiveVoidTraderData()
			}),
			widget.NewButtonWithIcon("Rotation Baro Data", baroinvestor.DucatIconResource(), func() {
				ui.loadVoidTraderData()
			}),
		}
	case viewArbitrations:
		ui.headerButtons.Objects = []fyne.CanvasObject{
			widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
				ui.loadArbitrationsData()
			}),
		}
	default:
		ui.headerButtons.Objects = nil
	}
	ui.headerButtons.Refresh()
}

func (ui *ui) loadVoidTraderData() {
	ui.status.SetText("Working...")
	ui.progress.Show()
	ui.progress.SetValue(0)
	ui.header.SetText("")

	go func() {
		rows, window, err := baroinvestor.BuildRows(
			func(text string) {
				fyne.Do(func() {
					ui.status.SetText(text)
				})
			},
			func(done, total int, currentItem string) {
				fyne.Do(func() {
					if total > 0 {
						ui.progress.SetValue(float64(done) / float64(total))
					}
					if currentItem != "" {
						ui.status.SetText(fmt.Sprintf("Checking %d/%d: %s", done+1, total, currentItem))
					}
				})
			}, true,
		)

		fyne.Do(func() {
			ui.progress.Hide()
			if err != nil {
				ui.status.SetText("Error: " + err.Error())
				return
			}

			ui.status.SetText(fmt.Sprintf("Loaded %d items.", len(rows)))
			ui.header.SetText(fmt.Sprintf(
				"Current Baro Ki'Teer stock — averages computed for %s to %s",
				window.Start.Format("02.01.2006"),
				window.End.Format("02.01.2006"),
			))
			ui.setHeaderButtons(viewBaro)

			table := baroinvestor.BuildResultsTable(rows)
			ui.resultsHolder.Objects = []fyne.CanvasObject{table}
			ui.resultsHolder.Refresh()
		})
	}()
}

func (ui *ui) loadLiveVoidTraderData() {
	ui.status.SetText("Working...")
	ui.progress.Show()
	ui.progress.SetValue(0)
	ui.header.SetText("")

	go func() {
		rows, window, err := baroinvestor.BuildRows(
			func(text string) {
				fyne.Do(func() {
					ui.status.SetText(text)
				})
			},
			func(done, total int, currentItem string) {
				fyne.Do(func() {
					if total > 0 {
						ui.progress.SetValue(float64(done) / float64(total))
					}
					if currentItem != "" {
						ui.status.SetText(fmt.Sprintf("Checking %d/%d: %s", done+1, total, currentItem))
					}
				})
			}, false,
		)

		fyne.Do(func() {
			ui.progress.Hide()
			if err != nil {
				ui.status.SetText("Error: " + err.Error())
				return
			}

			ui.status.SetText(fmt.Sprintf("Loaded %d items.", len(rows)))
			ui.header.SetText(fmt.Sprintf(
				"Current Baro Ki'Teer stock — averages computed for %s to %s",
				window.Start.Format("02.01.2006"),
				window.End.Format("02.01.2006"),
			))
			ui.setHeaderButtons(viewLiveBaro)

			table := baroinvestor.BuildResultsTable(rows)
			ui.resultsHolder.Objects = []fyne.CanvasObject{table}
			ui.resultsHolder.Refresh()
		})
	}()
}

func (ui *ui) loadArbitrationsData() {
	ui.status.SetText("Working...")
	ui.progress.Show()
	ui.progress.SetValue(0)
	ui.header.SetText("")

	go func() {
		rows, window, err := arbitrations.BuildRows(
			func(text string) {
				fyne.Do(func() {
					ui.status.SetText(text)
				})
			},
			func(done, total int, currentItem string) {
				fyne.Do(func() {
					if total > 0 {
						ui.progress.SetValue(float64(done) / float64(total))
					}
					if currentItem != "" {
						ui.status.SetText(fmt.Sprintf("Checking %d/%d: %s", done+1, total, currentItem))
					}
				})
			},
		)

		fyne.Do(func() {
			ui.progress.Hide()
			if err != nil {
				ui.status.SetText("Error: " + err.Error())
				return
			}

			ui.status.SetText(fmt.Sprintf("Loaded %d items.", len(rows)))
			ui.header.SetText(fmt.Sprintf(
				"Arbitrations Shop Market data — averages computed for %s to %s",
				window.Start.Format("02.01.2006"),
				window.End.Format("02.01.2006"),
			))
			ui.setHeaderButtons(viewArbitrations)

			table := arbitrations.BuildResultsTable(rows)
			ui.resultsHolder.Objects = []fyne.CanvasObject{table}
			ui.resultsHolder.Refresh()
		})
	}()
}
