package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/seb-ehm/panasonic-comfort-cloud/comfortcloud"
)

func RunApp() {
	a := app.New()
	w := a.NewWindow("AC Controller")
	// Set a minimum window size
	w.Resize(fyne.NewSize(800, 600))

	client := comfortcloud.NewClient("", "", tokenFile)

	client, needsLogin := NeedsLogin(client)
	if needsLogin {
		w.SetContent(createLoginScreen(client, w))
	} else {
		w.SetContent(createOverviewScreen(client, w))
	}

	// Start control loops in the background
	go StartControlLoops()

	w.ShowAndRun()
}
