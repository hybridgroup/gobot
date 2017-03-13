// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
)

func main() {
	board := edison.NewAdaptor()
	sensor := aio.NewGroveTemperatureSensorDriver(board, "0")

	work := func() {
		gobot.Every(500*time.Millisecond, func() {
			fmt.Println("current temp (c): ", sensor.Temperature())
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{sensor},
		work,
	)

	robot.Start()
}
