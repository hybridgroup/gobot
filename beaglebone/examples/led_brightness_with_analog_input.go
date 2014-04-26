package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-beaglebone"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {
	beaglebone := new(gobotBeaglebone.Beaglebone)
	beaglebone.Name = "beaglebone"

	sensor := gobotGPIO.NewAnalogSensor(beaglebone)
	sensor.Name = "sensor"
	sensor.Pin = "P9_33"

	led := gobotGPIO.NewLed(beaglebone)
	led.Name = "led"
	led.Pin = "P9_14"

	work := func() {
		gobot.Every("0.1s", func() {
			val := sensor.Read()
			brightness := uint8(gobotGPIO.ToPwm(val))
			fmt.Println("sensor", val)
			fmt.Println("brightness", brightness)
			led.Brightness(brightness)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{beaglebone},
		Devices:     []gobot.Device{sensor, led},
		Work:        work,
	}

	robot.Start()
}
