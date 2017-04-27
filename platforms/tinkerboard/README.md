# Tinker Board

The ASUS Tinker Board is a single board SoC computer based on the Rockchip RK3288 processor. It has built-in GPIO, PWM, SPI, and I2C interfaces.

For more info about the Tinker Board, go to [https://www.asus.com/uk/Single-Board-Computer/Tinker-Board/](https://www.asus.com/uk/Single-Board-Computer/Tinker-Board/).

## How to Install

We recommend updating to the latest Debian TinkerOS when using the Tinker Board.

You would normally install Go and Gobot on your workstation. Once installed, cross compile your program on your workstation, transfer the final executable to your Tinker Board, and run the program on the Tinker Board as documented here.

```
go get -d -u gobot.io/x/gobot/...
```

### Enabling GPIO pins

To enable use of the Tinker Board GPIO pins, you need to perform the following steps as a one-time configuration. Once your Tinker Board has been configured, you do not need to do so again.

Note that these configuration steps must be performed on the Tinker Board itself. The easiest is to login to the Tinker Board via SSH:

```
ssh linaro@192.168.1.xxx
```

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

To enable use of the Tinker Board I2C, you need to perform the following steps as a one-time configuration. Once your Tinker Board has been configured, you do not need to do so again.

Note that these configuration steps must be performed on the Tinker Board itself. The easiest is to login to the Tinker Board via SSH:

```
ssh linaro@192.168.1.xxx
```

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

Once you have compiled your code, you can you can upload your program and execute it on the Tinkerboard from your workstation using the `scp` and `ssh` commands like this:

```bash
$ scp tinkerboard_blink linaro@192.168.1.xxx:/home/linaro/
$ ssh -t linaro@192.168.1.xxx "./tinkerboard_blink"
```
