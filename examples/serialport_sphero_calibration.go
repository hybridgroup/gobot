//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/api"
	"gobot.io/x/gobot/v2/drivers/serial/sphero"
	"gobot.io/x/gobot/v2/platforms/keyboard"
	"gobot.io/x/gobot/v2/platforms/serialport"
)

func main() {
	manager := gobot.NewManager()
	a := api.NewAPI(manager)
	a.Start()

	ballConn := serialport.NewAdaptor("/dev/rfcomm0")
	ball := sphero.NewSpheroDriver(ballConn)

	keys := keyboard.NewDriver()

	calibrating := false

	work := func() {
		_ = keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			switch key.Key {
			case keyboard.ArrowUp:
				if calibrating {
					break
				}
				ball.Roll(100, 0)
			case keyboard.ArrowDown:
				if calibrating {
					break
				}
				ball.Roll(100, 100)
			case keyboard.ArrowLeft:
				ball.Roll(100, 270)
			case keyboard.ArrowRight:
				ball.Roll(100, 90)
			case keyboard.Spacebar:
				if calibrating {
					ball.FinishCalibration()
				} else {
					ball.StartCalibration()
				}
				calibrating = !calibrating
			}
		})
	}

	robot := gobot.NewRobot("sphero-calibration",
		[]gobot.Connection{ballConn},
		[]gobot.Device{ball, keys},
		work,
	)

	manager.AddRobot(robot)

	if err := manager.Start(); err != nil {
		panic(err)
	}
}
