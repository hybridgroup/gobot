//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth address or name as the first param:

	go run examples/bleclient_r2d2.go D2-1234

 NOTE: sudo is required to use BLE in Linux
 NOTE: GODEBUG=cgocheck=0 is required in the prefix of the command to use BLE in OsX
*/

//nolint:gosec // ok here
package main

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/sphero"
	"gobot.io/x/gobot/v2/platforms/bleclient"
	"os"
	"time"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	r2 := sphero.NewR2D2Driver(bleAdaptor)

	work := func() {
		r2.SetDomePosition(0)
		// gobot.Rand(90) will get a random color value up to 90
		r := uint8(180)
		g := uint8(180)
		b := uint8(45)
		r2.SetRGB(r, g, b)
		time.Sleep(1 * time.Second)
		r2.SetDomePosition(160)
		time.Sleep(1 * time.Second)
		r2.PlaySound(1737 /* ALARM1 */, sphero.PlaybackImmediate)
		time.Sleep(500 * time.Millisecond)
		r2.SetDomePosition(-90)
	}

	robot := gobot.NewRobot("r2Bot",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{r2},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
