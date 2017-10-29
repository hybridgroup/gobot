package ble

import (
	blelib "github.com/go-ble/ble"
	"github.com/go-ble/ble/darwin"
)

func defaultDevice(impl string) (d blelib.Device, err error) {
	return darwin.NewDevice()
}
