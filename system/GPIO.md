# GPIOs

This document describes some basics for developers. This is useful to understand programming in gobot's [digital pin driver](./digitalpin_access.go).

## GPIOs with sysfs

> Kernel SYSFS ABI is deprecated since Linux 4.8, see <https://www.kernel.org/doc/html/latest/admin-guide/gpio/sysfs.html>.
> For GPIO's we still use the Kernel SYSFS ABI.

## GPIOs with character devices

This document provides some test possibilities by using the new character device feature since Kernel 4.8. Please check
your Kernel version before using this. Install of "gpiod" is mandatory for the tests.

```sh
uname -a
Linux raspi 5.15.61+ #1579 Fri Aug 26 11:08:59 BST 2022 armv6l GNU/Linux

sudo apt install gpiod
```

> For work on character device user space drivers, please refer to our [issue #775](https://github.com/hybridgroup/gobot/issues/775).

## Check available GPIO banks

Example for Tinkerboard (RK3288) with TinkerOS:

```sh
ls -la /sys/class/gpio/
total 0
drwxr-xr-x  2 root root    0 Nov 12 06:53 .
drwxr-xr-x 66 root root    0 Feb 14  2019 ..
--w-------  1 root root 4096 Feb 14  2019 export
lrwxrwxrwx  1 root root    0 Feb 14  2019 gpiochip0 -> ../../devices/platform/pinctrl/gpio/gpiochip0
lrwxrwxrwx  1 root root    0 Feb 14  2019 gpiochip120 -> ../../devices/platform/pinctrl/gpio/gpiochip120
lrwxrwxrwx  1 root root    0 Feb 14  2019 gpiochip152 -> ../../devices/platform/pinctrl/gpio/gpiochip152
lrwxrwxrwx  1 root root    0 Feb 14  2019 gpiochip184 -> ../../devices/platform/pinctrl/gpio/gpiochip184
lrwxrwxrwx  1 root root    0 Feb 14  2019 gpiochip216 -> ../../devices/platform/pinctrl/gpio/gpiochip216
lrwxrwxrwx  1 root root    0 Feb 14  2019 gpiochip24 -> ../../devices/platform/pinctrl/gpio/gpiochip24
lrwxrwxrwx  1 root root    0 Feb 14  2019 gpiochip248 -> ../../devices/platform/pinctrl/gpio/gpiochip248
lrwxrwxrwx  1 root root    0 Feb 14  2019 gpiochip56 -> ../../devices/platform/pinctrl/gpio/gpiochip56
lrwxrwxrwx  1 root root    0 Feb 14  2019 gpiochip88 -> ../../devices/platform/pinctrl/gpio/gpiochip88
```

Example for "Raspberry Pi Model B Rev 2" with "Raspbian GNU/Linux 11 (bullseye)" (26 pin header):

```sh
ls -la /sys/class/gpio/
total 0
drwxrwxr-x  2 root gpio    0 Nov 13 09:33 .
drwxr-xr-x 58 root root    0 Sep 13 02:58 ..
--w--w----  1 root gpio 4096 Nov 13 09:33 export
lrwxrwxrwx  1 root gpio    0 Nov 13 09:33 gpiochip0 -> ../../devices/platform/soc/20200000.gpio/gpio/gpiochip0
--w--w----  1 root gpio 4096 Nov 13 09:33 unexport

gpiodetect
gpiochip0 [pinctrl-bcm2835] (54 lines)

gpioinfo 
gpiochip0 - 54 lines:
  gpiochip0 - 54 lines:
  ...
  line   2:       "SDA1"       unused   input  active-high 
  line   3:       "SCL1"       unused   input  active-high 
  line   4:  "GPIO_GCLK"       unused   input  active-high 
  ...
  line   7:  "SPI_CE1_N"       unused   input  active-high 
  line   8:  "SPI_CE0_N"       unused   input  active-high 
  line   9:    "SPI_SDI"       unused   input  active-high
  line  10:    "SPI_SDO"       unused   input  active-high
  line  11:   "SPI_SCLK"       unused   input  active-high 
  ...
  line  14:       "TXD0"       unused   input  active-high 
  line  15:       "RXD0"       unused   input  active-high 
  ...
  line  17:     "GPIO17"       unused   input  active-high 
  line  18:     "GPIO18"       unused   input  active-high 
  ...
  line  22:     "GPIO22"       unused   input  active-high 
  line  23:     "GPIO23"       unused   input  active-high 
  line  24:     "GPIO24"       unused   input  active-high 
  line  25:     "GPIO25"       unused   input  active-high 
  ...
```

## General GPIO tests

For Tinkerboard and in general for all other boards:

* the name on system level differ from the header name (normally pin1..pin40)
* the mapping is done in gobot by a file named something like [pinmap.go](../platforms/asus/tinkerboard/pinmap.go)
* for the next tests the system level name is needed

Connect an oscilloscope or at least a meter to the pin (used header pin26 for example). For the output tests a LED with
a sufficient resistor to 3.3.V (e.g. 290 Ohm) can be used to detect the output state.  

> On Tinkerboard the pin26 relates to gpio251. For raspi it relates to gpio7.

### Creating gpio251 (sysfs Tinkerboard)

> Needs to be "root" for this to work. Switch to root user by "su -".

```sh
echo "251" > /sys/class/gpio/export
```

investigate result:

```sh
# cat /sys/class/gpio/gpio251/active_low 
0

# cat /sys/class/gpio/gpio251/direction 
in

# cat /sys/class/gpio/gpio251/edge 
none

# cat /sys/class/gpio/gpio251/value 
1
```

> The value can float to "1", if the input is open.

### Test input behavior of gpio251 (sysfs Tinkerboard)

> Be careful with connecting the input to GND or 3.3V directly, instead use an resistor with minimum 300 Ohm.

Connect the input header pin26 to GND.

```sh
# cat /sys/class/gpio/gpio251/direction 
in
# cat /sys/class/gpio/gpio251/value 
0
```

Connect the input header pin26 to +3.3V with an resistor (e.g. 1kOhm).

```sh
# cat /sys/class/gpio/gpio251/value 
1
```

### Test edge detection behavior of gpio251 (sysfs Tinkerboard)

investigate status:

```sh
# cat /sys/class/gpio/gpio251/edge
none
```

The file exists only if the pin can be configured as an interrupt generating input pin. To activate edge detection,
"rising", "falling", or "both" needs to be set.

```sh
# cat /sys/class/gpio/gpio251/value
1
```

If edge detection is activated, a poll will return only when the interrupt was triggered. The new value is written to
the beginning of the file.

> Not tested yet, not supported by gobot yet.

### Test output behavior of gpio251 (sysfs Tinkerboard)

Connect the output header pin26 to +3.3V with an resistor (e.g. 1kOhm leads to ~0.3mA, 300Ohm leads to ~10mA).

```sh
# echo "out" > /sys/class/gpio/gpio251/direction
# cat /sys/class/gpio/gpio251/direction
out
# cat /sys/class/gpio/gpio251/value
0

# echo "1" > /sys/class/gpio/gpio251/value
# cat /sys/class/gpio/gpio251/value
1

# echo "0" > /sys/class/gpio/gpio251/value
# cat /sys/class/gpio/gpio251/value
0
```

The meter should show "0V" for values of "0" and "3.3V", when the value was set to "1".

> For armbian and Tinkerboard the value remains to "1", although it was set to "0". In this case the pin is not usable
> as output.

### Test inverse output behavior of gpio251 (sysfs Tinkerboard)

```sh
# cat /sys/class/gpio/gpio251/value
0
# cat /sys/class/gpio/gpio251/active_low
0

# echo "1" > /sys/class/gpio/gpio251/active_low
# cat /sys/class/gpio/gpio251/value
1

# echo "0" > /sys/class/gpio/gpio251/value
# cat /sys/class/gpio/gpio251/value
0
```

The meter should show "0V" for values of "1" and "3.3V", when the value was set to "0".

### Test input behavior of gpio7 (cdev Raspi)

> Use --help to get some information of the command usage and options, e.g. "gpioget --help". Be careful with connecting
> the input to GND or 3.3V directly, instead use an resistor with minimum 300Ohm. Prefer to use the pull-up or pull-down
> feature, if working.

```sh
sudo gpioget 0 7
1
gpioinfo | grep 'line   7'
  line   7:  "SPI_CE1_N"       unused   input  active-high

sudo gpioget --bias=pull-down 0 7
0
```

>The value can float to "1", if the input is open. Most likely the raspi device has an internal pull-up resistor.
>Setting the bias is not possible for sysfs usage. This is one of the advantages of the new character device Kernel feature.

### Test output behavior of gpio7 (cdev Raspi)

```sh
sudo gpioset 0 7=0
gpioinfo | grep 'line   7'
  line   7:  "SPI_CE1_N"       unused  output  active-high

sudo gpioset 0 7=1
gpioinfo | grep 'line   7'
  line   7:  "SPI_CE1_N"       unused  output  active-high
```

The meter should show "0V" for values of "0" and "3.3V", when the value was set to "1".  
A connected LED with pull-up resistor lights up for setting to "0" (inverse).

### Test inverse output behavior of gpio7 (cdev Raspi)

```sh
sudo gpioset -l 0 7=0
sudo gpioset -l 0 7=1
```

The meter should show "0V" for values of "1" and "3.3V", when the value was set to "0" (inverse logic).  
A connected LED with pull-up resistor lights up for setting to "1" (inverse reversed).

> The gpioinfo seems to do not recognize the "active-low" set.

## Links

* <https://www.kernel.org/doc/html/latest/admin-guide/gpio/sysfs.html>
* <https://embeddedbits.org/linux-kernel-gpio-user-space-interface-embeddedbits-2/>
