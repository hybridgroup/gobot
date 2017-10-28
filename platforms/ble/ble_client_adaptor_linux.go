package ble

import (
	blelib "github.com/hybridgroup/ble"
	"github.com/hybridgroup/ble/linux"
)

func defaultDevice(impl string) (d blelib.Device, err error) {
	return linux.NewDevice()
}
