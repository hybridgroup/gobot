// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/ble"
)

func main() {
	bleAdaptor := ble.NewClientAdaptor(os.Args[1])
	info := ble.NewDeviceInformationDriver(bleAdaptor)

	work := func() {
		fmt.Println("Model number:", info.GetModelNumber())
		fmt.Println("Firmware rev:", info.GetFirmwareRevision())
		fmt.Println("Hardware rev:", info.GetHardwareRevision())
		fmt.Println("Manufacturer name:", info.GetManufacturerName())
		fmt.Println("PnPId:", info.GetPnPId())
	}

	robot := gobot.NewRobot("bleBot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{info},
		work,
	)

	robot.Start()
}
