# SPI

This package provides drivers for [SPI](https://en.wikipedia.org/wiki/Serial_Peripheral_Interface_Bus) devices.

## Getting Started

Please refer to the main [README.md](https://github.com/hybridgroup/gobot/blob/release/README.md)

## Hardware Support

Gobot has a extensible system for connecting to hardware devices.

The following SPI Devices are currently supported:

- APA102 Programmable LEDs
- MCP3002 Analog/Digital Converter
- MCP3004 Analog/Digital Converter
- MCP3008 Analog/Digital Converter
- MCP3202 Analog/Digital Converter
- MCP3204 Analog/Digital Converter
- MCP3208 Analog/Digital Converter
- MCP3304 Analog/Digital Converter
- MFRC522 RFID Card Reader
- SSD1306 OLED Display Controller
- GoPiGo3 Robot

The following SPI system drivers are currently supported:

- SPI by `/dev/spidevX.Y` with the awesome [periph.io](https://periph.io/) which currently only works on Linux systems
- SPI via GPIO's
