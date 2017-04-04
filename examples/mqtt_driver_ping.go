// +build example
//
// Do not build by default.

// TO RUN:
//  go run ./examples/mqtt_driver_ping.go <SERVER>
//
// EXAMPLE:
//	go run ./examples/mqtt_driver_ping.go ssl://iot.eclipse.org:8883
//
package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/mqtt"
)

func main() {
	mqttAdaptor := mqtt.NewAdaptor(os.Args[1], "pinger")
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
		[]gobot.Device{helloDriver, holaDriver},
		work,
	)

	robot.Start()
}
