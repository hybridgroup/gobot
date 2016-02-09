# MQTT

MQTT is an Internet of Things connectivity protocol featuring a lightweight publish/subscribe messaging transport. It is useful for its small code footprint and minimal network bandwidth usage.

This repository contains the Gobot adaptor/drivers to connect to MQTT servers. It uses the Paho MQTT Golang client package (https://eclipse.org/paho/) created and maintained by the Eclipse Foundation (https://github.com/eclipse) thank you!

For more info about the MQTT machine to machine messaging standard, go to http://mqtt.org/

## How to Install

Install running:

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/mqtt
```

## How to Use

Before running the example, make sure you have an MQTT message broker running somewhere you can connect to

```go
package main

import (
  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot/platforms/mqtt"
  "fmt"
  "time"
)

func main() {
  gbot := gobot.NewGobot()

  mqttAdaptor := mqtt.NewMqttAdaptor("server", "tcp://0.0.0.0:1883", "pinger")

  work := func() {
    mqttAdaptor.On("hello", func(data []byte) {
      fmt.Println("hello")
    })
    mqttAdaptor.On("hola", func(data []byte) {
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
```

## Supported Features

* Publish messages
* Respond to incoming message events

## Contributing

For our contribution guidelines, please go to https://github.com/hybridgroup/gobot/blob/master/CONTRIBUTING.md

## License

Copyright (c) 2013-2016 The Hybrid Group. Licensed under the Apache 2.0 license.
