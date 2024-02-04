//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"os"
	"strconv"
	"time"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/ble/sphero"
	"gobot.io/x/gobot/v2/platforms/bleclient"
	"gobot.io/x/gobot/v2/platforms/mqtt"
)

const (
	FRENTE    = 0
	DERECHA   = 90
	ATRAS     = 180
	IZQUIERDA = 270
)

func main() {
	bleAdaptor := bleclient.NewAdaptor(os.Args[1])
	ollie := sphero.NewOllieDriver(bleAdaptor)

	mqttAdaptor := mqtt.NewAdaptor("tcp://iot.eclipse.org:1883", "ollie")

	work := func() {
		ollie.SetRGB(255, 0, 255)

		mqttAdaptor.On("sensores/dial", func(msg mqtt.Message) {
			val, _ := strconv.Atoi(string(msg.Payload()))

			if val > 2000 {
				ollie.SetRGB(0, 255, 0)
				return
			}
			if val > 1000 {
				ollie.SetRGB(255, 255, 0)
				return
			}
			ollie.SetRGB(255, 0, 0)
		})

		mqttAdaptor.On("rover/frente", func(msg mqtt.Message) {
			ollie.Roll(40, FRENTE)
			gobot.After(1*time.Second, func() {
				ollie.Stop()
			})
		})

		mqttAdaptor.On("rover/derecha", func(msg mqtt.Message) {
			ollie.Roll(40, DERECHA)
			gobot.After(1*time.Second, func() {
				ollie.Stop()
			})
		})

		mqttAdaptor.On("rover/atras", func(msg mqtt.Message) {
			ollie.Roll(40, ATRAS)
			gobot.After(1*time.Second, func() {
				ollie.Stop()
			})
		})

		mqttAdaptor.On("rover/izquierda", func(msg mqtt.Message) {
			ollie.Roll(40, IZQUIERDA)
			gobot.After(1*time.Second, func() {
				ollie.Stop()
			})
		})
	}

	robot := gobot.NewRobot("ollieBot",
		[]gobot.Connection{bleAdaptor, mqttAdaptor},
		[]gobot.Device{ollie},
		work,
	)

	robot.Start()
}
