# FriendlyELEC NanoPC-T6

The FriendlyELEC NanoPC-T6 is a single board SoC computer based on the Rockchip RK3588 arm64 processor. It has built-in
GPIO, I2C, PWM, SPI, 1-Wire, MIPI CSI and MIPI DSI interfaces.

For more info about the FriendlyELEC NanoPC-T6, go to [https://wiki.friendlyelec.com/wiki/index.php/NanoPC-T6](https://wiki.friendlyelec.com/wiki/index.php/NanoPC-T6).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Tested OS:

* [armbian](https://www.armbian.com/nanopct6/) with "Armbian 24.11.1 Bookworm Minimal / IOT"

## Configuration steps for the OS

### System access and configuration basics

Please follow the instructions of the OS provider. A ssh access is used in this guide.

```sh
ssh <user>@192.168.1.xxx
```

### Enabling hardware drivers

Not all drivers are enabled by default. You can have a look at the configuration file, to find out what is enabled at
your system:

```sh
cat /boot/armbianEnv.txt
```

```sh
sudo apt install armbian-config
sudo armbian-config
```

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
r := nanopct6.NewAdaptor()
led := gpio.NewLedDriver(r, "7")
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```sh
GOARCH=arm64 GOOS=linux go build -o output/ examples/nanopct6_blink.go
```

Once you have compiled your code, you can upload your program and execute it on the board from your workstation
using the `scp` and `ssh` commands like this:

```sh
scp nanopct6_blink <user>@192.168.1.xxx:~
ssh -t <user>@192.168.1.xxx "./nanopct6_blink"
```

## Troubleshooting
