# C.H.I.P.

The [C.H.I.P.](http://www.getchip.com/) is a small, inexpensive ARM based single board computer, with many different IO interfaces available on the [pin headers](http://docs.getchip.com/#pin-headers).

We recommend updating to the latest Debian OS when using the C.H.I.P., however Gobot should also support older versions of the OS, should your application require this.

For documentation about the C.H.I.P. platform click [here](http://docs.getchip.com/).

## How to Install
```
go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/chip
```

To be able to use the built in device tree overlay manager, you need to install the required
modified device tree compiler "dtc" from https://github.com/atenart/dtc.

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
package main

import (
    "fmt"

    "gobot.io/x/gobot"
    "gobot.io/x/gobot/drivers/gpio"
    "gobot.io/x/gobot/platforms/chip"
)

func main() {
    chipAdaptor := chip.NewAdaptor()
    button := gpio.NewButtonDriver(chipAdaptor, "XIO-P0")

    work := func() {
        gobot.On(button.Event("push"), func(data interface{}) {
            fmt.Println("button pressed")
        })

        gobot.On(button.Event("release"), func(data interface{}) {
            fmt.Println("button released")
        })
    }

    robot := gobot.NewRobot("buttonBot",
        []gobot.Connection{chipAdaptor},
        []gobot.Device{button},
        work,
    )

    robot.Start()
}
```

## Using the overlay manager

The C.H.I.P platform adapter includes an overlay manager for building
and installing device tree overlays for optional features PWM0 and SPI2.

```go
package main

import (
    "fmt"
    "log"
    "time"

    "gobot.io/x/gobot"
    "gobot.io/x/gobot/drivers/gpio"
    "gobot.io/x/gobot/platforms/chip"
)

func main() {
    chipAdaptor := chip.NewAdaptor()

    if err := chip.BuildAndInstallOverlays(); err != nil {
        log.Println(err)
        log.Fatal("Failed to build and install overlays")
    }

    if err := chip.LoadOverlay("PWM0"); err != nil {
        fmt.Printf("Failed to load overlay: %v\n", err)
        log.Fatal(err)
    }

    servo := gpio.NewServoDriver(chipAdaptor, "PWM0")

    work := func() {
        gobot.Every(10*time.Second, func() {
            err := servo.Move(0)
            if err != nil {
                fmt.Printf("Failed to move servo: %v\n", err)
                return
            }
            time.Sleep(1 * time.Second)
            err = servo.Move(180)
            time.Sleep(1 * time.Second)
        })
    }

    robot := gobot.NewRobot("servoBot",
        []gobot.Connection{chipAdaptor},
        []gobot.Device{servo},
        work,
    )

    robot.Start()
}

```

## How to Connect

### Compiling

Compile your Gobot program like this:

```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/chip_button.go
```

Then you can simply upload your program to the CHIP and execute it with

```bash
$ scp chip_button root@192.168.1.xx:
$ ssh -t root@192.168.1.xx "./chip_button"
```
