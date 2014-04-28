package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/firmata"
	"github.com/hybridgroup/gobot/gpio"
)

func main() {
	firmataAdaptor := firmata.NewFirmataAdaptor()
	firmataAdaptor.Name = "firmata"
	firmataAdaptor.Port = "/dev/ttyACM0"

	sensor := gpio.NewAnalogSensor(firmataAdaptor)
	sensor.Name = "sensor"
	sensor.Pin = "0"

	led := gpio.NewLed(firmataAdaptor)
	led.Name = "led"
	led.Pin = "3"

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
		Connections: []gobot.Connection{firmataAdaptor},
		Devices:     []gobot.Device{sensor, led},
		Work:        work,
	}

	robot.Start()
}
