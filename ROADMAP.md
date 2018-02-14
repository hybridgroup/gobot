# Roadmap

This is the roadmap of what we as a community want to see happen with Gobot. It should be considered more as a statement of direction then a list of tasks.

Requests for changes to the roadmap should be made in the form of pull requests to this document.

Anything tied to any implementation, including requests for platform support, bug reports, or other specifics should still be made by creating a new issue here:

https://github.com/hybridgroup/gobot/issues

## core

- standardized logging
- use Context to allow for graceful exits.

## api

- ability to plug in your own router to handle API calls, for example to serve a custom web app.
- restrict API calls to only specific set of entrypoints.
- serve other transports/protocols other than HTTP/REST for example CoAP.

## gpio

- support for epoll/interrupt based gpio events.
- helper method for interrupts to handle "ping" timing-based devices.
- Windows 10 support.
- use variadic constructor functions to allow for additional params, similar to i2c drivers.

## aio

- support for epoll based aio events possible?
- Windows 10 support.
- use variadic constructor functions to allow for additional params, similar to i2c drivers.

## i2c

- ensure that SMBUS operations are working as expected.
- add support for the following i2c devices:
   - HMC5883L
   - LSM303DLHC
   - MAG3110
   - MMA8452
   - PCF8591
   - T5403
   - TMP006
   - VCNL4000

## 1-wire

- add support for 1-wire protocol.

## serial

- create a common serial Adaptor, so different serial devices such as GPS, LIDAR etc only need to implement drivers.

## ble

- improve the ble package to allow support for multiple peripherals.
- Windows 10 support.
