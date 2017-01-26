package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mqtt"
)

func main() {
	mqttAdaptor := mqtt.NewAdaptor("tcp://test.mosquitto.org:1883", "pinger")
	mqttAdaptor.SetAutoReconnect(true)

	holaDriver := mqtt.NewDriver(mqttAdaptor, "hola")
	helloDriver := mqtt.NewDriver(mqttAdaptor, "hello")

	work := func() {
		helloDriver.On(mqtt.Data, func(data interface{}) {
			fmt.Println("hello")
		})

		holaDriver.On(mqtt.Data, func(data interface{}) {
			fmt.Println("hola")
		})

		data := []byte("o")
		gobot.Every(1*time.Second, func() {
			helloDriver.Publish(data)
		})

		gobot.Every(5*time.Second, func() {
			holaDriver.Publish(data)
		})
	}

	robot := gobot.NewRobot("mqttBot",
		[]gobot.Connection{mqttAdaptor},
		work,
	)

	robot.Start()
}
