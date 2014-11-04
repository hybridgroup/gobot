package main

import (
  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot/platforms/mqtt"
  "fmt"
  "time"
)

func main() {
  gbot := gobot.NewGobot()

  mqttAdaptor := mqtt.NewMqttAdaptor("server", "tcp://0.0.0.0:1883")

  work := func() {
    mqttAdaptor.On("hello", func(data interface{}) {
      fmt.Println("hello")
    })
    mqttAdaptor.On("hola", func(data interface{}) {
      fmt.Println("hola")
    })
    data := []byte("o")
    gobot.Every(1*time.Second, func() {
      mqttAdaptor.Publish("hello", data)
    })
    gobot.Every(5*time.Second, func() {
      mqttAdaptor.Publish("hola", data)
    })
  }

  robot := gobot.NewRobot("mqttBot",
    []gobot.Connection{mqttAdaptor},
    work,
  )

  gbot.AddRobot(robot)

  gbot.Start()
}
