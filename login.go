package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/seb-ehm/panasonic-comfort-cloud/comfortcloud"
)

func NeedsLogin(client *comfortcloud.Client) (*comfortcloud.Client, bool) {
	//Check if token file can be used instead of username / password
	err := client.Login()
	if err == nil {
		return client, false
	}
	return nil, true
	//return client, false
}

func createLoginScreen(client *comfortcloud.Client, window fyne.Window) fyne.CanvasObject {

	username := widget.NewEntry()
	password := widget.NewPasswordEntry()

	// Create login form
	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Username", Widget: username},
			{Text: "Password", Widget: password},
		},
		OnSubmit: func() {
			client = comfortcloud.NewClient(username.Text, password.Text, tokenFile)
			err := client.Login()
			if err != nil {
				dialog.ShowError(err, window)
			} else {
				window.SetContent(createOverviewScreen(client, window))
			}
		},
	}

	// Wrap the form in a fixed-size container
	formContainer := container.NewGridWrap(fyne.NewSize(300, form.MinSize().Height), form)

	// Center the form
	centeredForm := container.NewCenter(formContainer)

	return centeredForm
}
