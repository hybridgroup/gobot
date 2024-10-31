# 1-wire bus

This document describes some basics for developers. This is useful to understand programming in gobot's [1-wire driver](./onewiredevice_sysfs.go).

## 1-wire with sysfs

If the 1-wire bus is enabled on the board, the bus data can be connected to one pin of the used platform. The enabling
activates the Kernel drivers for common devices (family drivers), which are than mapped to the sysfs, see
<https://www.kernel.org/doc/Documentation/w1/w1.generic>.

## Check available 1-wire devices

Example for Tinkerboard (RK3288) with armbian and 3 connected temperature sensors DS18B20:

```sh
ls -la /sys/bus/w1/devices/
insgesamt 0
drwxr-xr-x 2 root root 0 29. Okt 08:58 .
drwxr-xr-x 4 root root 0 29. Okt 08:58 ..
lrwxrwxrwx 1 root root 0 31. Okt 07:55 28-072261452f18 -> ../../../devices/w1_bus_master1/28-072261452f18
lrwxrwxrwx 1 root root 0 31. Okt 07:55 28-08225482b0de -> ../../../devices/w1_bus_master1/28-08225482b0de
lrwxrwxrwx 1 root root 0 31. Okt 07:55 28-1e40710a6461 -> ../../../devices/w1_bus_master1/28-1e40710a6461
lrwxrwxrwx 1 root root 0 29. Okt 08:58 w1_bus_master1 -> ../../../devices/w1_bus_master1
```

Within a device folder different files are available for typical access.

```sh
ls -la /sys/bus/w1/devices/28-072261452f18/
insgesamt 0
drwxr-xr-x 4 root root    0 29. Okt 08:58 .
drwxr-xr-x 6 root root    0 29. Okt 08:58 ..
-rw-r--r-- 1 root root 4096 31. Okt 07:57 alarms
-rw-r--r-- 1 root root 4096 31. Okt 07:57 conv_time
lrwxrwxrwx 1 root root    0 31. Okt 07:57 driver -> ../../../bus/w1/drivers/w1_slave_driver
--w------- 1 root root 4096 31. Okt 07:57 eeprom_cmd
-r--r--r-- 1 root root 4096 31. Okt 07:57 ext_power
-rw-r--r-- 1 root root 4096 31. Okt 07:57 features
drwxr-xr-x 3 root root    0 29. Okt 08:58 hwmon
-r--r--r-- 1 root root 4096 31. Okt 07:57 id
-r--r--r-- 1 root root 4096 31. Okt 07:57 name
drwxr-xr-x 2 root root    0 31. Okt 07:13 power
-rw-r--r-- 1 root root 4096 31. Okt 11:10 resolution
lrwxrwxrwx 1 root root    0 29. Okt 08:58 subsystem -> ../../../bus/w1
-r--r--r-- 1 root root 4096 31. Okt 07:57 temperature
-rw-r--r-- 1 root root 4096 29. Okt 08:58 uevent
-rw-r--r-- 1 root root 4096 31. Okt 07:57 w1_slave
```

This files depends on the family driver.

## Different access levels and modes

Currently gobot supports only direct access to the devices in automatic search mode of the controller device. The
implementation is similar to the sysfs access of the analog pin driver.

E.g. if the cyclic device search should be avoided, the access to the controller device is needed, see Kernel
documentation. If this will be implemented in the future, have in mind that more than one controller devices are
possible. The gobot's 1-wire architecture can be changed then similar to SPI or I2C.

## Troubleshooting

If something is not working, please check this points:

Is the correct gpio used: `cat /sys/kernel/debug/gpio`
Does the base path exist: `ls /sys/bus/w1/`
Is the Kernel module loaded: `lsmod | grep wire`
Is the onewire support activated in the device tree: `dtc -I fs -O dts /sys/firmware/devicetree/base | grep onewire`
Is there an according overlay on the system: `locate w1-gpio`
Is the overlay loading configured: `cat /boot/armbianEnv.txt`
Is the content of the overlay correct: `dtc -I dtb -O dts -o <name>-w1-gpio.dtbo.dts /boot/dtb/overlay/<name>-w1-gpio.dtbo`
