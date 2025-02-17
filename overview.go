package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/seb-ehm/panasonic-comfort-cloud/comfortcloud"
	"image/color"
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
			// Template for each item
			nameLabel := widget.NewLabelWithStyle("Device Name", fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
			indoorLabel := widget.NewLabel("🏠 22.0°C")
			outdoorLabel := canvas.NewText("🌤️ 15.0°C", color.RGBA{128, 128, 128, 255}) // Gray color
			powerToggle := widget.NewCheck("", nil)
			tempLabel := widget.NewLabel("23.0°C")
			increaseTemp := widget.NewButton("+", nil)
			decreaseTemp := widget.NewButton("-", nil)

			// Layout with horizontal box
			tempControls := container.NewHBox(decreaseTemp, tempLabel, increaseTemp)
			infoContainer := container.NewVBox(nameLabel, indoorLabel, outdoorLabel)
			controlContainer := container.NewVBox(powerToggle, tempControls)

			return container.NewBorder(nil, nil, infoContainer, controlContainer)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			// Populate data
			device := deviceList[i]
			containers := o.(*fyne.Container).Objects
			infoContainer := containers[0].(*fyne.Container)
			controlContainer := containers[1].(*fyne.Container)

			nameLabel := infoContainer.Objects[0].(*widget.Label)
			indoorLabel := infoContainer.Objects[1].(*widget.Label)
			outdoorLabel := infoContainer.Objects[2].(*canvas.Text)

			powerToggle := controlContainer.Objects[0].(*widget.Check)
			tempControls := controlContainer.Objects[1].(*fyne.Container)
			tempLabel := tempControls.Objects[1].(*widget.Label)
			increaseTemp := tempControls.Objects[2].(*widget.Button)
			decreaseTemp := tempControls.Objects[0].(*widget.Button)

			// Update content
			nameLabel.SetText(device.DeviceName)
			indoorLabel.SetText(fmt.Sprintf("🏠 %.1f°C", device.Parameters.InsideTemperature))
			outdoorLabel.Text = fmt.Sprintf("🌤️ %.1f°C", device.Parameters.OutTemperature)
			outdoorLabel.Refresh()
			powerToggle.SetChecked(device.Parameters.Operate == comfortcloud.PowerOn)
			tempLabel.SetText(fmt.Sprintf("%.1f°C", device.Parameters.TemperatureSet))

			// Handlers
			powerToggle.OnChanged = func(checked bool) {
				var operate comfortcloud.Power
				if checked {
					operate = comfortcloud.PowerOn
				} else {
					operate = comfortcloud.PowerOff
				}
				deviceList[i].Parameters.Operate = operate
				fmt.Println("Power toggled:", checked)
			}
			increaseTemp.OnTapped = func() {
				deviceList[i].Parameters.TemperatureSet += 0.5
				tempLabel.SetText(fmt.Sprintf("%.1f°C", deviceList[i].Parameters.TemperatureSet))
			}
			decreaseTemp.OnTapped = func() {
				deviceList[i].Parameters.TemperatureSet -= 0.5
				tempLabel.SetText(fmt.Sprintf("%.1f°C", deviceList[i].Parameters.TemperatureSet))
			}
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
