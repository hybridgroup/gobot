package main

import (
	"github.com/hybridgroup/gobot"
	"time"
)

func main() {

	beaglebone := new(Beaglebone)
	beaglebone.Name = "Beaglebone"

	led := NewLed(beaglebone)
	led.Driver = Driver{
		Name: "led",
		Pin:  "P9_12",
	}

	connections := []interface{}{
		beaglebone,
	}
	devices := []interface{}{
		led,
	}

	work := func() {
		Every(1000*time.Millisecond, func() { led.Toggle() })
	}

	robot := Robot{
		Connections: connections,
		Devices:     devices,
		Work:        work,
	}

	robot.Start()
}
