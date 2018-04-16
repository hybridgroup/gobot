# SPI

This package provides drivers for [spi](https://en.wikipedia.org/wiki/Serial_Peripheral_Interface_Bus) devices. 

It currently must be used along with platforms such as the [Raspberry Pi](https://gobot.io/documentation/platforms/raspi) and [Beaglebone Black](https://gobot.io/documentation/platforms/beaglebone) that have adaptors that implement the needed SPI interface. 

The SPI implementation uses the awesome [periph.io](https://periph.io/) which currently only works on Linux systems.

## Getting Started

## Installing
```
go get -d -u gobot.io/x/gobot/...
```

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. 

The following spi Devices are currently supported:

- APA102 Programmable LEDs
- MCP3002 Analog/Digital Converter
- MCP3004 Analog/Digital Converter
- MCP3008 Analog/Digital Converter
- MCP3202 Analog/Digital Converter
- MCP3204 Analog/Digital Converter
- MCP3208 Analog/Digital Converter
- MCP3304 Analog/Digital Converter
- GoPiGo3 Robot

Drivers wanted! :)

The following spi Adaptors are currently supported:

- Raspberry Pi

Adaptors wanted too!
