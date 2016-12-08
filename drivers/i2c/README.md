# I2C

This package provides drivers for [i2c](https://en.wikipedia.org/wiki/I%C2%B2C)devices. It must be used along with an adaptor such as [firmata](https://gobot.io/x/gobot/platforms/firmata) that supports the needed interfaces for i2c devices.

## Getting Started

## Installing
```
go get -d -u gobot.io/x/gobot/... && go install gobot.io/x/gobot/platforms/i2c
```

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. The following i2c devices are currently supported:

- BlinkM
- Grove Digital Accelerometer
- Grove RGB LCD
- HMC6352 Compass
- JHD1313M1 RGB LCD Display
- LIDAR-Lite
- MCP23017 Port Expander
- MMA7660 3-Axis Accelerometer
- MPL115A2 Barometer
- MPU6050 Accelerometer/Gyroscope
- Wii Nunchuck Controller

More drivers are coming soon...
