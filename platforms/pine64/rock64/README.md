# Pine ROCK64

The Pine ROCK64 is a single board SoC computer based on the Rockchip RK3328 arm64 processor. It has built-in GPIO and
I2C interfaces. SPI is most likely not usable (not tested), because in use by the SPI FLASH 128M memory chip.

For more info about the Pine ROCK64, go to [https://pine64.org/documentation/ROCK64/](https://pine64.org/documentation/ROCK64/).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Tested OS:

* [armbian](https://www.armbian.com/rock64/) with Debian

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
r := rock64.NewAdaptor()
led := gpio.NewLedDriver(r, "7")
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```sh
GOARCH=arm64 GOOS=linux go build -o output/ examples/rock64_blink.go
```

Once you have compiled your code, you can upload your program and execute it on the board from your workstation
using the `scp` and `ssh` commands like this:

```sh
scp rock64_blink <user>@192.168.1.xxx:~
ssh -t <user>@192.168.1.xxx "./rock64_blink"
```

## Troubleshooting

### I2C-0 overlay

With the armbian-config sometimes the overlays can not properly applied (different open Bugs). To ensure your overlay
is applied have a look into your /boot/boot.cmd and search for the name of the used shell variable(s). This name needs
to be used in your /boot/armbianEnv.txt.

```sh cat /boot/boot.cmd | grep overlay_file
for overlay_file in ${overlays}; do
	if load ${devtype} ${devnum}:${distro_bootpart} ${load_addr} ${prefix}dtb/rockchip/overlay/${overlay_prefix}-${overlay_file}.dtbo; then
		echo "Applying kernel provided DT overlay ${overlay_prefix}-${overlay_file}.dtbo"
for overlay_file in ${user_overlays}; do
	if load ${devtype} ${devnum}:${distro_bootpart} ${load_addr} ${prefix}overlay-user/${overlay_file}.dtbo; then
		echo "Applying user provided DT overlay ${overlay_file}.dtbo"
```

In the example above the variable is named `overlays`. So your /boot/armbianEnv.txt must contain this variable.

```sh cat /boot/armbianEnv.txt | grep overlay
overlay_prefix=rockchip
overlays=rk3328-i2c0
```

In some buggy versions the variable is named "fdt_overlays", just rename the variable in your "armbianEnv.txt" to match
the boot script.

As you can see in the boot script, the real file name is a concatenate of different variables `${overlay_prefix}-${overlay_file}.dtbo`.
This file must exist in the folder `${prefix}dtb/rockchip/overlay` (prefix="/boot/"). So for the i2c-0 overlay:

```sh ls -la /boot/dtb/rockchip/overlay/ | grep i2c0
-rw-r--r-- 1 root root   218 Nov 25 19:15 rockchip-rk3328-i2c0.dtbo
-rw-r--r-- 1 root root   223 Nov 25 19:15 rockchip-rk3568-hk-i2c0.dtbo
```

...means the entry in the armbianEnv.txt sould be set to "overlays=rk3328-i2c0".

The variable can contain a space separated list.

### PWM

There are 3 PWMs on the chip (pwm0, pwm1, pwm2). Unfortunately all pins are shared with the PMIC, so i2c-1 (pwm0, pwm1)
can not be deactivated, because it is mandatory for the i2c communication to PMIC address 0x18. Simply an activation of
pwm0 or pwm1 with an overlay leads to the Kernel can not be loaded anymore.