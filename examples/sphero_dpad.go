package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/sphero"
)

func main() {
	gbot := gobot.NewMaster()
	a := api.NewAPI(gbot)
	a.Start()

	conn := sphero.NewAdaptor("/dev/rfcomm0")
	ball := sphero.NewSpheroDriver(conn)

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

	gbot.AddRobot(robot)

	gbot.Start()
}
