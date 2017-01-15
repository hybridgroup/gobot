package ble

import (
	blelib "github.com/currantlabs/ble"
	"github.com/currantlabs/ble/darwin"
)

func defaultDevice(impl string) (d blelib.Device, err error) {
	return darwin.NewDevice()
}
