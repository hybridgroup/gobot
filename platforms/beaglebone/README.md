# Beaglebone

This package provides the Gobot adaptor for the [Beaglebone Black](http://beagleboard.org/Products/BeagleBone+Black/)

## Getting Started

## Installing
```
go get github.com/hybridgroup/gobot && go install github.com/hybridgroup/platforms/beaglebone
```

## Cross compiling for the Beaglebone Black
You must first configure your Go environment for arm linux cross compiling

```bash
$ cd $GOROOT/src
$ GOOS=linux GOARCH=arm ./make.bash --no-clean
```

Then compile your Gobot program with
```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/beaglebone_blink.go
```

If you are running the official Angstrom or Debian linux through the usb->ethernet connection, you can simply upload your program and execute it with
``` bash
$ scp beaglebone_blink root@192.168.7.2:/home/root/
$ ssh -t root@192.168.7.2 "./beaglebone_blink"
```

## Example

```go
package main

import (
        "github.com/hybridgroup/gobot"
        "github.com/hybridgroup/gobot/platforms/beaglebone"
        "github.com/hybridgroup/gobot/platforms/gpio"
        "time"
)

func main() {
        gbot := gobot.NewGobot()

        adaptor := beaglebone.NewBeagleboneAdaptor("beaglebone")
        led := gpio.NewLedDriver(adaptor, "led", "P9_12")

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