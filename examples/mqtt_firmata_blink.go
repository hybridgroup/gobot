package main

import (
  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot/platforms/mqtt"
  "github.com/hybridgroup/gobot/platforms/firmata"
  "github.com/hybridgroup/gobot/platforms/gpio"
)

func main() {
  gbot := gobot.NewGobot()

  mqttAdaptor := mqtt.NewMqttAdaptor("server", "localhost:1883")
  firmataAdaptor := firmata.NewFirmataAdaptor("arduino", "/dev/ttyACM0")
  led := gpio.NewLedDriver(firmataAdaptor, "led", "13")

  work := func() {
    mqttAdaptor.On('lights/on', func(data interface{}) {
      led.On()
    })
    mqttAdaptor.On('lights/off', func(data interface{}) {
      led.Off()
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
