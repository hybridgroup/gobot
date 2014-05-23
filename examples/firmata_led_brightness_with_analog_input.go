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
	sensor := gpio.NewAnalogSensor(firmataAdaptor, "sensor", "0")
	led := gpio.NewLed(firmataAdaptor, "led", "3")

	work := func() {
		gobot.Every(0.1*time.Second, func() {
			val := sensor.Read()
			brightness := uint8(gpio.ToPwm(val))
			fmt.Println("sensor", val)
			fmt.Println("brightness", brightness)
			led.Brightness(brightness)
		})
	}

	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("sensorBot", []gobot.Connection{firmataAdaptor}, []gobot.Device{sensor, led}, work))

	gbot.Start()
}
