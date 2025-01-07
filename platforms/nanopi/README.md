# NanoPi Boards

The FriendlyARM NanoPi Boards are single board SoC computers with different hardware design. It has built-in GPIO, PWM,
SPI, and I2C interfaces.

For more info about the NanoPi Boards, go to [https://wiki.friendlyelec.com/wiki/index.php/Main_Page](https://wiki.friendlyelec.com/wiki/index.php/Main_Page).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Tested OS:

* [armbian](https://www.armbian.com/nanopi-neo/) with Debian or Ubuntu

### System access and configuration basics

Please follow the installation instructions for the chosen OS.

### Enabling hardware drivers

Please follow the configuration instructions for the chosen OS.

E.g. for armbian:

```sh
sudo armbian-config
```

After configuration was changed, an reboot is necessary.

```sh
sudo reboot
```

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
r := nanopi.NewNeoAdaptor()
led := gpio.NewLedDriver(r, "7")
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```sh
GOARM=7 GOARCH=arm GOOS=linux go build -o output/ examples/nanopi_blink.go
```

Once you have compiled your code, you can upload your program and execute it on the board from your workstation
using the `scp` and `ssh` commands like this:

```sh
scp output/nanopi_blink nan@nanopineo:~
ssh -t nan@nanopineo "./nanopi_blink"
```

## GPIO's

At least for NEO, nearly all 14 GPIO's supports the advanced pin options "bias", "drive", "debounce" and "edge detection".

> Configure of edge detection will cause an initial event. GPIO header pins 19, 21, 23, 24 - do NOT support "debounce" and
> "edge detection" for NEO. Using unsupported options leads to reconfigure errors with text "no such device or address".

## PWM

A single PWM output is available at UART0-RX (UART_RXD0, internal PA5). So the UART0 needs to be disabled. The sunxi-
overlay (e.g. for armbian) disables the UART0 and the kernel console at ttyS0. The related kernel module needs to be
loaded: `sudo modprobe pwm-sun4i`. The default frequency is 100Hz.

## I2C

The default bus number is set to 0, which is connected to header pins 3 (PA12-SDA) and 5 (PA11-SCL). At least for NEO
rev.1.4 it is possible to activate bus 1, which is connected to "USB/Audio/IR" header pins 9 "PCM0_CLK/I2S0_BCK"
(PA19-SDA) and 8 "PCM0_SYNC/I2S0_LRC" (PA18-SCL). Armbian allows to activate bus 2 (PE12-SCL, PE13-SDA), which pins are
not wired for NEO and NEO2, but we do not block it at adaptor side.

## SPI

There is a known issue on [armbian](https://forum.armbian.com/topic/20033-51525-breaks-spi-on-nanopi-neo-and-does-not-create-devspidev00/)
for later Kernels.
