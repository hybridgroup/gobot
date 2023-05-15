# system

This document describes some basics for developers.

> The gobot package name is a little bit misleading, because also contains `/dev/i2c-*` usage with syscall for i2c
> interface and the file system mock for unit testing.

## Change available features on TinkerOS (Tinkerboard)

Open the file /boot/config.txt and add lines/remove comment sign. Then save the file, close and reboot.

### For PWM on TinkerOS (Tinkerboard)

intf:pwm2=on
intf:pwm3=on

### For i2c on TinkerOS (Tinkerboard)

intf:i2c1=on
intf:i2c4=on

Device tree overlays can also be provided, e.g. "overlay=i2c-100kbit-tinker-overlay".

Create a dts-file "i2c1-100kbit-tinker-overlay.dts" with this content:

```c
// Definitions for i2c1 with 100 kbit
/dts-v1/;
/plugin/;

/ {
  compatible = "rockchip,rk3288-evb-rk808-linux", "rockchip,rk3288";

  fragment@0 {
    target = <&i2c1>;
    __overlay__ {
      status = "okay";
      clock-frequency = <100000>;
    };
  };
};
```

Compile the tree overlay:

```sh
dtc -@ -I dts -O dtb -o i2c1-100kbit-tinker-overlay.dtbo i2c1-100kbit-tinker-overlay.dts
```

Copy the file to "/boot/overlays/i2c-100kbit-tinker-overlay.dtbo"

> Package "device-tree-compiler" needs to be installed for "dtc" command.

## Change available features on raspi

Start "raspi-config" and go to the menu "interface options". After change something save and reboot the device.

## Change available features on armbian

Install the armbian-config utility (howto is shown after each login) and start it.

### For i2c on armbian

Open the related menu entry and activate/deactivate the needed features.

### For PWM on armbian

Go to the device tree menu - this opens the device tree for editing.  
Looking for the pwm entries. There are some marked with "disabled". Change the entry to "okay" for your needed pwm.

Example for Tinkerboard:

ff680020 => pwm2, pin33
ff680030 => pwm3, pin32

## Next steps for developers

* test [gpio](GPIO.md)
* test [pwm](PWM.md)
* background information for [i2c](I2C.md) in gobot

## Links

* <https://www.digi.com/resources/documentation/digidocs/embedded/dey/3.0/cc8x/bsp_r_create-dt-overlays>
