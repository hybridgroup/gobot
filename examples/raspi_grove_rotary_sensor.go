// +build example
//
// Do not build by default.

package main

import (
	"fmt"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/aio"
	"gobot.io/x/gobot/drivers/i2c"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	board := raspi.NewAdaptor()
	ads1015 := i2c.NewADS1015Driver(board)
	sensor := aio.NewGroveRotaryDriver(ads1015, "0")

	work := func() {
		sensor.On(aio.Data, func(data interface{}) {
			fmt.Println("sensor", data)
		})
	}

	robot := gobot.NewRobot("sensorBot",
		[]gobot.Connection{board},
		[]gobot.Device{ads1015, sensor},
		work,
	)

	robot.Start()
}
