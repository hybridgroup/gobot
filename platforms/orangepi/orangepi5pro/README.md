# OrangePi 5 Pro

The OrangePi 5 Pro is a single board SoC computer based on the Rockchip RK3588S arm64 processor. It has built-in
GPIO, I2C, PWM, SPI and MIPI DSI interfaces.

For more info about the OrangePi 5 Pro, go to [http://www.orangepi.org/html/hardWare/computerAndMicrocontrollers/details/Orange-Pi-5-Pro.html](http://www.orangepi.org/html/hardWare/computerAndMicrocontrollers/details/Orange-Pi-5-Pro.html).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Tested OS:

* [armbian](https://www.armbian.com/orange-pi-5-pro/): "Armbian_community_25.5.0-trunk.4_Orangepi5pro_bookworm_vendor_6.1.99_minimal" (1-wire not working)
* [Debian server image](http://www.orangepi.org/html/hardWare/computerAndMicrocontrollers/service-and-support/Orange-Pi-5-Pro.html): "Orangepi5pro_1.0.4_debian_bookworm_server_linux6.1.43"

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
r := orangepi5pro.NewAdaptor()
led := gpio.NewLedDriver(r, "7")
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```sh
GOARCH=arm64 GOOS=linux go build -o output/ examples/orangepi5pro_blink.go
```

Once you have compiled your code, you can upload your program and execute it on the board from your workstation
using the `scp` and `ssh` commands like this:

```sh
scp output/orangepi5pro_blink <user>@192.168.1.xxx:~
ssh -t <user>@192.168.1.xxx "./orangepi5pro_blink"
```

## Troubleshooting
