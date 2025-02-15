package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/seb-ehm/panasonic-comfort-cloud/comfortcloud"
	"os"
)

func checkConfigFile() bool {
	// Check if the config file exists and meets preconditions
	_, err := os.Stat("config.json")
	return os.IsNotExist(err)
}

func NeedsLogin() bool {
	return true
}

func createLoginScreen(client *comfortcloud.Client, window fyne.Window) fyne.CanvasObject {
	username := widget.NewEntry()
	password := widget.NewPasswordEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: username},
			{Text: "Password", Widget: password},
		},
		OnSubmit: func() {
			// Simulate login logic
			if username.Text == "admin" && password.Text == "password" {
				// Switch to the overview screen
				window.SetContent(createOverviewScreen(client, window))
			} else {
				dialog.ShowError(fmt.Errorf("invalid credentials"), window)
			}
		},
	}

	return container.NewVBox(form)
}
