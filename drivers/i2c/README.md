# I2C

This package provides drivers for [i2c](https://en.wikipedia.org/wiki/I%C2%B2C)devices. It must be used along with an adaptor such as [firmata](https://gobot.io/x/gobot/platforms/firmata) that supports the needed interfaces for i2c devices.

## Getting Started

## Installing
```
go get -d -u gobot.io/x/gobot/...
```

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. The following i2c devices are currently supported:

- Adafruit Motor Hat
- BlinkM LED
- BMP180 Barometric Pressure/Temperature/Altitude Sensor
- DRV2605L Haptic Controller
- Grove Digital Accelerometer
- Grove RGB LCD
- HMC6352 Compass
- JHD1313M1 LCD Display w/RGB Backlight
- L3GD20H 3-Axis Gyroscope
- LIDAR-Lite
- MCP23017 Port Expander
- MMA7660 3-Axis Accelerometer
- MPL115A2 Barometer
- MPU6050 Accelerometer/Gyroscope
- SHT3x-D Temperature/Humidity
- TSL2561 Digital Luminosity/Lux/Light Sensor
- Wii Nunchuck Controller

More drivers are coming soon...

## Using A Different Bus or Address

You can set a different I2C address or I2C bus than the default when initializing your I2C drivers by using optional parameters. Here is an example:

```go
blinkm := i2c.NewBlinkMDriver(e, i2c.WithBus(0), i2c.WithAddress(0x09))
```
