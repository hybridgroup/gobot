# gobot-firmata

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This library provides an adaptor for microcontrollers such as Arduino that support the Firmata protocol (http://firmata.org/wiki/Main_Page)

[![Build Status](https://travis-ci.org/hybridgroup/gobot-firmata.svg?branch=master)](https://travis-ci.org/hybridgroup/gobot-firmata) [![Coverage Status](https://coveralls.io/repos/hybridgroup/gobot-firmata/badge.png)](https://coveralls.io/r/hybridgroup/gobot-firmata)

## Getting Started

Install the library with: `go get -u github.com/hybridgroup/gobot-firmata`

## Example

```go
package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot-firmata"
	"github.com/hybridgroup/gobot-gpio"
)

func main() {

	firmata := new(gobotFirmata.FirmataAdaptor)
	firmata.Name = "firmata"
	firmata.Port = "/dev/ttyACM0"

	led := gobotGPIO.NewLed(firmata)
	led.Name = "led"
	led.Pin = "13"

	work := func() {
		gobot.Every("1s", func() {
			led.Toggle()
		})
	}

	robot := gobot.Robot{
		Connections: []gobot.Connection{firmata},
		Devices:     []gobot.Device{led},
		Work:        work,
	}

	robot.Start()
}
```
## Hardware Support
The following firmata devices have been tested and are currently supported:

  - [Arduino uno r3](http://arduino.cc/en/Main/arduinoBoardUno)
  - [Teensy 3.0](http://www.pjrc.com/store/teensy3.html)

More devices are coming soon...

## Documentation
We're busy adding documentation to our web site at http://gobot.io/ please check there as we continue to work on Gobot

Thank you!

## Contributing
In lieu of a formal styleguide, take care to maintain the existing coding style. Add unit tests for any new or changed functionality.

## License
Copyright (c) 2013 The Hybrid Group. Licensed under the Apache 2.0 license.
