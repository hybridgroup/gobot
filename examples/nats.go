// +build example
//
// Do not build by default.

package main

import (
	"fmt"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/nats"
)

func main() {
	natsAdaptor := nats.NewAdaptorWithAuth("localhost:4222", 1234, "user", "pass")

	work := func() {
		natsAdaptor.On("hello", func(msg nats.Message) {
			fmt.Println("hello")
		})
		natsAdaptor.On("hola", func(msg nats.Message) {
			fmt.Println("hola")
		})
		data := []byte("o")
		gobot.Every(1*time.Second, func() {
			natsAdaptor.Publish("hello", data)
		})
		gobot.Every(5*time.Second, func() {
			natsAdaptor.Publish("hola", data)
		})
	}

	robot := gobot.NewRobot("natsBot",
		[]gobot.Connection{natsAdaptor},
		work,
	)

	robot.Start()
}
