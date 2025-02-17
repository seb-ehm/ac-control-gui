package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seb-ehm/panasonic-comfort-cloud/comfortcloud"
	"sync"
)

var (
	deviceList []comfortcloud.Device
	deviceLock sync.RWMutex
)

// Dummy function to get devices (replace with actual API call)
func getDevices(client *comfortcloud.Client) error {
	deviceLock.RLock()
	defer deviceLock.RUnlock()
	devices, err := client.GetDevices()
	if err != nil {
		return err
	}
	// Return a copy to prevent external modifications
	devicesCopy := make([]comfortcloud.Device, len(devices))
	copy(devicesCopy, devices)
	deviceList = devicesCopy
	fmt.Println(deviceList)
	return nil

}

func createOverviewScreen(client *comfortcloud.Client, window fyne.Window) fyne.CanvasObject {
	// Acquire read lock before accessing deviceList
	err := getDevices(client)
	errorLabel := widget.NewLabel("")
	if err != nil {
		errorLabel.Text = fmt.Sprintf("Error getting devices: %v", err)
	}
	deviceLock.RLock()
	defer deviceLock.RUnlock()

	// Create a list widget to display AC systems
	acList := widget.NewList(
		func() int {
			return len(deviceList)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("AC System that has a lot of data in it")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			currentDevice := deviceList[i]
			label := fmt.Sprintf("%s, %s, %f", currentDevice.DeviceName, currentDevice.Parameters.Operate, currentDevice.Parameters.TemperatureSet)
			o.(*widget.Label).SetText(label)
		},
	)

	// Handle AC selection
	acList.OnSelected = func(id widget.ListItemID) {
		deviceLock.RLock()
		selectedDevice := deviceList[id]
		deviceLock.RUnlock()

		window.SetContent(createDetailScreen(client, window, selectedDevice))
	}

	return container.NewBorder(nil, nil, nil, nil, acList)
}
