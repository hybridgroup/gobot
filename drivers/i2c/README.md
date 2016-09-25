# I2C

This package provides drivers for [i2c](https://en.wikipedia.org/wiki/I%C2%B2C)devices . It is normally not used directly, but instead is registered by an adaptor such as [firmata](https://github.com/hybridgroup/gobot/platforms/firmata) that supports the needed interfaces for i2c devices.

## Getting Started

## Installing
```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/i2c
```

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. The following i2c devices are currently supported:

- BlinkM
- HMC6352 Digital Compass
- MPL115A2 Barometer/Temperature Sensor
- MPU6050 Accelerometer/Gyroscope
- Wii Nunchuck Controller

More drivers are coming soon...
