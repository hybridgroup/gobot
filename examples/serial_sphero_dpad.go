//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/drivers/serial"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
	master := gobot.NewMaster()
	a := api.NewAPI(master)
	a.Start()

	conn := serialport.NewAdaptor("/dev/rfcomm0")
	ball := serial.NewSpheroDriver(conn)

	robot := gobot.NewRobot("sphero-dpad",
		[]gobot.Connection{conn},
		[]gobot.Device{ball},
	)

	robot.AddCommand("move", func(params map[string]interface{}) interface{} {
		direction := params["direction"].(string)

		switch direction {
		case "up":
			ball.Roll(100, 0)
		case "down":
			ball.Roll(100, 180)
		case "left":
			ball.Roll(100, 270)
		case "right":
			ball.Roll(100, 90)
		}

		time.Sleep(2 * time.Second)
		ball.Stop()
		return "ok"
	})

	master.AddRobot(robot)

	if err := master.Start(); err != nil {
		panic(err)
	}
}
