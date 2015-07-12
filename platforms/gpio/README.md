# GPIO

This package provides drivers for [General Purpose Input/Output (GPIO)](https://en.wikipedia.org/wiki/General_Purpose_Input/Output) devices . It is normally not used directly, but instead is registered by an adaptor such as [firmata](https://github.com/hybridgroup/gobot/platforms/firmata) that supports the needed interfaces for GPIO devices.

## Getting Started

## Installing
```
go get -d -u github.com/hybridgroup/gobot/... && go install github.com/hybridgroup/gobot/platforms/gpio
```

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. The following GPIO devices are currently supported:

  - Analog Sensor
  - Button
  - Direct Pin
  - LED
  - Makey Button
  - Motor
  - Servo

More drivers are coming soon...