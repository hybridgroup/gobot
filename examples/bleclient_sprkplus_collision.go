//go:build example
// +build example

//
// Do not build by default.

/*
 How to run
 Pass the Bluetooth address or name as the first param:

	go run examples/bb8-collision.go BB-1234

 NOTE: sudo is required to use BLE in Linux
*/

package main

import (
	"fmt"
	"os"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/sphero"
	"gobot.io/x/gobot/v2/platforms/bleclient"
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	ball := sphero.NewSPRKPlusDriver(bleAdaptor)

	work := func() {
		_ = ball.On("collision", func(data interface{}) {
			fmt.Printf("collision detected = %+v \n", data)
			ball.SetRGB(255, 0, 0)
		})

		ball.SetRGB(0, 255, 0)
		ball.Roll(80, 0)
	}

	robot := gobot.NewRobot("sprkplus",
		[]gobot.Connection{bleAdaptor},
		[]gobot.Device{ball},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
