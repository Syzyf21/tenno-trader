package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
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

	u := newUI(w)

	sidebar := u.buildSidebar()
	content := u.buildContent()

	split := container.NewHSplit(sidebar, content)
	split.Offset = 0.2

	w.SetContent(split)
	w.ShowAndRun()
}
