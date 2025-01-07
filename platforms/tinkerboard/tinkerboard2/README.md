# Tinker Board 2

The ASUS Tinker Board 2 is a single board SoC computer based on the Rockchip RK3399 processor (arm64). It has built-in
GPIO, I2C, PWM, SPI, 1-Wire, MIPI CSI and MIPI DSI interfaces.

For more info about the Tinker Board, go to [https://tinker-board.asus.com/series/tinker-board-2.html/](https://tinker-board.asus.com/series/tinker-board-2.html).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Tested OS:

* [armbian](https://www.armbian.com/tinkerboard-2/) with Debian

### System access and configuration basics

Use `sudo armbian-config` or see description for [Tinker Board](../README.md).

### Enabling hardware drivers

See description for [Tinker Board](../README.md).

### Enabling GPIO pins

See description for [Tinker Board](../README.md).

### Enabling I2C

See description for [Tinker Board](../README.md).

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
r := tinkerboard2.NewAdaptor()
led := gpio.NewLedDriver(r, "7")
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```sh
GOARCH=arm64 GOOS=linux go build -o output/ examples/tinkerboard2_yl40.go
```

Once you have compiled your code, you can upload your program and execute it on the Tinker Board 2 from your workstation
using the `scp` and `ssh` commands like this:

```sh
scp output/tinkerboard2_yl40 <user>@192.168.1.xxx:~
ssh -t <user>@192.168.1.xxx "./tinkerboard2_yl40"
```

## Troubleshooting

### I2C

The version "Armbian_community 25.2.0-trunk.124 bookworm" contains by default the overlay for i2c7 (header pins 27, 28)
and i2c8 (connection on DSI_1, DSI_2). An overlay for i2c6 (header pins 3, 5) needs to be created manually based on one
of the other overlays.

### PWM

#### Investigate state

```sh
# ls -la /sys/class/pwm/
ls -la /sys/class/pwm/
total 0
drwxr-xr-x  2 root root 0 Jan 18  2013 .
drwxr-xr-x 77 root root 0 Jan 18  2013 ..
lrwxrwxrwx  1 root root 0 Jan 18  2013 pwmchip0 -> ../../devices/platform/ff420020.pwm/pwm/pwmchip0
```

looking for one of the following items in the path:
ff420000 => pwm0, pin32
ff420010 => pwm1, pin33
ff420020 => pwm2, already activated, but internally, not usable
ff420030 => pwm3, pin26

#### Activate

The version "Armbian_community 25.2.0-trunk.124 bookworm" contains no overlay for pwm0, pwm1 or pwm3. This needs to be
created based on another overlay, e.g. `rockchip-rk3568-hk-pwm1.dtbo`. After activation of your preferred pwmchip,
proceed like described for [Tinker Board](../README).
