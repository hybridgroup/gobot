# SPI

This package provides drivers for [spi](https://en.wikipedia.org/wiki/Serial_Peripheral_Interface_Bus) devices. It must be used along with the [raspberry pi](https://gobot.io/documentation/platforms/raspi) adaptor that supports the needed interfaces for spi devices. This uses the experimental [spi package](https://github.com/golang/exp/tree/master/io/spi) which only works on linux systems.

## Getting Started

## Installing
```
go get -d -u gobot.io/x/gobot/...
```

## Hardware Support
Gobot has a extensible system for connecting to hardware devices. The following spi devices are currently supported:

More drivers are coming soon...

## Using A Different Bus or Address

You can set a different SPI address or SPI bus than the default when initializing your SPI drivers by using optional parameters. Here is an example:

```go
blinkm := spi.NewGoPiGo3DriverDriver(e, spi.WithBus(0), spi.WithAddress(0x10), spi.With)
```