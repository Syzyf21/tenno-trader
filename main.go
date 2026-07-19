package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.NewWithID("github.com/Syzyf21/tenno-trader")
	w := a.NewWindow("Tenno Trader")
	w.Resize(fyne.NewSize(1300, 780))
	w.CenterOnScreen()

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

	loadData := func() {
		status.SetText("Working...")
		progress.Show()
		progress.SetValue(0)
		header.SetText("")

		go func() {
			rows, window, err := buildRows(
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

				table := buildResultsTable(rows)
				resultsHolder.Objects = []fyne.CanvasObject{table}
				resultsHolder.Refresh()
			})
		}()
	}

	sidebar := buildSidebar(loadData)

	split := container.NewHSplit(sidebar, content)
	split.Offset = 0.2

	w.SetContent(split)
	w.ShowAndRun()
}
