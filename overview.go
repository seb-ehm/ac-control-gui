package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seb-ehm/panasonic-comfort-cloud/comfortcloud"
	"sync"
)

var (
	acState   sync.Map // Shared state for all ACs
	stateLock sync.RWMutex
)

// Dummy function to get devices (replace with actual API call)
func getDevices() []comfortcloud.Device {
	return []comfortcloud.Device{
		{DeviceGuid: "1", DeviceName: "Living Room"},
		{DeviceGuid: "2", DeviceName: "Bedroom"},
		{DeviceGuid: "3", DeviceName: "Another Room"},
	}
}

func createOverviewScreen(client *comfortcloud.Client, window fyne.Window) fyne.CanvasObject {
	// Create a list widget to display AC systems
	acList := widget.NewList(
		func() int {
			count := 0
			acState.Range(func(key, value interface{}) bool {
				count++
				return true
			})
			return count
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("AC System")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			var device comfortcloud.Device
			acState.Range(func(key, value interface{}) bool {
				if i == 0 {
					device = value.(comfortcloud.Device)
					return false
				}
				i--
				return true
			})
			o.(*widget.Label).SetText(device.DeviceName)
		},
	)

	// Handle AC selection
	acList.OnSelected = func(id widget.ListItemID) {
		var device comfortcloud.Device
		acState.Range(func(key, value interface{}) bool {
			if id == 0 {
				device = value.(comfortcloud.Device)
				return false
			}
			id--
			return true
		})
		window.SetContent(createDetailScreen(client, window, device))
	}

	return container.NewBorder(nil, nil, nil, nil, acList)
}
