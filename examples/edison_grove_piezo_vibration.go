//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/aio"
	"gobot.io/x/gobot/v2/platforms/intel-iot/edison"
)

func main() {
	board := edison.NewAdaptor()
	sensor := aio.NewGrovePiezoVibrationSensorDriver(board, "0")

	work := func() {
		sensor.On(aio.Vibration, func(data interface{}) {
			fmt.Println("got one!")
		})
	}

	robot := gobot.NewRobot("bot",
		[]gobot.Connection{board},
		[]gobot.Device{sensor},
		work,
	)

	robot.Start()
}
