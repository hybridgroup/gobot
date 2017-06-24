// +build example
//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "13")
	neo := firmata.NewNeopixelDriver(firmataAdaptor, "6", 5)

	work := func() {
		cols := []uint32{0xff0000, 0x00ff00, 0x0000ff, 0xffff00, 0xffffff}

		for i := 0; i < 50; i++ {
			// rotate the slice of colors
			x, c := cols[0], cols[1:]
			cols = append(c, x)

			for j := 0; j < 5; j++ {
				neo.SetPixel(uint16(j), cols[j])
			}

			neo.Show()
			time.Sleep(100 * time.Millisecond)
		}

		neo.Off()
	}

	robot := gobot.NewRobot("neoBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{neo, led},
		work,
	)

	robot.Start()
}
