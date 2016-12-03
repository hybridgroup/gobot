# Beaglebone

The BeagleBone is an ARM based single board computer, with many different GPIO interfaces built in.

For more info about the BeagleBone platform click [here](http://beagleboard.org/Products/BeagleBone+Black).

## How to Install

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/beaglebone
```

## How to Use

```go
package main

import (
	"time"

	"github.com/hybridgroup/gobot"
	"github.com/hybridgroup/gobot/drivers/gpio"
	"github.com/hybridgroup/gobot/platforms/beaglebone"
)

func main() {
	beagleboneAdaptor := beaglebone.NewAdaptor()
	led := gpio.NewLedDriver(beagleboneAdaptor, "P9_12")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{beagleboneAdaptor},
		[]gobot.Device{led},
		work,
	)

	robot.Start()
}
```

## How to Connect

### Compiling

Simply compile your Gobot program like this:

```bash
$ GOARCH=arm GOOS=linux go build examples/beaglebone_blink.go
```

If you are running the official Debian Linux through the usb->ethernet connection, or are connected to the board using WiFi, you can simply upload your program and execute it with the `scp` command like this:

```bash
$ scp beaglebone_blink root@192.168.7.2:/home/root/
$ ssh -t root@192.168.7.2 "./beaglebone_blink"
```

### Updating your board to the latest OS

We recommend updating your BeagleBone to the latest Debian OS. It is very easy to do this using the Etcher (https://etcher.io/) utility program.

First, download the latest BeagleBone OS from http://beagleboard.org/latest-images

Now, use Etcher to create an SD card with the OS image you have downloaded.

Once you have created the SD card, boot your BeagleBone using the new image as follows:

- Insert SD card into your (powered-down) board, hold down the USER/BOOT button (if using Black) and apply power, either by the USB cable or 5V adapter.

- If all you want to do it boot once from the SD card, it should now be booting.

- If using BeagleBone Black and desire to write the image to your on-board eMMC, you'll need to follow the instructions at http://elinux.org/Beagleboard:BeagleBoneBlack_Debian#Flashing_eMMC. When the flashing is complete, all 4 USRx LEDs will be steady on or off. The latest Debian flasher images automatically power down the board upon completion. This can take up to 45 minutes. Power-down your board, remove the SD card and apply power again to be complete.

These instructions come from the Beagleboard web site's "Getting Started" page located here:

http://beagleboard.org/getting-started
