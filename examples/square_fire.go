package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/api"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/intel-iot/edison"
)

func main() {
	gbot := gobot.NewGobot()
	a := api.NewAPI(gbot)
	a.Start()

	board := edison.NewEdisonAdaptor("edison")
	red := gpio.NewLedDriver(board, "red", "3")
	green := gpio.NewLedDriver(board, "green", "5")
	blue := gpio.NewLedDriver(board, "blue", "6")

	button := gpio.NewButtonDriver(board, "button", "7")

	enabled := true
	work := func() {
		red.Brightness(0xff)
		green.Brightness(0x00)
		blue.Brightness(0x00)

		flash := false
		on := true

		gobot.Every(50*time.Millisecond, func() {
			if enabled {
				if flash {
					if on {
						red.Brightness(0x00)
						green.Brightness(0xff)
						blue.Brightness(0x00)
						on = false
					} else {
						red.Brightness(0x00)
						green.Brightness(0x00)
						blue.Brightness(0xff)
						on = true
					}
				}
			}
		})

		gobot.On(button.Event("push"), func(data interface{}) {
			flash = true
		})

		gobot.On(button.Event("release"), func(data interface{}) {
			flash = false
			red.Brightness(0x00)
			green.Brightness(0x00)
			blue.Brightness(0xff)
		})
	}

	robot := gobot.NewRobot(
		"square of fire",
		[]gobot.Connection{board},
		[]gobot.Device{red, green, blue, button},
		work,
	)

	robot.AddCommand("enable", func(params map[string]interface{}) interface{} {
		enabled = !enabled
		return enabled
	})

	gbot.AddRobot(robot)

	gbot.Start()
}
