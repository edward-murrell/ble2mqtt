package adapt

import "tinygo.org/x/bluetooth"

type Device interface {
	Name() string
	GetState() any                                                          // This needs to return a struct that will be turned into JSON.
	UpdateDevice(update *bluetooth.ScanResult) (change bool, failure error) // Update device with whatever data was received by from the BLE sensor.
}
