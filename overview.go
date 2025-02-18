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
	/*var err error
	devices, err := make([]comfortcloud.Device, 0), nil

	devices = append(devices, comfortcloud.Device{"CS-Z35XKEW+4886", "3", "Stube", +3, 0, 0, true, true, comfortcloud.ModeAvl{1}, comfortcloud.Parameters{0, 3, 22.5, 0, 0, 2, 2, 0, 0, 0, 1, 0, 0, 0, 0, 0, 18.0, 2.0, 0},
		"CS-Z35XKEW", "62f3b46f6ef4d593a0694d2cfb1ddc32e4351c", 1, false})
	*/
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

type DeviceItem struct {
	NameLabel    *widget.Label
	IndoorLabel  *widget.Label
	OutdoorLabel *canvas.Text
	PowerToggle  *widget.Check
	TempLabel    *widget.Label
	IncreaseTemp *widget.Button
	DecreaseTemp *widget.Button
}

func createListItem(device comfortcloud.Device) *DeviceItem {
	nameLabel := widget.NewLabelWithStyle(device.DeviceName, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	indoorLabel := widget.NewLabel(fmt.Sprintf("🏠 %.1f°C", device.Parameters.InsideTemperature))
	outdoorLabel := canvas.NewText(fmt.Sprintf("🌤️ %.1f°C", device.Parameters.OutTemperature), color.RGBA{128, 128, 128, 255}) // Gray color
	powerToggle := widget.NewCheck("", nil)
	tempLabel := widget.NewLabel(fmt.Sprintf("%.1f°C", device.Parameters.TemperatureSet))
	increaseTemp := widget.NewButton("+", nil)
	decreaseTemp := widget.NewButton("-", nil)

	return &DeviceItem{
		NameLabel:    nameLabel,
		IndoorLabel:  indoorLabel,
		OutdoorLabel: outdoorLabel,
		PowerToggle:  powerToggle,
		TempLabel:    tempLabel,
		IncreaseTemp: increaseTemp,
		DecreaseTemp: decreaseTemp,
	}
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
			// Create a new container with placeholders
			nameLabel := canvas.NewText("", color.White)
			nameLabel.Alignment = fyne.TextAlignCenter
			nameLabel.TextSize = 20
			nameLabel.TextStyle.Bold = true
			indoorLabel := widget.NewLabel("")
			outdoorLabel := canvas.NewText("", color.Gray{Y: 100})
			powerToggle := widget.NewCheck("", nil)
			tempLabel := widget.NewLabel("")
			increaseTemp := widget.NewButton("+", nil)
			decreaseTemp := widget.NewButton("-", nil)

			// Create layout structure
			tempControls := container.NewHBox(decreaseTemp, tempLabel, increaseTemp)
			infoContainer := container.NewHBox(nameLabel, container.NewVBox(indoorLabel, outdoorLabel))
			controlContainer := container.NewVBox(powerToggle, tempControls)
			finalContainer := container.NewBorder(nil, nil, infoContainer, controlContainer)

			// Store the elements in a wrapper container
			wrapper := container.NewVBox(finalContainer)

			// Return wrapper (same reference will be used in the update function)
			return wrapper
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			device := deviceList[i]
			wrapper := o.(*fyne.Container) // The wrapper container

			// Extract the finalContainer inside the wrapper
			finalContainer := wrapper.Objects[0].(*fyne.Container)

			// Extract all objects inside finalContainer
			infoContainer := finalContainer.Objects[0].(*fyne.Container)
			controlContainer := finalContainer.Objects[1].(*fyne.Container)

			// Extract individual widgets
			nameLabel := infoContainer.Objects[0].(*canvas.Text)
			indoorLabel := infoContainer.Objects[1].(*fyne.Container).Objects[0].(*widget.Label)
			outdoorLabel := infoContainer.Objects[1].(*fyne.Container).Objects[1].(*canvas.Text)

			powerToggle := controlContainer.Objects[0].(*widget.Check)
			tempControls := controlContainer.Objects[1].(*fyne.Container)
			decreaseTemp := tempControls.Objects[0].(*widget.Button)
			tempLabel := tempControls.Objects[1].(*widget.Label)
			increaseTemp := tempControls.Objects[2].(*widget.Button)

			// Update UI elements with device data
			nameLabel.Text = device.DeviceName
			indoorLabel.SetText(fmt.Sprintf("🏠 %.1f°C", device.Parameters.InsideTemperature))
			outdoorLabel.Text = fmt.Sprintf("🌤️ %.1f°C", device.Parameters.OutTemperature)
			tempLabel.SetText(fmt.Sprintf("%.1f°C", device.Parameters.TemperatureSet))

			// Update Power Toggle
			powerToggle.SetChecked(device.Parameters.Operate == comfortcloud.PowerOn)
			powerToggle.OnChanged = func(checked bool) {
				var operate comfortcloud.Power
				if checked {
					operate = comfortcloud.PowerOn
				} else {
					operate = comfortcloud.PowerOff
				}
				err := client.SetDevice(deviceList[i].DeviceHashGuid, comfortcloud.WithPower(operate))
				if err != nil {
					errorLabel.Text = fmt.Sprintf("Error setting device: %v", err)
				} else {
					errorLabel.Text = ""
					deviceList[i].Parameters.Operate = operate
				}

			}

			// Update buttons
			increaseTemp.OnTapped = func() {
				temperature := deviceList[i].Parameters.TemperatureSet
				temperature += 0.5
				err := client.SetDevice(deviceList[i].DeviceHashGuid, comfortcloud.WithTemperature(temperature))
				if err != nil {
					errorLabel.Text = fmt.Sprintf("Error setting device: %v", err)
				} else {
					errorLabel.Text = ""
					deviceList[i].Parameters.TemperatureSet = temperature
					tempLabel.SetText(fmt.Sprintf("%.1f°C", deviceList[i].Parameters.TemperatureSet))
				}

			}
			decreaseTemp.OnTapped = func() {
				temperature := deviceList[i].Parameters.TemperatureSet
				temperature -= 0.5
				err := client.SetDevice(deviceList[i].DeviceHashGuid, comfortcloud.WithTemperature(temperature))
				if err != nil {
					errorLabel.Text = fmt.Sprintf("Error setting device: %v", err)
				} else {
					errorLabel.Text = ""
					deviceList[i].Parameters.TemperatureSet = temperature
					tempLabel.SetText(fmt.Sprintf("%.1f°C", deviceList[i].Parameters.TemperatureSet))
				}
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
