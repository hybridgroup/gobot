//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/firmata"
	"gobot.io/x/gobot/v2/platforms/mqtt"
)

func main() {
	mqttAdaptor := mqtt.NewAdaptor("tcp://test.mosquitto.org:1883", "blinker")
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "13")

	work := func() {
		mqttAdaptor.On("lights/on", func(msg mqtt.Message) {
			led.On()
		})
		mqttAdaptor.On("lights/off", func(msg mqtt.Message) {
			led.Off()
		})
		data := []byte("")
		gobot.Every(1*time.Second, func() {
			mqttAdaptor.Publish("lights/on", data)
		})
		gobot.Every(2*time.Second, func() {
			mqttAdaptor.Publish("lights/off", data)
		})
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{mqttAdaptor, firmataAdaptor},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
