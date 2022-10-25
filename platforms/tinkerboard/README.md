# Tinker Board

The ASUS Tinker Board is a single board SoC computer based on the Rockchip RK3288 processor. It has built-in GPIO, PWM, SPI, and I2C interfaces.

For more info about the Tinker Board, go to [https://www.asus.com/uk/Single-Board-Computer/Tinker-Board/](https://www.asus.com/uk/Single-Board-Computer/Tinker-Board/).

## How to Install

We recommend updating to the [latest Debian TinkerOS](https://github.com/TinkerBoard/debian_kernel/releases) when using the Tinker Board.

You would normally install Go and Gobot on your workstation. Once installed, cross compile your program on your workstation, transfer the final executable to your Tinker Board, and run the program on the Tinker Board as documented here.

```
go get -d -u gobot.io/x/gobot/...
```

### System access and configuration basics

Some configuration steps are needed to enable drivers and simplify the interaction with your Tinker Board. Once your Tinker Board has been configured, you do not need to do so again.

Note that these configuration steps must be performed on the Tinker Board itself. The easiest is to login to the Tinker Board via SSH (option "-4" is used to force IPv4, which is needed for some versions of TinkerOS):

```
ssh -4 linaro@192.168.1.xxx
```

### Enabling hardware drivers

Not all drivers are enabled by default. You can have a look at the configuration file, to find out what is enabled at your system:

```
cat /boot/config.txt
```

This file can be modified by "vi" or "nano", it is self explanatory:

```
sudo vi /boot/config.txt
```

Newer versions of Tinker Board provide an user interface for configuration with:

```
sudo tinker-config
```

After configuration was changed, an reboot is necessary.

```
sudo reboot
```

### Enabling GPIO pins

#### Create a group "gpio"

Create a Linux group named "gpio" by running the following command:

```
sudo groupadd -f --system gpio
```

If you already have a "gpio" group, you can skip to the next step.

#### Add the "linaro" user to the new "gpio" group

Add the user "linaro" to be a member of the Linux group named "gpio" by running the following command:

```
sudo usermod -a -G gpio linaro
```

If you already have added the "gpio" group, you can skip to the next step.

#### Add a "udev" rules file

Create a new "udev" rules file for the GPIO on the Tinker Board by running the following command:

```
sudo vi /etc/udev/rules.d/91-gpio.rules
```

And add the following contents to the file:

```
SUBSYSTEM=="gpio", KERNEL=="gpiochip*", ACTION=="add", PROGRAM="/bin/sh -c 'chown root:gpio /sys/class/gpio/export /sys/class/gpio/unexport ; chmod 220 /sys/class/gpio/export /sys/class/gpio/unexport'"
SUBSYSTEM=="gpio", KERNEL=="gpio*", ACTION=="add", PROGRAM="/bin/sh -c 'chown root:gpio /sys%p/active_low /sys%p/direction /sys%p/edge /sys%p/value ; chmod 660 /sys%p/active_low /sys%p/direction /sys%p/edge /sys%p/value'"
```

Press the "Esc" key, then press the ":" key and then the "q" key, and then press the "Enter" key. This should save your file. After rebooting your Tinker Board, you should be able to run your Gobot code that uses GPIO.

### Enabling I2C

#### Create a group "i2c"

Create a Linux group named "i2c" by running the following command:

```
sudo groupadd -f --system i2c
```

If you already have a "i2c" group, you can skip to the next step.

#### Add the "linaro" user to the new "i2c" group

Add the user "linaro" to be a member of the Linux group named "i2c" by running the following command:

```
sudo usermod -a -G gpio linaro
```

If you already have added the "i2c" group, you can skip to the next step.

#### Add a "udev" rules file

Create a new "udev" rules file for the I2C on the Tinker Board by running the following command:

```
sudo vi /etc/udev/rules.d/92-i2c.rules
```

And add the following contents to the file:

```
KERNEL=="i2c-0"     , GROUP="i2c", MODE="0660"
KERNEL=="i2c-[1-9]*", GROUP="i2c", MODE="0666"
```

Press the "Esc" key, then press the ":" key and then the "q" key, and then press the "Enter" key. This should save your file. After rebooting your Tinker Board, you should be able to run your Gobot code that uses I2C.

## How to Use

The pin numbering used by your Gobot program should match the way your board is labeled right on the board itself.

```go
r := tinkerboard.NewAdaptor()
led := gpio.NewLedDriver(r, "7")
```

## How to Connect

### Compiling

Compile your Gobot program on your workstation like this:

```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/tinkerboard_blink.go
```

Once you have compiled your code, you can upload your program and execute it on the Tinkerboard from your workstation using the `scp` and `ssh` commands like this:

```bash
$ scp tinkerboard_blink linaro@192.168.1.xxx:/home/linaro/
$ ssh -t linaro@192.168.1.xxx "./tinkerboard_blink"
```

## Troubleshooting

### PWM

#### Investigate state
```
$ ls /sys/class/pwm/
pwmchip0  pwmchip1

```
pwmchip0 and pwmchip1 seems to be not usable (unknown, which functionality make use of that, maybe fan?)

#### Activate
When there is no pwm2, pwm3, this can be activated.
Open the file /boot/config.txt and add lines/remove comment sign and adjust:

intf:pwm2=on
intf:pwm3=on

Then save the file, close and reboot.

After reboot check the state:
```
$ ls /sys/class/pwm/
pwmchip0  pwmchip1  pwmchip2  pwmchip3

```
pwmchip2: the pin 33 is usable
pwmchip3: the pin 32 is usable

When only one pwm was enabled it will be always "pwmchip2". In this case, the activated pin can be found by investigating the symbolic link of the device:
```
# ls -la /sys/class/pwm/
total 0
drwxr-xr-x  2 root root 0 Apr 24 14:11 .
drwxr-xr-x 66 root root 0 Apr 24 14:09 ..
lrwxrwxrwx  1 root root 0 Apr 24 14:09 pwmchip0 -> ../../devices/platform/ff680000.pwm/pwm/pwmchip0
lrwxrwxrwx  1 root root 0 Apr 24 14:09 pwmchip1 -> ../../devices/platform/ff680010.pwm/pwm/pwmchip1
lrwxrwxrwx  1 root root 0 Apr 24 14:09 pwmchip2 -> ../../devices/platform/ff680030.pwm/pwm/pwmchip2
```

ff680020 => pwm2, pin33
ff680030 => pwm3, pin32

#### Test
For example only pwm3 was activated to use pin32. Connect an oscilloscope or at least a meter to the pin 32.

switch to root user by "su -"

investigate state:
```
# ls -la /sys/class/pwm/pwmchip2/
total 0
drwxr-xr-x 3 root root    0 Apr 24 14:17 .
drwxr-xr-x 3 root root    0 Apr 24 14:17 ..
lrwxrwxrwx 1 root root    0 Apr 24 14:17 device -> ../../../ff680030.pwm
--w------- 1 root root 4096 Apr 24 14:17 export
-r--r--r-- 1 root root 4096 Apr 24 14:17 npwm
drwxr-xr-x 2 root root    0 Apr 24 14:17 power
lrwxrwxrwx 1 root root    0 Apr 24 14:17 subsystem -> ../../../../../class/pwm
-rw-r--r-- 1 root root 4096 Apr 24 14:17 uevent
--w------- 1 root root 4096 Apr 24 14:17 unexport
```
#### Creating pwm0
`echo 0 > /sys/class/pwm/pwmchip2/enable`
investigate result:
```
# ls /sys/class/pwm/pwmchip2/
device	export	npwm  power  pwm0  subsystem  uevent  unexport
# ls /sys/class/pwm/pwmchip2/pwm0/
capture  duty_cycle  enable  period  polarity  power  uevent
# cat /sys/class/pwm/pwmchip2/pwm0/period 
0
# cat /sys/class/pwm/pwmchip2/pwm0/duty_cycle 
0
# cat /sys/class/pwm/pwmchip2/pwm0/enable 
0
# cat /sys/class/pwm/pwmchip2/pwm0/polarity 
inversed
```
#### Initialization
Note: Before writing the period all other write actions will cause an error "-bash: echo: write error: Invalid argument"
```
echo 10000000 > /sys/class/pwm/pwmchip2/pwm0/period # this is a frequency divider for 1GHz (1000 will produce a frequency of 1MHz, 1000000 will cause a frequency of 1kHz, her we got 100Hz)
echo "normal" > /sys/class/pwm/pwmchip2/pwm0/polarity
echo 3000000 > /sys/class/pwm/pwmchip2/pwm0/duty_cycle # this means 30%
echo 1  > /sys/class/pwm/pwmchip2/pwm0/enable

```
Now we should measure a value of around 1V with the meter, because the basis value is 3.3V and 30% leads to 1V.

Try to inverse the sequence:
`echo "inversed" > /sys/class/pwm/pwmchip2/pwm0/polarity`
Now we should measure a value of around 2.3V with the meter, which is the difference of 1V to 3.3V.

If we have attached an oscilloscope we can play around with the values for period and duty_cycle and see what happen.
