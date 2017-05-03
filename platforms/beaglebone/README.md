# Beaglebone

The BeagleBone is an ARM based single board computer, with lots of GPIO, I2C, and analog interfaces built in.

The Gobot adaptor for the BeagleBone should support all of the various BeagleBone boards such as the BeagleBone Black, SeeedStudio BeagleBone Green, SeeedStudio BeagleBone Green Wireless, and others that use the latest Debian and standard "Cape Manager" interfaces.

For more info about the BeagleBone platform go to  [http://beagleboard.org/getting-started](http://beagleboard.org/getting-started).


## How to Install

We recommend updating to the latest Debian Jessie OS when using the BeagleBone. The current Gobot only supports 4.x versions of the OS. If you need support for older versions of the OS, you will need to use Gobot v1.4.

You would normally install Go and Gobot on your workstation. Once installed, cross compile your program on your workstation, transfer the final executable to your BeagleBone, and run the program on the BeagleBone itself as documented here.

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

Gobot also has support for the four built-in LEDs on the BeagleBone Black, by referring to them as `usr0`, `usr1`, `usr2`, and `usr3`.

```go
package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/beaglebone"
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

Compile your Gobot program on your workstation like this:

```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/beaglebone_blink.go
```

Once you have compiled your code, you can you can upload your program and execute it on the BeagleBone from your workstation using the `scp` and `ssh` commands like this:

```bash
$ scp beaglebone_blink root@192.168.7.2:/home/root/
$ ssh -t root@192.168.7.2 "./beaglebone_blink"
```

In order to run the preceeding commands, you must be running the official Debian Linux through the usb->ethernet connection, or be connected to the board using WiFi.

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
