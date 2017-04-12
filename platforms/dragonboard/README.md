# DragonBoard™ 410c

The [DragonBoard 410c](http://www.96boards.org/product/dragonboard410c/), a product of Arrow Electronics, is the development board based on the mid-tier Qualcomm® Snapdragon™ 410E processor. It features advanced processing power, Wi-Fi, Bluetooth connectivity, and GPS, all packed into a board the size of a credit card.

## How to Install

Make sure you are using the latest Linaro Debian image. Both AArch32 and AArch64 work™ though you should stick to 64bit as OS internals may be different and aren't tested.

You would normally install Go and Gobot on your workstation. Once installed, cross compile your program on your workstation, transfer the final executable to your DragonBoard and run the program on the DragonBoard itself as documented here.

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself. See [here](https://www.96boards.org/db410c-getting-started/HardwareDocs/HWUserManual.md/).

```go
package main

import (
    "fmt"

    "gobot.io/x/gobot"
    "gobot.io/x/gobot/drivers/gpio"
    "gobot.io/x/gobot/platforms/dragonboard"
)

func main() {
    dragonAdaptor := dragonboard.NewAdaptor()
    button := gpio.NewButtonDriver(dragonAdaptor, "GPIO_A")

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

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```bash
$ GOARCH=arm64 GOOS=linux go build examples/dragon_button.go
```

Once you have compiled your code, you can you can upload your program and execute it on the DragonBoard from your workstation using the `scp` and `ssh` commands like this:

```bash
$ scp dragon_button root@192.168.1.xx:
$ ssh -t root@192.168.1.xx "./dragon_button"
```
