package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/seb-ehm/panasonic-comfort-cloud/comfortcloud"
)

func RunApp() {
	a := app.New()
	w := a.NewWindow("AC Controller")
	var client *comfortcloud.Client
	// Set the initial content to the login screen

	w.SetContent(createLoginScreen(client, w))

	// Start control loops in the background
	go StartControlLoops()

	w.ShowAndRun()
}
