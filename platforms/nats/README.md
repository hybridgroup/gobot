# NATS

NATS is a lightweight messaging protocol perfect for your IoT/Robotics projects. It operates over TCP, offers a great
number of features but an incredibly simple Pub Sub style model of communicating broadcast messages. NATS is blazingly
fast as it is written in Go.

This repository contains the Gobot adaptor/drivers to connect to NATS servers. It uses the NATS Go Client available at
<https://github.com/nats-io/nats>. The NATS project is a project part of the [CNCF](https://www.cncf.io/). Find more
information on setting up a NATS server and its capability at <http://nats.io/>.

The NATS messaging protocol (<http://www.nats.io/documentation/internals/nats-protocol-demo/>) is really easy to work
with and can be practiced by setting up a NATS server using Go or Docker. For information on setting up a server using
the source code, visit <https://github.com/nats-io/gnatsd>. For information on the Docker image up on Docker Hub, see
<https://hub.docker.com/_/nats/>. Getting the server set up is very easy. The server itself is Golang, can be built for
different architectures and installs in a small footprint. This is an excellent way to get communications going between
your IoT and Robotics projects.

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

## How to Use

Before running the example, make sure you have an NATS server running somewhere you can connect to (for demo purposes
you could also use the public endpoint `demo.nats.io:4222`).

```go
package main

import (
  "fmt"
  "time"

  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/platforms/nats"
)

func main() {
  natsAdaptor := nats.NewAdaptor("nats://localhost:4222", 1234)

  work := func() {
    natsAdaptor.On("hello", func(msg nats.Message) {
      fmt.Printf("[Received on %q] %s\n", msg.Subject, string(msg.Data))
    })
    natsAdaptor.On("hola", func(msg nats.Message) {
      fmt.Printf("[Received on %q] %s\n", msg.Subject, string(msg.Data))
    })

    data := []byte("Hello Gobot!")
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

  if err := robot.Start(); err != nil {
		panic(err)
	}
}
```

To run with TLS enabled, set the URL scheme prefix to `tls://`. Make sure the NATS server has TLS enabled and use the
NATS option parameters to pass in the TLS settings to the adaptor. Refer to the [github.com/nats-io/go-nats](https://github.com/nats-io/go-nats)
README for more TLS option parameters.

```go
package main

import (
  "fmt"
  "time"

  natsio "github.com/nats-io/go-nats"
  "gobot.io/x/gobot/v2"
  "gobot.io/x/gobot/v2/platforms/nats"
)

func main() {
  natsAdaptor := nats.NewAdaptor("tls://localhost:4222", 1234, natsio.RootCAs("certs/ca.pem"))

  work := func() {
    natsAdaptor.On("hello", func(msg nats.Message) {
      fmt.Printf("[Received on %q] %s\n", msg.Subject, string(msg.Data))
    })
    natsAdaptor.On("hola", func(msg nats.Message) {
      fmt.Printf("[Received on %q] %s\n", msg.Subject, string(msg.Data))
    })

    data := []byte("Hello Gobot!")
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

  if err := robot.Start(); err != nil {
		panic(err)
	}
}
```

### Supported Features

* Publish messages
* Respond to incoming message events
* Support for Username/password authentication
* Support for NATS adaptor options to support TLS

### Upcoming Features

* Encoded messages (JSON)
* Exposing more NATS Features (tls)
* Simplified tests

## License

Copyright (c) 2013-2018 The Hybrid Group. Licensed under the Apache 2.0 license.
