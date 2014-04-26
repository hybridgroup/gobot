# gobot-sphero

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This library provides an adaptor and driver for the Sphero robot from Orbotix (http://www.gosphero.com/)

[![Build Status](https://travis-ci.org/hybridgroup/gobot-sphero.svg?branch=master)](https://travis-ci.org/hybridgroup/gobot-sphero) [![Coverage Status](https://coveralls.io/repos/hybridgroup/gobot-sphero/badge.png)](https://coveralls.io/r/hybridgroup/gobot-sphero)

## Getting Started

Install the library with: `go get -u github.com/hybridgroup/gobot-sphero`

## Example

```go
package main
import (
  "github.com/hybridgroup/gobot"
  "github.com/hybridgroup/gobot-sphero"
  "fmt"
)

func main() {

  spheroAdaptor := new(gobotSphero.SpheroAdaptor)
  spheroAdaptor.Name = "Sphero"
  spheroAdaptor.Port = "127.0.0.1:4560"

  sphero := gobotSphero.NewSphero(spheroAdaptor)
  sphero.Name = "Sphero"

  connections := []gobot.Connection {
    spheroAdaptor,
  }
  devices := []gobot.Device {
    sphero,
  }

  work := func(){
    gobot.Every("2s", func(){ 
      sphero.Roll(100, uint16(gobot.Random(0, 360))) 
    })
  }
  
  robot := gobot.Robot{
      Connections: connections, 
      Devices: devices,
      Work: work,
  }

  robot.Start()
}
```

## Documentation
We're busy adding documentation to our web site at http://gobot.io/ please check there as we continue to work on Gobot

Thank you!

## Contributing
In lieu of a formal styleguide, take care to maintain the existing coding style. Add unit tests for any new or changed functionality.

## License
Copyright (c) 2013 The Hybrid Group. Licensed under the Apache 2.0 license.
