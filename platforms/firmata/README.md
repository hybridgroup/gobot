# Firmata

This package provides the adaptor for microcontrollers such as Arduino that support the [Firmata](http://firmata.org/wiki/Main_Page) protocol

## Getting Started

```
go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/gobot/platforms/firmata
```

## Example

```go
package main

import (
	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/platforms/firmata"
	"github.com/hybridgroup/gobot/platforms/gpio"
	"time"
)

func main() {
	gbot := gobot.NewGobot()
	adaptor := firmata.NewFirmataAdaptor("myFirmata", "/dev/ttyACM0")
	led := gpio.NewLedDriver(adaptor, "myLed", "13")
	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}
	gbot.Robots = append(gbot.Robots,
		gobot.NewRobot("blinkBot", []gobot.Connection{adaptor}, []gobot.Device{led}, work))
	gbot.Start()
}
```
## Hardware Support
The following firmata devices have been tested and are currently supported:

  - [Arduino uno r3](http://arduino.cc/en/Main/arduinoBoardUno)
  - [Teensy 3.0](http://www.pjrc.com/store/teensy3.html)

More devices are coming soon...