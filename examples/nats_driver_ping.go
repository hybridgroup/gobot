// +build example
//
// Do not build by default.

// TO RUN:
//  go run ./examples/nats_driver_ping.go <SERVER>
//
// EXAMPLE:
//	go run ./examples/nats_driver_ping.go tls://nats.demo.io:4443
//
package main

import (
	"fmt"
	"os"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/nats"
)

func main() {
	natsAdaptor := nats.NewAdaptor(os.Args[1], 1234)

	holaDriver := nats.NewDriver(natsAdaptor, "hola")
	helloDriver := nats.NewDriver(natsAdaptor, "hello")

	work := func() {
		helloDriver.On(nats.Data, func(msg nats.Message) {
			fmt.Println("hello")
		})

		holaDriver.On(nats.Data, func(msg nats.Message) {
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

	robot := gobot.NewRobot("natsBot",
		[]gobot.Connection{natsAdaptor},
		[]gobot.Device{helloDriver, holaDriver},
		work,
	)

	robot.Start()
}
