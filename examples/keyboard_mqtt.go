//go:build example
// +build example

//
// Do not build by default.

package main

import (
	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/platforms/keyboard"
	"gobot.io/x/gobot/v2/platforms/mqtt"
)

func main() {
	keys := keyboard.NewDriver()
	mqttAdaptor := mqtt.NewAdaptor("tcp://iot.eclipse.org:1883", "conductor")

	work := func() {
		_ = keys.On(keyboard.Key, func(data interface{}) {
			key := data.(keyboard.KeyEvent)

			switch key.Key {
			case keyboard.ArrowUp:
				mqttAdaptor.Publish("rover/frente", []byte{})
			case keyboard.ArrowRight:
				mqttAdaptor.Publish("rover/derecha", []byte{})
			case keyboard.ArrowDown:
				mqttAdaptor.Publish("rover/atras", []byte{})
			case keyboard.ArrowLeft:
				mqttAdaptor.Publish("rover/izquierda", []byte{})
			}
		})
	}

	robot := gobot.NewRobot("keyboardbot",
		[]gobot.Connection{mqttAdaptor},
		[]gobot.Device{keys},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
