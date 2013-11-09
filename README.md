# Gobot

A Go library for robotics and physical computing 

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

It provides a simple, yet powerful way to create solutions that incorporate multiple, different hardware devices at the same time.

Want to use Ruby or Javascript on robots? Check out our sister projects Artoo (http://artoo.io) and Cylon.js (http://cylonjs.com/)


## Getting Started

Install the library with: `go get -u github.com/hybridgroup/gobot`

Then install additional libraries for whatever hardware support you want to use from your robot. For example, `go get -u github.com/hybridgroup/gobot-sphero` to use Gobot with a Sphero.

## Examples

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

  connections := []interface{} {
    spheroAdaptor,
  }
  devices := []interface{} {
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

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. The following robotics and physical computing platforms are currently supported:

  - [Beaglebone Black](http://beagleboard.org/Products/BeagleBone+Black/) <=> [Library](https://github.com/hybridgroup/gobot-beaglebone)
  - [Sphero](http://www.gosphero.com/) <=> [Library](https://github.com/hybridgroup/gobot-sphero)

More platforms and drivers are coming soon...

## Documentation
We're busy adding documentation to our web site at http://gobot.io/ please check there as we continue to work on Gobot

Thank you!

## Contributing
In lieu of a formal styleguide, take care to maintain the existing coding style.
Add unit tests for any new or changed functionality

## License
Copyright (c) 2013 The Hybrid Group. Licensed under the Apache 2.0 license.
