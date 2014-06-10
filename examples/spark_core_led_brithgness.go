package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/spark"
	"time"
)

func main() {
	master := gobot.NewGobot()

	sparkCore := spark.NewSparkCoreAdaptor("spark", "device_id", "access_token")
	led := gpio.NewLedDriver(sparkCore, "led", "A1")

	work := func() {
		brightness := uint8(0)
		fadeAmount := uint8(25)

		gobot.Every(500*time.Millisecond, func() {
			led.Brightness(brightness)
			brightness = brightness + fadeAmount
			if brightness == 0 || brightness == 255 {
				fadeAmount = -fadeAmount
			}
		})
	}
	master.Robots = append(master.Robots,
		gobot.NewRobot("spark", []gobot.Connection{sparkCore}, []gobot.Device{led}, work))

	master.Start()
}
