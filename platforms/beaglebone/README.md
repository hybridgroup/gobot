# gobot-beaglebone

Gobot (http://gobot.io/) is a library for robotics and physical computing using Go

This library provides an adaptor and driver for the Beaglebone Black (http://beagleboard.org/Products/BeagleBone+Black/)

[![Build Status](https://travis-ci.org/hybridgroup/gobot-beaglebone.svg?branch=master)](https://travis-ci.org/hybridgroup/gobot-beaglebone) [![Coverage Status](https://coveralls.io/repos/hybridgroup/gobot-beaglebone/badge.png?branch=master)](https://coveralls.io/r/hybridgroup/gobot-beaglebone?branch=master)

## Getting Started

Install the library with: `go get -u github.com/hybridgroup/gobot-beaglebone`

## Cross compiling for the Beaglebone Black
You must first configure your Go environment for arm linux cross compiling

```bash
$ cd $GOROOT/src
$ GOOS=linux GOARCH=arm ./make.bash --no-clean
```

Then compile your Gobot program with
```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/blink.go
```

If you are running the default Angstrom linux through the usb->ethernet connection, you can simply upload your program and execute it with
``` bash
$ scp blink root@192.168.7.2:/home/root/
$ ssh -t root@192.168.7.2 "./blink"
```

## Example

```go
package main

import (
        "github.com/hybridgroup/gobot"
        "github.com/hybridgroup/gobot-beaglebone"
        "github.com/hybridgroup/gobot-gpio"
)

func main() {
        beaglebone := new(gobotBeaglebone.Beaglebone)
        beaglebone.Name = "beaglebone"

        led := gobotGPIO.NewLed(beaglebone)
        led.Name = "led"
        led.Pin = "P9_12"

        work := func() {
                gobot.Every("1s", func() { led.Toggle() })
        }

        robot := gobot.Robot{
                Connections: []interface{}{beaglebone},
                Devices:     []interface{}{led},
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
