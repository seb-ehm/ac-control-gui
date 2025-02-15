package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/seb-ehm/panasonic-comfort-cloud/comfortcloud"
)

func createDetailScreen(client *comfortcloud.Client, window fyne.Window, device comfortcloud.Device) fyne.CanvasObject {
	// Display AC details
	nameLabel := widget.NewLabel("Name: " + device.DeviceName)
	statusLabel := widget.NewLabel(fmt.Sprintf("Status: %v", device.Parameters.Operate))
	roomTempLabel := widget.NewLabel(fmt.Sprintf("Room Temp: %.1f°C", float64(device.Parameters.InsideTemperature)))
	customerTempLabel := widget.NewLabel(fmt.Sprintf("Customer Set Temp: %.1f°C", float64(device.Parameters.TemperatureSet)))

	// Manual control buttons
	onButton := widget.NewButton("Turn On", func() {
		err := client.SetDevice(device.DeviceGuid, comfortcloud.WithPower(comfortcloud.PowerOn))
		if err != nil {
			dialog.ShowError(err, window)
		}
	})
	offButton := widget.NewButton("Turn Off", func() {
		err := client.SetDevice(device.DeviceGuid, comfortcloud.WithPower(comfortcloud.PowerOff))
		if err != nil {
			dialog.ShowError(err, window)
		}
	})

	// Back button to return to the overview
	backButton := widget.NewButton("Back", func() {
		window.SetContent(createOverviewScreen(client, window))
	})

	return container.NewVBox(
		nameLabel,
		statusLabel,
		roomTempLabel,
		customerTempLabel,
		container.NewHBox(onButton, offButton),
		backButton,
	)
}
