# Raspberry Pi

The Raspberry Pi is an inexpensive and popular ARM based single board computer with digital & PWM GPIO, and i2c interfaces
built in.

The Gobot adaptor for the Raspberry Pi should support all of the various Raspberry Pi boards such as the
Raspberry Pi 4 Model B, Raspberry Pi 3 Model B, Raspberry Pi 2 Model B, Raspberry Pi 1 Model A+, Raspberry Pi Zero,
and Raspberry Pi Zero W.

For more info about the Raspberry Pi platform, click [here](http://www.raspberrypi.org/).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

We recommend updating to the latest Raspian Jessie OS when using the Raspberry Pi, however Gobot should also support
older versions of the OS, should your application require this.

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
package main

import (
"time"

"gobot.io/x/gobot/v2"
"gobot.io/x/gobot/v2/drivers/gpio"
"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
  r := raspi.NewAdaptor()
  led := gpio.NewLedDriver(r, "7")

  work := func() {
    gobot.Every(1*time.Second, func() {
      if err := led.Toggle(); err != nil {
				fmt.Println(err)
			}
    })
  }

  robot := gobot.NewRobot("blinkBot",
    []gobot.Connection{r},
    []gobot.Device{led},
    work,
  )

  if err := robot.Start(); err != nil {
		panic(err)
	}
}
```

## Compiling

Compile your Gobot program on your workstation like this:

```sh
GOARM=6 GOARCH=arm GOOS=linux go build examples/raspi_blink.go
```

Use the following `GOARM` values to compile depending on which model Raspberry Pi you are using:

`GOARM=6` (Raspberry Pi A, A+, B, B+, Zero)
`GOARM=7` (Raspberry Pi 2, 3)

Once you have compiled your code, you can upload your program and execute it on the Raspberry Pi from your workstation
using the `scp` and `ssh` commands like this:

```sh
scp raspi_blink pi@192.168.1.xxx:/home/pi/
ssh -t pi@192.168.1.xxx "./raspi_blink"
```

## Enabling PWM output on GPIO pins

### Using Linux Kernel sysfs implementation

The PWM needs to be enabled in the device tree. Please read `/boot/overlays/README` of your device. Usually "pwm0" can
be activated for all raspi variants with a line `dtoverlay=pwm,pin=18,func=2` added to `/boot/config.txt`. The number
relates to "GPIO18", not the header number, which is "12" in this case.

Now the pin can be used with gobot by the pwm channel name, e.g. for our example above:

```go
...
// create the adaptor with a 50Hz default frequency for usage with servos
a := NewAdaptor(adaptors.WithPWMDefaultPeriod(20000000))
// move servo connected with header pin 12 to 90°
a.ServoWrite("pwm0", 90)
...
```

> If the activation fails or something strange happen, maybe the audio driver conflicts with the PWM. Please deactivate
> the audio device tree overlay in `/boot/config.txt` to avoid conflicts.

### Using pi-blaster

For support PWM on all pins, you may use a program called pi-blaster. You can follow the instructions for install in
the pi-blaster repo here: <https://github.com/sarfata/pi-blaster>

For using a PWM for servo, the default 100Hz period needs to be adjusted to 50Hz in the source code of the driver.
Please refer to <https://github.com/sarfata/pi-blaster#how-to-adjust-the-frequency-and-the-resolution-of-the-pwm>.

It is not possible to change the period from gobot side.

Now the pin can be used with gobot by the header number, e.g.:

```go
...
// create the adaptor with usage of pi-blaster instead of default sysfs, 50Hz default is given for calculate
// duty cycle for servos but will not change everything for the pi-blaster driver, see description above
a := NewAdaptor(adaptors.WithPWMUsePiBlaster(), adaptors.WithPWMDefaultPeriod(20000000))
// move servo to 90°
a.ServoWrite("11", 90)
// this will not work like expected, see description
a.SetPeriod("11", 20000000)
...
```
