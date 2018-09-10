# UP2

The UP2 Board aka "Up Squared" is a single board SoC computer based on the Intel Apollo Lake processor. It has built-in GPIO, PWM, SPI, and I2C interfaces.

For more info about the UP2 Board, go to [http://www.up-board.org/upsquared/](http://www.up-board.org/upsquared/).

## How to Install

### Update operating system and BIOS on UP2 board

We recommend updating to the latest Ubuntu 16.04 and BIOS v3.3 when using the UP2 board. To update your UP2 OS go to:

https://downloads.up-community.org/download/up-squared-iot-grove-development-kit-ubuntu-16-04-server-image/

To update your UP2 BIOS, go to:

https://downloads.up-community.org/download/up-squared-uefi-bios-v3-3/

Once your UP2 has been updated, you will need to provide permission to the `upsquared` user to access the GPIO or I2C subsystems on the board.

### Configuring GPIO on UP2 board

To access the GPIO subsystem, you will need to create a new group, add your user to the group, and then add a UDEV rule.

First, run the following commands on the board itself:

```
sudo groupadd gpio
sudo adduser upsquared gpio
```

Now, add a new UDEV rule to the UP2 board. Add the following text as a new UDEV rule named `/etc/udev/rules.d/99-gpio.rules`:

```
SUBSYSTEM=="gpio*", PROGRAM="/bin/sh -c '\
        chown -R root:gpiouser /sys/class/gpio && chmod -R 770 /sys/class/gpio;\
        chown -R root:gpiouser /sys/devices/virtual/gpio && chmod -R 770 /sys/devices/virtual/gpio;\
        chown -R root:gpiouser /sys$devpath && chmod -R 770 /sys$devpath\
'"
```

### Configuring built-in LEDs on UP2 board

To use the built-in LEDs you will need to create a new group, add your user to the group, and then add a UDEV rule.

First, run the following commands on the board itself:

```
sudo groupadd leds
sudo adduser upsquared leds
```

Now add the following text as a new UDEV rule named `/etc/udev/rules.d/99-leds.rules`:

```
SUBSYSTEM=="leds*", PROGRAM="/bin/sh -c '\
        chown -R root:leds /sys/class/leds && chmod -R 770 /sys/class/leds;\
        chown -R root:leds /sys/devices/platform/up-pinctrl/leds && chmod -R 770 /sys/devices/platform/up-pinctrl/leds;\
        chown -R root:leds /sys/devices/platform/AANT0F01:00/upboard-led.* && chmod -R 770 /sys/devices/platform/AANT0F01:00/upboard-led.*;\
'"
```

### Configuring I2C on UP2 board

To access the I2C subsystem, run the following command:

```
sudo usermod -aG i2c upsquared
```

You should reboot your UP2 board after making these changes for them to take effect.

**IMPORTANT NOTE REGARDING I2C:** 
The current UP2 firmware is not able to scan for I2C devices using the `i2cdetect` command line tool. If you run this tool, it will cause the I2C subsystem to malfunction until you reboot your system. That means at this time, do not use `i2cdetect` on the UP2 board.

### Local setup

You would normally install Go and Gobot on your local workstation. Once installed, cross compile your program on your workstation, transfer the final executable to your UP2, and run the program on the UP2 as documented below.

```
go get -d -u gobot.io/x/gobot/...
```

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
r := up2.NewAdaptor()
led := gpio.NewLedDriver(r, "13")
```

You can also use the values `up2.LEDRed`, `up2.LEDBlue`, `up2.LEDGreen`, and `up2.LEDYellow` as pin reference to access the 4 built-in LEDs. For example:

```go
r := up2.NewAdaptor()
led := gpio.NewLedDriver(r, up2.LEDRed)
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```bash
$ GOARCH=amd64 GOOS=linux go build examples/up2_blink.go
```

Once you have compiled your code, you can you can upload your program and execute it on the UP2 from your workstation using the `scp` and `ssh` commands like this:

```bash
$ scp up2_blink upsquared@192.168.1.xxx:/home/upsquared/
$ ssh -t upsquared@192.168.1.xxx "./up2_blink"
```
