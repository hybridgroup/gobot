package ble

import (
	blelib "github.com/hybridgroup/ble"
	"github.com/hybridgroup/ble/darwin"
)

func defaultDevice(impl string) (d blelib.Device, err error) {
	return darwin.NewDevice()
}
