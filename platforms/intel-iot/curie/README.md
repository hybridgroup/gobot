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

		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})

		gobot.Every(100*time.Millisecond, func() {
			imu.ReadAccelerometer()
			imu.ReadGyroscope()
			imu.ReadTemperature()
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

To setup your Arduino environment:

- Install the latest Arduino, if you have not done so yet.
- Install the "Intel Curie Boards" board files using the "Board Manager". You can find it in the Arduino IDE under the "Tools" menu. Choose "Boards > Boards Manager".

	Search for the "Intel Curie Boards" package in the "Boards Manager" dialog, and then install the latest version.

- Download the ZIP file for the ConfigurableFirmata library. You can download the latest version of the ConfigurableFirmata from here:

	[https://github.com/firmata/ConfigurableFirmata/archive/master.zip](https://github.com/firmata/ConfigurableFirmata/archive/master.zip)

	Once you have downloaded ConfigurableFirmata, install it by using the "Library Manager". You can find it in the Arduino IDE under the "Sketch" menu. Choose "Include Library > Add .ZIP Library". Select the ZIP file for the ConfigurableFirmata library that you just downloaded.

- Download the ZIP file for the FirmataCurieIMU library. You can download the latest version of FirmataCurieIMU from here:

	[https://github.com/intel-iot-devkit/firmata-curie-imu/archive/master.zip](https://github.com/intel-iot-devkit/firmata-curie-imu/archive/master.zip)

	Once you have downloaded the FirmataCurieIMU library, install it by using the "Library Manager". You can find it in the Arduino IDE under the "Sketch" menu. Choose "Include Library > Add .ZIP Library". Select the ZIP file for the FirmataCurieIMU library that you just downloaded.

- Linux only: On some Linux distributions, additional device rules are required in order to connect to the board. Run the following command then unplug the board and plug it back in before proceeding:

  ```sh
  curl -sL https://raw.githubusercontent.com/01org/corelibs-arduino101/master/scripts/create_dfu_udev_rule | sudo -E bash -
  ```

Now you are ready to install your firmware. You must decide if you want to connect via the serial port, or using Bluetooth LE.

### Serial Port

To use your Intel Curie connected via serial port, you should use the sketch located here:

[https://github.com/intel-iot-devkit/firmata-curie-imu/blob/master/examples/everythingIMU/everythingIMU.ino](https://github.com/intel-iot-devkit/firmata-curie-imu/blob/master/examples/everythingIMU/everythingIMU.ino)

Once you have loaded this sketch on your Intel Curie, you can run your Gobot code to communicate with it. Leave your Arduino 101 or TinyTILE connected using the serial cable that you used to flash the firmware, and refer to that same serial port name in your Gobot code.

### Bluetooth LE

To use your Intel Curie connected via Bluetooth LE, you should use the sketch located here:

[https://github.com/intel-iot-devkit/firmata-curie-imu/blob/master/examples/bleIMU/bleIMU.ino](https://github.com/intel-iot-devkit/firmata-curie-imu/blob/master/examples/bleIMU/bleIMU.ino)

Once you have loaded this sketch on your Intel Curie, you can run your Gobot code to communicate with it.

Power up your Arduino 101 or TinyTILE using a battery or other power source, and connect using the BLE address or name. The default BLE name is "FIRMATA".
