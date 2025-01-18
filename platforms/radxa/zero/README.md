# Radxa Zero

The Radxa Zero is a single board SoC computer based on the Amlogic S905Y2 arm64 processor. It has built-in
GPIO, I2C, PWM, SPI, 1-Wire and ADC interfaces.

For more info about the Radxa Zero, go to [https://docs.radxa.com/en/zero/zero](https://docs.radxa.com/en/zero/zero).

## How to Install

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

Tested OS:

* [dietPi](https://dietpi.com/downloads/images/DietPi_RadxaZero-ARMv8-Bookworm.img.xz) a minimal image with Debian Bookworm

## Configuration steps for the OS

### WLAN access

There is no LAN network interface on the board, but WLAN is sufficient for our needs. After copy over the image to your
SD card you can modify the file dietpi.txt before plug in the card to [make it work on first boot](https://dietpi.com/docs/usage/#how-to-do-an-automatic-base-installation-at-first-boot-dietpi-automation).

```txt adjust dietpi.txt to your needs
AUTO_SETUP_NET_WIFI_ENABLED=1
AUTO_SETUP_NET_WIFI_COUNTRY_CODE=DE
```

Afterwards the WiFi login data needs to be provided (unencrypted file, will be removed after first boot):

```txt adjust dietpi-wifi.txt
aWIFI_SSID[0]="<name of your network>"
aWIFI_KEY[0]="<password>"
```

First login needs to be done with "root" or "dietpi", and you will start with a configuration procedure.

```sh
ssh root@DietPi
```

### System access and configuration basics

Please follow the instructions of the OS provider. A ssh access for (WLAN with dietpi user) is used in this guide.

```sh
ssh dietpi@DietPi
```

### Enabling hardware drivers in general

Not all drivers are enabled by default. You can have a look at the configuration file, to find out what is enabled at
your system:

```sh
cat /boot/dietpiEnv.txt
```

Missing interfaces needs to be enabled by DT-overlays (drop "meson" prefix and file extension).

```sh list available overlays
ls /boot/dtb/amlogic/overlay/
```

Please read [this GPIO page](https://wiki.radxa.com/Zero/hardware/gpio.) for meaning of the different part of file name (e.g. "ao" vs. "ee").
The [page about overlays](https://wiki.radxa.com/Device-tree-overlays#Meson_G12A_Available_Overlay_.28Radxa_Zero.29) will help you in addition to choose the right one.

### enable I2C

|  device  |SDA|SCL|DT overlay file name|
|----------|---|---|--------------------|
|/dev/i2c-1| 16| 13|meson-g12a-radxa-zero-i2c-ee-m1-gpiox-10-gpiox-11.dtbo|
|/dev/i2c-1| 24| 23|meson-g12a-radxa-zero-i2c-ee-m1-gpioh-6-gpioh-7.dtbo|
|/dev/i2c-3|  3|  5|meson-g12a-radxa-zero-i2c-ee-m3-gpioa-14-gpioa-15.dtbo|
|/dev/i2c-4|  7| 11|meson-g12a-radxa-zero-i2c-ao-m0-gpioao-2-gpioao-3.dtbo|


```sh /boot/dietpiEnv.txt example for i2c-1
...
overlays=g12a-radxa-zero-i2c-ee-m1-gpioh-6-gpioh-7
...
```

>The I2C device "/dev/i2c-3" was already enabled on dietPi after setup.

### enable SPI

```sh /boot/dietpiEnv.txt for SPI
...
overlays=g12a-radxa-zero-spi-spidev
...
```

>Most likely the overlay is currently defective - it contains "armbian" and only "disabled" spi devices.

### enable PWM

|pin     |  symbol |          DT path          |        driver      |DT overlay file name|
|--------|---------|---------------------------|--------------------|--------------------|
|32, PWMAO_C|pwm_AO_cd|/soc/bus@ff800000/pwm@2000 |meson-g12a-ao-pwm-cd|on by default|
|40, PWMAO_A|pwm_AO_ab|/soc/bus@ff800000/pwm@7000 |meson-g12a-ao-pwm-ab|meson-g12a-radxa-zero-pwmao-a-on-gpioao-11.dtbo|
|18, PWM_C  |pwm_cd   |/soc/bus@ffd00000/pwm@1a000|meson-g12a-ee-pwm   |meson-g12a-radxa-zero-pwm-c-on-gpiox-8.dtbo|
|21, PWM_F  |pwm_ef   |/soc/bus@ffd00000/pwm@19000|meson-g12a-ee-pwm   |on by default|

>PWMAO_B (channel 1 of pwm_AO_ab) and PWM_D (channel 1 of pwm_cd) not wired. PWMAO_D (channel 1 of pwm_AO_cd) in use by
>"regulator-vddcpu". PWM_E (channel 0 of pwm_ef) in use by wifi32k. PWMAO_C and PWM_F not really working, see
>troubleshooting section.

### enable 1-wire

The contained overlays maybe not working, because compatible with gxbb (Amlogic Meson S905). At least for my Zero V1.51
g12a (Amlogic Meson S905X2 and above) it does not work.

```sh /boot/dietpiEnv.txt for 1-wire
dtc -I dtb -O dts /boot/dtb/amlogic/overlay/meson-w1AB-gpio.dtbo | grep amlogic
...
compatible = "amlogic,meson-gxbb";
```

So create your own overlay as `/boot/overlay_user/meson-g12a-w1-gpioao-3.dts`...

```sh
/dts-v1/;
/plugin/;

/ {
	compatible = "radxa,zero", "amlogic,g12a";

	fragment@0 {
		target-path = "/";

		__overlay__ {
			w1: onewire {
				compatible = "w1-gpio";
				pinctrl-names = "default";
				/* GPIOAO_3=0x03, GPIOC_7=0x30 */
				/* GPIO_ACTIVE_HIGH=0, GPIO_ACTIVE_LOW=1 */
				/* GPIO_SINGLE_ENDED=2, GPIO_LINE_OPEN_DRAIN=4 */
				gpios = <0xffffffff 0x03 0x06>;
				status = "okay";
				phandle = <0x01>;
			};
		};
	};

	__fixups__ {
		/* gpio_ao or gpio for GPIOC_7 */
		gpio_ao = "/fragment@0/__overlay__/onewire:gpios:0";
	};
};
```

...compile it

```sh
dtc -@ -O dtb -b 0 -o meson-g12a-w1-gpioao-3.dtbo meson-g12a-w1-gpioao-3.dts
```

... and add it as follows:

```sh /boot/dietpiEnv.txt for 1-wire
...
user_overlays=<existing names> meson-g12a-w1-gpioao-3
...
```

## enable SAR ADC

The 12-bit ADC is enabled by default. The voltage range is 0..1.8V. Raw values can be read with pin 15 (channel 1) or
pin 26 (channel 2). Additionally some internal values can be accessed, e.g. gnd and 1/4 vdd. For debugging purposes more
information is provided, e.g. each item provides a label:

```sh
cat /sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/in_voltage9_label
gnd
cat /sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/in_voltage10_label
0.25vdd
```
```sh
cat /sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/calibbias
-4
cat /sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/calibscale
1.002447
cat /sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/in_voltage_scale
0.439453125
cat /sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/in_voltage9_raw
0
cat /sys/bus/platform/drivers/meson-saradc/ff809000.adc/iio:device0/in_voltage10_raw
1023
```

The `in_voltage_scale` (e.g. 0.439453125) is for calculation of "value = (raw + offset) * scale", value in millivolts.
The `calibbias` is the offset and `calibscale` the scale which is used for internal calibration, so we get an output
of 0..4095 for 0..1.8V input.

>The channel 10 is a 1/4 of full range, but if the channel 2 is in saturation state (above 1.8V = 4095), it becomes more
>and more wrong.

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
r := zero.NewAdaptor()
led := gpio.NewLedDriver(r, "7")
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```sh
GOARCH=arm64 GOOS=linux go build -o output/ examples/zero_blink.go
```

Once you have compiled your code, you can upload your program and execute it on the board from your workstation
using the `scp` and `ssh` commands like this:

```sh
scp zero_blink dietpi@DietPi:~
ssh -t dietpi@DietPi "./zero_blink"
```

## Troubleshooting

### scp fails

"bash: line 1: /usr/lib/sftp-server: No such file or directory"

The dietPi has only a limited package set, so sftp-server is missing.

```sh
sudo apt install openssh-sftp-server
```

### GPIO pin3, pin5 and pin7 not working like expected

e.g. for pin3:
`cdev.Export(): cdev.reconfigure(gpiochip0-63)-c.RequestLine(63, [0 2000000000 2]): invalid argument`

The pin does not support `adaptors.WithGpioDebounce(inPinNum, debounceTime)`.

`cdev.Export(): cdev.reconfigure(gpiochip1-0)-c.RequestLine(0, [0 2]): invalid argument`

The pin 8 is configured for UART.

>Some pins have low power or have a strong pullup resistor, so the expected voltage drop is maybe not possible.
>Pins 7-27 and 11-28 are bridged.

### PWMAO_C and PWM_F not really working

If you run `cat /sys/kernel/debug/pwm` when those PWMs are used you can see, that it is working "internally". Most
likely the pin itself is not enabled by DT on PWM usage. Currently there is no fix provided for that.
