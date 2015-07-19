# ODroid C1

The ODroid-C1 is an inexpensive and popular ARM based single board computer with digital & PWM GPIO, and i2c interfaces built in.

The ODroid-C1 is a credit-card-sized single-board computer developed in South Korea.

For more info about the ODroid-C1 device, click [here](http://www.hardkernel.com/main/products/prdt_info.php?g_code=G141578608433/).

## How to Install

First you must install the appropriate Go packages

```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/odroid
```

#### Cross compiling for the ODroid C1
You must first configure your Go environment for linux cross compiling

```bash
$ cd $GOROOT/src
$ GOOS=linux GOARCH=arm ./make.bash --no-clean

```

Then compile your Gobot program with

```bash
$ GOARM=7 GOARCH=arm GOOS=linux go build examples/odroidc1_blink.go
```

Then you can simply upload your program over the network from your host computer to the ODroid

```bash
$ scp odroidc1_blink odroid@192.168.1.xxx:/home/odroid/
```

and execute it on your ODroidC1 with

```bash
$ ./odroidc1_blink
```

## How to Use

```go
package main

import (
        "time"

        "github.com/hybridgroup/gobot"
        "github.com/hybridgroup/gobot/platforms/gpio"
        "github.com/hybridgroup/gobot/platforms/odroid/c1"
)

func main() {
        gbot := gobot.NewGobot()

        r := c1.NewODroidC1Adaptor("c1")
        led := gpio.NewLedDriver(r, "led", "7")

        work := func() {
                gobot.Every(1*time.Second, func() {
                        led.Toggle()
                })
        }

        robot := gobot.NewRobot("blinkBot",
                []gobot.Connection{r},
                []gobot.Device{led},
                work,
        )

        gbot.AddRobot(robot)

        gbot.Start()
}
```

## ODroid C1 Hardware Configuration

### PWM

If using pwm, ensure that the following has been run first to enable pwm, with `npwm` being set to the number of pwm ports:

```bash
sudo modprobe pwm-meson npwm=1
sudo modprobe pwm-ctrl
```

If using 1 pwm port, gpio pin 33 will be enabled.  If using 2, in addition gpio pin 19 will be enabled, but this will disable SPI which is shared on the same port.

### I2C

Dedicated pins for I2C are configured for GPIO, these pins can be configured as I2C bus while change the pin configuration. In order to change the configuration, you must load the driver.

```bash
sudo modprobe aml_i2c
```

If you have to load the driver every time whenever your ODROID-C1 start, simple you can register driver to /etc/modules.

```bash
sudo echo "aml_i2c" >> /etc/modules
```

### ADC

There 2 ADC input ports on the 40-pin header.
- ADC.AIN0 : Pin #40
- ADC.AIN1 : Pin #37

You can access the ADC inputs via sysfs nodes.

```bash
/sys/class/saradc/saradc_ch0 
/sys/class/saradc/saradc_ch1
```

ADC's maximum sample rate is 50kSPS with 10bit resolution (0~1023).
But the actual sample rate is 8kSPS if you access it via sysfs due to the limited file IO speed.

The ADC inputs are limited to 1.8Volt.
