# I2C

This package provides drivers for [i2c](https://en.wikipedia.org/wiki/I%C2%B2C)devices. It must be used along with an adaptor such as [firmata](https://gobot.io/x/gobot/platforms/firmata) that supports the needed interfaces for i2c devices.

## Getting Started

## Installing
```
go get -d -u gobot.io/x/gobot/...
```

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. The following i2c devices are currently supported:

- Adafruit 2x16 RGB-LCD with 5 keys
- Adafruit Motor Hat
- ADS1015 Analog to Digital Converter
- ADS1115 Analog to Digital Converter
- ADXL345 Digital Accelerometer
- BH1750 Digital Luminosity/Lux/Light Sensor
- BlinkM LED
- BME280 Barometric Pressure/Temperature/Altitude/Humidity Sensor
- BMP180 Barometric Pressure/Temperature/Altitude Sensor
- BMP280 Barometric Pressure/Temperature/Altitude Sensor
- BMP388 Barometric Pressure/Temperature/Altitude Sensor
- DRV2605L Haptic Controller
- Generic driver for read and write values to/from register address
- Grove Digital Accelerometer
- GrovePi Expansion Board
- Grove RGB LCD
- HMC6352 Compass
- HMC5883L 3-Axis Digital Compass
- INA3221 Voltage Monitor
- JHD1313M1 LCD Display w/RGB Backlight
- L3GD20H 3-Axis Gyroscope
- LIDAR-Lite
- MCP23017 Port Expander
- MMA7660 3-Axis Accelerometer
- MPL115A2 Barometric Pressure/Temperature
- MPU6050 Accelerometer/Gyroscope
- PCA9501 8-bit I/O port with interrupt, 2-kbit EEPROM
- PCA9685 16-channel 12-bit PWM/Servo Driver
- PCF8583 clock and calendar or event counter, 240 x 8-bit RAM
- PCF8591 8-bit 4xA/D & 1xD/A converter
- SHT2x Temperature/Humidity
- SHT3x-D Temperature/Humidity
- SSD1306 OLED Display Controller
- TSL2561 Digital Luminosity/Lux/Light Sensor
- Wii Nunchuck Controller
- YL-40 Brightness/Temperature sensor, Potentiometer, analog input, analog output Driver

More drivers are coming soon...

## Using A Different Bus or Address

You can set a different I2C address or I2C bus than the default when initializing your I2C drivers by using optional parameters. Here is an example:

```go
blinkm := i2c.NewBlinkMDriver(e, i2c.WithBus(0), i2c.WithAddress(0x09))
```
