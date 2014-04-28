package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/beaglebone"
	"github.com/hybridgroup/gobot/gpio"
)

func main() {
	beagleboneAdaptor := beaglebone.NewBeagleboneAdaptor()
	beagleboneAdaptor.Name = "beaglebone"

	sensor := gpio.NewAnalogSensorDriver(beagleboneAdaptor)
	sensor.Name = "sensor"
	sensor.Pin = "P9_33"

	led := gpio.NewLedDriver(beagleboneAdaptor)
	led.Name = "led"
	led.Pin = "P9_14"

	work := func() {
		gobot.Every("0.1s", func() {
			val := sensor.Read()
			brightness := uint8(gpio.ToPwm(val))
			fmt.Println("sensor", val)
			fmt.Println("brightness", brightness)
			led.Brightness(brightness)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{beagleboneAdaptor},
		Devices:     []gobot.Device{sensor, led},
		Work:        work,
	}

	robot.Start()
}
