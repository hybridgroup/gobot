# Radxa Rock Pi

The [Radxa Rock Pi board series](https://wiki.radxa.com/Rock4/getting_started) are clones of the popular Raspberry Pi Single Board Computers (SBCs) with GPIO/PWM/I2C functionalities built-in.

The Gobot adaptor is currently compatible with:

- Rock Pi 4
- Rock Pi 4C+

With the possibility to expand its compatibility into past and future models.

Check out the output of `cat /proc/device-tree/model` to see which model you have if you're not sure. The 4C+ model has a Rockchip 3399_T SoC, while the regular 4 has the 3399. Both are similar, but have slightly different GPIO pin configurations.

## How to Install

Make sure you've installed an official Linux image from Radxa with working drivers. See the [ROCK 4 Installation Wiki](https://wiki.radxa.com/Rock4/install) for your SBC setup.

> Some versions or Armbian ISOs do not detect the newer SoC chips.

As for your Gobot development, treat this as a regular Go package. It can be cross-compiled and copied over, or simply compiled on the SBC itself (tested and working with go 1.15.15 on linux/arm64, RockPi4C+).

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

That is, follow the **colored Pin# numbers** in the middle in the [GPIO mapping Wiki](https://wiki.radxa.com/Rock4/hardware/gpio). These have been translated for you into their corresponding underlying GPIO numbers.

```go
package main

import (
        "time"

        "gobot.io/x/gobot/v2"
        "gobot.io/x/gobot/v2/drivers/gpio"
        "gobot.io/x/gobot/v2/platforms/radxa/rockpi"
)

func main() {
        r := rockpi.NewAdaptor()
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

If you want to use I2C, RockPi4 offers three I2C buses: I2C2 (pins 27, 28), I2C6 (pins 29, 31) and I2C7 (pins 3, 5) of which I2C7 is the default.
Changing this is a matter of passing the right bus number using `i2c.WithBus`:

```go
a := rockpi.NewAdaptor()
adc2 := i2c.NewADS1115Driver(a, i2c.WithBus(6))
```

There are mapped to `/dev/i2c-[bus]`, just like the Gobot raspi implementation.

### What's currently supported?

- General digital pin GPIO access works, but not through the new character device driver.
- The I2C buses 2, 6, 7 work
- SPI buses 1, 2 work

PWM interaction is currently not yet supported. 

Please see the [official Radxa Rock Pi documentation](https://wiki.radxa.com/Rockpi4/dev/libmraa) on how `librmaa` can be utilized in combination with this module, and how these work.

### Compiling

Compile your Gobot program on your workstation like this:

```bash
$ GOARCH=arm64 GOOS=linux go build examples/rockpi_blink.go
```

Rock Pi 4s are ARM64 machines.

