package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	firmataAdaptor := firmata.NewFirmataAdaptor("firmata", "/dev/ttyACM0")
	sensor := gpio.NewAnalogSensorDriver(firmataAdaptor, "sensor", "0")
	led := gpio.NewLedDriver(firmataAdaptor, "led", "3")

	work := func() {
		gobot.Every(100*time.Millisecond, func() {
			val := sensor.Read()
			brightness := uint8(gobot.ToScale(gobot.FromScale(float64(val), 0, 1024), 0, 255))
			fmt.Println("sensor", val)
			fmt.Println("brightness", brightness)
			led.Brightness(brightness)
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("sensorBot", []gobot.Connection{firmataAdaptor}, []gobot.Device{sensor, led}, work))

	gbot.Start()
}
