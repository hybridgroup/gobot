package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"github.com/hybridgroup/gobot/platforms/mqtt"
)

func main() {
	gbot := gobot.NewGobot()

	mqttAdaptor := mqtt.NewMqttAdaptor("server", "tcp://test.mosquitto.org:1883", "blinker")
	firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "led", "13")

	work := func() {
		mqttAdaptor.On("lights/on", func(data []byte) {
			led.On()
		})
		mqttAdaptor.On("lights/off", func(data []byte) {
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

	gbot.AddRobot(robot)

	gbot.Start()
}
