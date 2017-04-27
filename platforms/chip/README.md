# C.H.I.P.

The [C.H.I.P.](http://www.getchip.com/) is a small, inexpensive ARM based single board computer, with many different IO interfaces available on the [pin headers](http://docs.getchip.com/#pin-headers).

For documentation about the C.H.I.P. platform click [here](http://docs.getchip.com/).

The [C.H.I.P. Pro](https://getchip.com/pages/chippro) is a version of C.H.I.P. intended for use in embedded product development. Here is info about the [C.H.I.P. Pro pin headers](https://docs.getchip.com/chip_pro.html#pin-descriptions).


## How to Install

We recommend updating to the latest Debian OS when using the C.H.I.P., however Gobot should also support older versions of the OS, should your application require this.

You would normally install Go and Gobot on your workstation. Once installed, cross compile your program on your workstation, transfer the final executable to your C.H.I.P and run the program on the C.H.I.P. itself as documented here.

```
go get -d -u gobot.io/x/gobot/...
```

### PWM support
Note that PWM might not be available in your kernel. In that case, you can install the required device tree overlay
from the command line using [Gort](https://gobot.io/x/gort) CLI commands on the C.H.I.P device.
Here are the steps:

Install the required patched device tree compiler as described in the [C.H.I.P docs](https://docs.getchip.com/dip.html#make-a-dtbo-device-tree-overlay-blob):
```
gort chip install dtc
```

Now, install the pwm overlay to activate pwm on the PWM0 pin:
```
gort chip install pwm
```

Reboot the device to make sure the init script loads the overlay on boot.


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

If you want to use the C.H.I.P. Pro, use the `NewProAdaptor()` function like this:

```go
chipProAdaptor := chip.NewProAdaptor()
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/chip_button.go
```

Once you have compiled your code, you can you can upload your program and execute it on the C.H.I.P. from your workstation using the `scp` and `ssh` commands like this:

```bash
$ scp chip_button root@192.168.1.xx:
$ ssh -t root@192.168.1.xx "./chip_button"
```
