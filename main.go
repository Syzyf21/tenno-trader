package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Syzyf21/tenno-trader/arbitrations"
	"github.com/Syzyf21/tenno-trader/baroinvestor"
	"github.com/Syzyf21/tenno-trader/internal"
)

func main() {
	a := app.NewWithID("github.com/Syzyf21/tenno-trader")
	w := a.NewWindow("Tenno Trader")
	w.Resize(fyne.NewSize(1300, 780))
	w.CenterOnScreen()

	_, err := internal.InitDB("./data_clean.db")
	if err != nil {
		fmt.Printf("Error opening database: %v\n", err)
		return
	}
	defer internal.DBConn.Close()
	fmt.Print("Connected to database successfully.\n")

	status := widget.NewLabel("Select \"Baro Investor\" on the left to check investing opportunities, if new Baro stock is available, it will be automatically loaded.")
	progress := widget.NewProgressBar()
	progress.Hide()

	resultsHolder := container.NewStack(widget.NewLabel(""))

	header := widget.NewLabel("")
	header.TextStyle = fyne.TextStyle{Bold: true}

	content := container.NewBorder(
		container.NewVBox(header, status, progress),
		nil, nil, nil,
		resultsHolder,
	)

	loadVoidTraderData := func() {
		status.SetText("Working...")
		progress.Show()
		progress.SetValue(0)
		header.SetText("")

		go func() {
			rows, window, err := baroinvestor.BuildRows(
				func(text string) {
					fyne.Do(func() {
						status.SetText(text)
					})
				},
				func(done, total int, currentItem string) {
					fyne.Do(func() {
						if total > 0 {
							progress.SetValue(float64(done) / float64(total))
						}
						if currentItem != "" {
							status.SetText(fmt.Sprintf("Checking %d/%d: %s", done+1, total, currentItem))
						}
					})
				},
			)

			fyne.Do(func() {
				progress.Hide()
				if err != nil {
					status.SetText("Error: " + err.Error())
					return
				}

				status.SetText(fmt.Sprintf("Loaded %d items.", len(rows)))
				header.SetText(fmt.Sprintf(
					"Baro Ki'Teer stock — averages computed for %s to %s",
					window.Start.Format("02.01.2006"),
					window.End.Format("02.01.2006"),
				))

				table := baroinvestor.BuildResultsTable(rows)
				resultsHolder.Objects = []fyne.CanvasObject{table}
				resultsHolder.Refresh()
			})
		}()
	}

	loadArbitrationsData := func() {
		status.SetText("Working...")
		progress.Show()
		progress.SetValue(0)
		header.SetText("")

		go func() {
			rows, window, err := arbitrations.BuildRows(
				func(text string) {
					fyne.Do(func() {
						status.SetText(text)
					})
				},
				func(done, total int, currentItem string) {
					fyne.Do(func() {
						if total > 0 {
							progress.SetValue(float64(done) / float64(total))
						}
						if currentItem != "" {
							status.SetText(fmt.Sprintf("Checking %d/%d: %s", done+1, total, currentItem))
						}
					})
				},
			)

			fyne.Do(func() {
				progress.Hide()
				if err != nil {
					status.SetText("Error: " + err.Error())
					return
				}

				status.SetText(fmt.Sprintf("Loaded %d items.", len(rows)))
				header.SetText(fmt.Sprintf(
					"Baro Ki'Teer stock — averages computed for %s to %s",
					window.Start.Format("02.01.2006"),
					window.End.Format("02.01.2006"),
				))

				table := arbitrations.BuildResultsTable(rows)
				resultsHolder.Objects = []fyne.CanvasObject{table}
				resultsHolder.Refresh()
			})
		}()
	}

	sidebar := buildSidebar(loadVoidTraderData, loadArbitrationsData)

	split := container.NewHSplit(sidebar, content)
	split.Offset = 0.2

	w.SetContent(split)
	w.ShowAndRun()
}
