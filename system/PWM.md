# PWM's

This document describes some basics for developers.

## Check the PWM features

Example for Tinkerboard:

```sh
# ls -la /sys/class/pwm/
total 0
drwxr-xr-x  2 root root 0 Apr 24 14:11 .
drwxr-xr-x 66 root root 0 Apr 24 14:09 ..
lrwxrwxrwx  1 root root 0 Apr 24 14:09 pwmchip0 -> ../../devices/platform/ff680000.pwm/pwm/pwmchip0
lrwxrwxrwx  1 root root 0 Apr 24 14:09 pwmchip1 -> ../../devices/platform/ff680010.pwm/pwm/pwmchip1
lrwxrwxrwx  1 root root 0 Apr 24 14:09 pwmchip2 -> ../../devices/platform/ff680030.pwm/pwm/pwmchip2
```

```sh
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

## General PWM tests

Connect an oscilloscope or at least a meter to the pin (used pin32 for example with Tinkerboard).  
Switch to root user by "su -".

### Investigate state of PWMs

For Tinkerboard:

* ff680000 and ff680010 seems to be not usable (unknown, which functionality make use of that, maybe fan?)
* if there is no ff680020, ff680030, this can be activated. See section [Change available features](README.md#change-available-features)

### Creating pwm0 on pwmchip2

```sh
echo 0 > /sys/class/pwm/pwmchip2/export
```

investigate result:

```sh
# ls /sys/class/pwm/pwmchip2/
device  export  npwm  power  pwm0  subsystem  uevent  unexport

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

### Initialization of pwm0

```sh
echo 10000000 > /sys/class/pwm/pwmchip2/pwm0/period
echo "normal" > /sys/class/pwm/pwmchip2/pwm0/polarity
echo 3000000 > /sys/class/pwm/pwmchip2/pwm0/duty_cycle # this means 30%
echo 1  > /sys/class/pwm/pwmchip2/pwm0/enable

```

> Before writing the period, all other write actions will cause an error "-bash: echo: write error: Invalid argument".
> The "period" is in nanoseconds (or a frequency divider for 1GHz), 1000 will produce a frequency of 1MHz, 1000000 will
> cause a frequency of 1kHz. For the example 10000000 we have 100Hz.

Now we should measure a value of around 1V with the meter, because the basis value is 3.3V and 30% leads to 1V.

Try to inverse the sequence:
`echo "inversed" > /sys/class/pwm/pwmchip2/pwm0/polarity`
Now we should measure a value of around 2.3V with the meter, which is the difference of 1V to 3.3V.

If we have attached an oscilloscope we can play around with the values for period and duty_cycle and see what happen.

## Links

* <https://docs.kernel.org/driver-api/pwm.html>
