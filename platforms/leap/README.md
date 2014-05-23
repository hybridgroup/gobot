# gobot-leapmotion

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This library provides an adaptor and driver for the Leap Motion (https://www.leapmotion.com/)

## Getting Started

Install the library with: `go get -u github.com/hybridgroup/gobot-leapmotion`

## Example

```go
package main

import (
	"fmt"
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-leapmotion"
)

func main() {
	leapAdaptor := new(gobotLeap.LeapAdaptor)
	leapAdaptor.Name = "leap"
	leapAdaptor.Port = "127.0.0.1:6437"

	leap := gobotLeap.NewLeap(leapAdaptor)
	leap.Name = "leap"

	work := func() {
		gobot.On(leap.Events["Message"], func(data interface{}) {
			fmt.Println(data)
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{leapAdaptor},
		Devices:     []gobot.Device{leap},
		Work:        work,
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

