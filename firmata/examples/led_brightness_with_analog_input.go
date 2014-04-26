package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-firmata"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {
	firmata := new(gobotFirmata.FirmataAdaptor)
	firmata.Name = "firmata"
	firmata.Port = "/dev/ttyACM0"

	sensor := gobotGPIO.NewAnalogSensor(firmata)
	sensor.Name = "sensor"
	sensor.Pin = "0"

	led := gobotGPIO.NewLed(firmata)
	led.Name = "led"
	led.Pin = "3"

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
		Connections: []gobot.Connection{firmata},
		Devices:     []gobot.Device{sensor, led},
		Work:        work,
	}

	robot.Start()
}
