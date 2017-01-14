package ble

import (
	blelib "github.com/currantlabs/ble"
	"github.com/currantlabs/ble/linux"
)

func defaultDevice(impl string) (d blelib.Device, err error) {
	return linux.NewDevice()
}
