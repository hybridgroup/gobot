// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/platforms/intel-iot/edison"
)

func main() {
	board := edison.NewAdaptor()
	sensor := aio.NewGroveRotaryDriver(board, "0")

	work := func() {
		sensor.On(aio.Data, func(data interface{}) {
			fmt.Println("sensor", data)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{sensor},
		work,
	)

	robot.Start()
}
