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

	led := gobotGPIO.NewLed(firmata)
	led.Name = "led"
	led.Pin = "13"

	connections := []interface{}{
		firmata,
	}
	devices := []interface{}{
		led,
	}

	work := func() {
		gobot.Every("1s", func() {
			led.Toggle()
			if led.IsOn() {
				fmt.Println("On")
			} else {
				fmt.Println("Off")
			}
		})
	}

	robot := gobot.Robot{
		Connections: connections,
		Devices:     devices,
		Work:        work,
	}

	robot.Start()
}
