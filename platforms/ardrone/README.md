# Gobot for ardrone

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This repository contains the Gobot adaptor for ardrone.

For more information about Gobot, check out the github repo at
https://github.com/hybridgroup/gobot

[![Build Status](https://travis-ci.org/hybridgroup/gobot-ardrone.svg?branch=master)](https://travis-ci.org/hybridgroup/gobot-ardrone) [![Coverage Status](https://coveralls.io/repos/hybridgroup/gobot-ardrone/badge.png)](https://coveralls.io/r/hybridgroup/gobot-ardrone)

## Installing
```
go get github.com/hybridgroup/gobot-ardrone
```
## Using
```go
package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-ardrone"
)

func main() {

	ardroneAdaptor := new(gobotArdrone.ArdroneAdaptor)
	ardroneAdaptor.Name = "Drone"

	drone := gobotArdrone.NewArdrone(ardroneAdaptor)
	drone.Name = "Drone"

	work := func() {
		drone.TakeOff()
		gobot.On(drone.Events["Flying"], func(data interface{}) {
			gobot.After("1s", func() {
				drone.Right(0.1)
			})
			gobot.After("2s", func() {
				drone.Left(0.1)
			})
			gobot.After("3s", func() {
				drone.Land()
			})
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{ardroneAdaptor},
		Devices:     []gobot.Device{drone},
		Work:        work,
	}

	robot.Start()
}
```

## License

Copyright (c) 2013-2014 The Hybrid Group. Licensed under the Apache 2.0 license.
