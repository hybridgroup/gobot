# Curie

The Intel Curie is a tiny computer for the Internet of Things. It is the processor used by the [Arduino/Genuino 101](https://www.arduino.cc/en/Main/ArduinoBoard101) and the [Intel TinyTILE](https://software.intel.com/en-us/node/675623).

In addition to the GPIO, ADC, and I2C hardware interfaces, the Curie also has a built-in Inertial Measurement Unit (IMU), including an accelerometer, gyroscope, and thermometer.

For more info about the Intel Curie platform go to: [https://www.intel.com/content/www/us/en/products/boards-kits/curie.html](https://www.intel.com/content/www/us/en/products/boards-kits/curie.html).

## How to Install

You would normally install Go and Gobot on your computer. When you execute the Gobot program code, it communicates with the connected microcontroller using the [Firmata protocol](https://github.com/firmata/protocol), either using a serial port, or the Bluetooth LE wireless interface.

```
go get -d -u gobot.io/x/gobot/...
```

## How To Use

```go
package main

import (
	"log"
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
	"gobot.io/x/gobot/platforms/intel-iot/curie"
)

func main() {
	firmataAdaptor := firmata.NewAdaptor("/dev/ttyACM0")
	led := gpio.NewLedDriver(firmataAdaptor, "13")
	imu := curie.NewIMUDriver(firmataAdaptor)

	work := func() {
		imu.On("Accelerometer", func(data interface{}) {
			log.Println("Accelerometer", data)
		})

		imu.On("Gyroscope", func(data interface{}) {
			log.Println("Gyroscope", data)
		})

		imu.On("Temperature", func(data interface{}) {
			log.Println("Temperature", data)
		})

		imu.On("Motion", func(data interface{}) {
			log.Println("Motion", data)
		})

		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})

		gobot.Every(100*time.Millisecond, func() {
			imu.ReadAccelerometer()
			imu.ReadGyroscope()
			imu.ReadTemperature()
			imu.ReadMotion()
		})
	}

	robot := gobot.NewRobot("curieBot",
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{imu, led},
		work,
	)

	robot.Start()
}
```

## How to Connect

### Installing Firmware

You need to flash your Intel Curie with firmware that uses ConfigurableFirmata along with the FirmataCurieIMU plugin. There are 2 versions of this firmware, once that allows connecting using the serial interface, the other using a Bluetooth LE connection.

### Serial Port

Connect your Arduino 101 or TinyTILE using a serial cable, and connect using the correct serial port name.

### Bluetooth LE

Power up your Arduino 101 or TinyTILE using a battery or other power source, and connect using the BLE address or name.
