# MQTT

MQTT is an Internet of Things connectivity protocol featuring a lightweight publish/subscribe messaging transport. It is
useful for its small code footprint and minimal network bandwidth usage.

This repository contains the Gobot adaptor/driver to connect to MQTT servers. It uses the Paho MQTT Golang client package
(<https://eclipse.org/paho/>) created and maintained by the Eclipse Foundation (<https://github.com/eclipse>) thank you!

For more info about the MQTT machine to machine messaging standard, go to <http://mqtt.org/>.

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

## How to Use

Before running the example, make sure you have an MQTT message broker running somewhere you can connect to

```go
package main

import (
  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/platforms/mqtt"
  "fmt"
  "time"
)

func main() {
  mqttAdaptor := mqtt.NewAdaptor("tcp://0.0.0.0:1883", "pinger")

  work := func() {
    mqttAdaptor.On("hello", func(msg mqtt.Message) {
      fmt.Println(msg)
    })
    mqttAdaptor.On("hola", func(msg mqtt.Message) {
      fmt.Println(msg)
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

  robot.Start()
}
```

## Supported Features

* Publish messages
* Respond to incoming message events

## License

Copyright (c) 2013-2018 The Hybrid Group. Licensed under the Apache 2.0 license.
