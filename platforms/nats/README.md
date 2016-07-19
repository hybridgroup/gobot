# NATS

NATS is a lightweight messaging protocol perfect for your IoT/Robotics projects. It operates over TCP, offers a great number of features but an incredibly simple Pub Sub style model of communicating broadcast messages. NATS is blazingly fast as it is written in Go. 

This repository contains the Gobot adaptor/drivers to connect to NATS servers. It uses the NATS Go Client available at https://github.com/nats-io/nats. The NATS project is maintained by Nats.io and sponsored by Apcera. Find more information on setting up a NATS server and its capability at http://nats.io/.

The NATS messaging protocol (http://www.nats.io/documentation/internals/nats-protocol-demo/) is really easy to work with and can be practiced by setting up a NATS server using Go or Docker. For information on setting up a server using the source code, visit https://github.com/nats-io/gnatsd. For information on the Docker image up on Docker Hub, see https://hub.docker.com/_/nats/. Getting the server set up is very easy. The server itself is Golang, can be built for different architectures and installs in a small footprint. This is an excellent way to get communications going between your IoT and Robotics projects.

## How to Install

Install running:

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/nats
```

## How to Use

Before running the example, make sure you have an NATS server running somewhere you can connect to

```go
package main

import (
  "fmt"
  "time"

  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot/platforms/nats"
)

func main() {
  gbot := gobot.NewGobot()

  natsAdaptor := nats.NewNatsAdaptor("nats", "localhost:4222", 1234)

  work := func() {
    natsAdaptor.On("hello", func(data []byte) {
      fmt.Println("hello")
    })
    natsAdaptor.On("hola", func(data []byte) {
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

  gbot.AddRobot(robot)

  gbot.Start()
}
```

## Supported Features

* Publish messages
* Respond to incoming message events

## Upcoming Features

* Support for Username/password
* Encoded messages (JSON)
* Exposing more NATS Features (tls)
* Simplified tests

## Contributing

For our contribution guidelines, please go to https://github.com/hybridgroup/gobot/blob/master/CONTRIBUTING.md

## License

Copyright (c) 2013-2016 The Hybrid Group. Licensed under the Apache 2.0 license.
